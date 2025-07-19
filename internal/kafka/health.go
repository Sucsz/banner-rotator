package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

// CheckConnection пытается установить TCP‑соединение с любым из брокеров
// и сразу закрывает его. Возвращает ошибку, если ни один брокер недоступен.
func CheckConnection(brokers []string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Попробуем первый брокер — обычно достаточно
	conn, err := kafka.DialContext(ctx, "tcp", brokers[0])
	if err != nil {
		return err
	}
	return conn.Close()
}
