# ---------- Stage 1: Build ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app


# Кэш зависимостей
COPY go.mod ./
RUN go mod tidy
RUN go mod download

# Копируем исходники
COPY . .

# Сборка бинарника
RUN go build -ldflags="-s -w" -o app ./cmd/main.go

# ---------- Stage 2: Runtime ----------
FROM alpine:latest


WORKDIR /app


# Копируем бинарник из Stage 1
COPY --from=builder /app/app .

# Копируем конфиг
COPY config/local.yaml /app/config/local.yaml

ENV CONFIG_PATH=/app/config/local.yaml

EXPOSE 8080

CMD ["./app"]
