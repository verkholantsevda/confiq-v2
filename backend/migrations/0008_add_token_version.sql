ALTER TABLE users ADD COLUMN token_version INTEGER NOT NULL DEFAULT 0;

-- Токены, выданные до обновления, не содержат token_version в claims.
-- Увеличиваем версию, чтобы разово разлогинить все старые сессии.
UPDATE users SET token_version = 1;
