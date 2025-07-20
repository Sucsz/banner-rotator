package bandit

import (
	"context"

	"github.com/Sucsz/banner-rotator/internal/db/dao"
	"github.com/Sucsz/banner-rotator/internal/service/egreedy"
)

// BannerSelector выбирает баннер и сразу инкрементит показ.
type BannerSelector interface {
	Select(ctx context.Context, slotID, groupID int64) (bannerID int64, err error)
}

// Config параметры алгоритма.
type Config struct {
	// Доля случайных выборов (0.0–1.0).
	Epsilon float64
}

// NewBandit создаёт ε‑greedy селектор.
// statDAO — для работы со статистикой (клики/показы),
// slotDAO — для получения списка баннеров в слоте.
func NewBandit(
	cfg Config,
	statDAO dao.StatDAO,
	slotDAO dao.BannerSlotDAO,
) BannerSelector {
	return egreedy.NewEpsilonGreedy(
		cfg.Epsilon,
		statDAO,
		slotDAO,
	)
}
