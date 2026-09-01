# ---------- Build ----------
FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache \
    git \
    gcc \
    musl-dev \
    sqlite-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags="-s -w" \
    -o confiq \
    ./cmd/server

# ---------- Runtime ----------
FROM alpine:3.22

RUN apk add --no-cache ca-certificates sqlite-libs

WORKDIR /app

# Копируем бинарник
COPY --from=builder /app/confiq .

# Копируем папку с миграциями (добавленная строка)
COPY --from=builder /app/migrations ./migrations

EXPOSE 8089

CMD ["./confiq"]