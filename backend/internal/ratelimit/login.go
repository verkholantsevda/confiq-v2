package ratelimit

import (
	"sync"
	"time"
)

// bucket хранит времена неудачных попыток и, при превышении лимита,
// момент до которого доступ запрещён.
type bucket struct {
	attempts     []time.Time
	blockedUntil time.Time
}

func (b *bucket) prune(now time.Time, window time.Duration) {
	if b == nil || len(b.attempts) == 0 {
		return
	}

	cutoff := now.Add(-window)
	i := 0
	for i < len(b.attempts) && b.attempts[i].Before(cutoff) {
		i++
	}
	if i > 0 {
		b.attempts = append([]time.Time(nil), b.attempts[i:]...)
	}
}

// LoginGuard — простой in-memory лимитер для логина и проверки TOTP.
// Экземпляр один на процесс; при нескольких репликах понадобится общий store.
type LoginGuard struct {
	mu sync.Mutex

	accounts map[string]*bucket
	totp     map[string]*bucket
	ips      map[string]*bucket

	accountLimit  int
	accountWindow time.Duration
	accountBlock  time.Duration
	totpLimit     int
	totpWindow    time.Duration
	totpBlock     time.Duration
	ipLimit       int
	ipWindow      time.Duration
}

func NewLoginGuard() *LoginGuard {
	return &LoginGuard{
		accounts: make(map[string]*bucket),
		totp:     make(map[string]*bucket),
		ips:      make(map[string]*bucket),

		accountLimit:  5,
		accountWindow: 15 * time.Minute,
		accountBlock:  15 * time.Minute,
		totpLimit:     5,
		totpWindow:    5 * time.Minute,
		totpBlock:     5 * time.Minute,
		ipLimit:       20,
		ipWindow:      time.Minute,
	}
}

// Check возвращает false и время до снятия ограничения, если вход запрещён.
func (g *LoginGuard) Check(username, ip string) (bool, time.Duration) {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now()

	if d := blockedFor(g.accounts[username], now); d > 0 {
		return false, d
	}
	if d := blockedFor(g.totp[username], now); d > 0 {
		return false, d
	}

	if b := g.ips[ip]; b != nil {
		b.prune(now, g.ipWindow)
		if len(b.attempts) >= g.ipLimit {
			return false, b.attempts[0].Add(g.ipWindow).Sub(now)
		}
	}

	return true, 0
}

func (g *LoginGuard) Fail(username, ip string, totpFailure bool) {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now()

	if totpFailure {
		g.record(g.totp, username, now, g.totpWindow, g.totpLimit, g.totpBlock)
	} else {
		g.record(g.accounts, username, now, g.accountWindow, g.accountLimit, g.accountBlock)
	}
	g.record(g.ips, ip, now, g.ipWindow, g.ipLimit, 0)

	g.cleanup(now)
}

func (g *LoginGuard) Success(username, ip string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	delete(g.accounts, username)
	delete(g.totp, username)
	delete(g.ips, ip)
}

func (g *LoginGuard) record(
	buckets map[string]*bucket,
	key string,
	now time.Time,
	window time.Duration,
	limit int,
	block time.Duration,
) {
	b := buckets[key]
	if b == nil {
		b = &bucket{}
		buckets[key] = b
	}

	b.prune(now, window)
	b.attempts = append(b.attempts, now)

	if block > 0 && len(b.attempts) >= limit {
		b.blockedUntil = now.Add(block)
		b.attempts = nil
	}
}

func (g *LoginGuard) cleanup(now time.Time) {
	for key, b := range g.accounts {
		b.prune(now, g.accountWindow)
		if len(b.attempts) == 0 && !b.blockedUntil.After(now) {
			delete(g.accounts, key)
		}
	}

	for key, b := range g.totp {
		b.prune(now, g.totpWindow)
		if len(b.attempts) == 0 && !b.blockedUntil.After(now) {
			delete(g.totp, key)
		}
	}

	for key, b := range g.ips {
		b.prune(now, g.ipWindow)
		if len(b.attempts) == 0 {
			delete(g.ips, key)
		}
	}
}

func blockedFor(b *bucket, now time.Time) time.Duration {
	if b == nil || !b.blockedUntil.After(now) {
		return 0
	}

	return b.blockedUntil.Sub(now)
}
