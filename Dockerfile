# ─── Стадия сборки ────────────────────────────────────────────────────────────
FROM golang:1.23.4 AS builder

# Отключаем cgo и заставляем Go использовать собственный резолвер DNS
ENV CGO_ENABLED=0
ENV GODEBUG=netdns=go

WORKDIR /app

# Кэшируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Собираем строго статический бинарник (-static) и обрезаем отладочные символы (-s -w)
RUN go build -tags netgo \
    -ldflags="-s -w -extldflags '-static'" \
    -o banner-rotator ./cmd/main.go

# ─── Финальный минимальный образ ─────────────────────────────────────────────
FROM scratch

# Копируем бинарник из builder-стадии
COPY --from=builder /app/banner-rotator /banner-rotator

# Документируем порт (пробросит compose)
EXPOSE 8080

# Запуск приложения
ENTRYPOINT ["/banner-rotator"]
