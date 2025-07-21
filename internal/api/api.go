package api

import (
	"github.com/Sucsz/banner-rotator/internal/db/dao"
	"github.com/Sucsz/banner-rotator/internal/kafka"
	"github.com/Sucsz/banner-rotator/internal/service/bandit"
)

// API включает все зависимости для HTTP‑хендлеров.
type API struct {
	Selector      bandit.BannerSelector
	BannerDAO     dao.BannerDAO
	BannerSlotDAO dao.BannerSlotDAO
	StatDAO       dao.StatDAO
	Producer      kafka.Producer
}

// NewAPI создаёт новый API‑объект со всеми зависимостями.
func NewAPI(
	selector bandit.BannerSelector,
	producer kafka.Producer,
	bannerSlotDAO dao.BannerSlotDAO,
) *API {
	return &API{
		Selector:      selector,
		Producer:      producer,
		BannerSlotDAO: bannerSlotDAO,
	}
}
