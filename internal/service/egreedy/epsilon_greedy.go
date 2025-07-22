// Package egreedy реализует алгоритм ε‑жадного выбора баннеров.
package egreedy

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/Sucsz/banner-rotator/internal/db/dao"
)

// Service — алгоритм ε‑greedy, безопасный для конкурентного использования.
type Service struct {
	eps     float64
	statDAO dao.StatDAO
	slotDAO dao.BannerSlotDAO

	mu  sync.Mutex // защита rnd
	rnd *rand.Rand
}

// NewEpsilonGreedy создаёт Service с генератором, засеянным текущим UnixNano.
//
//nolint:gosec
func NewEpsilonGreedy(
	eps float64,
	statDAO dao.StatDAO,
	slotDAO dao.BannerSlotDAO,
) *Service {
	src := rand.NewSource(time.Now().UnixNano())
	return &Service{
		eps:     eps,
		statDAO: statDAO,
		slotDAO: slotDAO,
		rnd:     rand.New(src),
	}
}

// NewEpsilonGreedyWithRND создаёт Service с уже готовым rnd (для тестов).
func NewEpsilonGreedyWithRND(
	eps float64,
	statDAO dao.StatDAO,
	slotDAO dao.BannerSlotDAO,
	rnd *rand.Rand,
) *Service {
	return &Service{
		eps:     eps,
		statDAO: statDAO,
		slotDAO: slotDAO,
		rnd:     rnd,
	}
}

// Select выбирает баннер для показа: с вероятностью eps — случайный (explore),
// иначе — лучший по CTR (exploit). После выбора инкрементит показ.
func (s *Service) Select(ctx context.Context, slotID, groupID int64) (bannerID int64, err error) {
	// 1) Получаем список баннеров
	ids, err := s.slotDAO.GetBannersBySlot(ctx, slotID)
	if err != nil {
		return 0, err
	}

	// 2) Случайное число в [0,1), реализуем вероятность
	s.mu.Lock()
	r := s.rnd.Float64()
	s.mu.Unlock()

	if r < s.eps {
		// explore: случайный индекс
		s.mu.Lock()
		idx := s.rnd.Intn(len(ids))
		s.mu.Unlock()
		bannerID = ids[idx]
	} else {
		// exploit: лучший по CTR
		var bestID int64
		var bestCTR float64
		for _, id := range ids {
			st, err := s.statDAO.Get(ctx, slotID, id, groupID)
			if err != nil {
				return 0, err
			}
			ctr := float64(st.Clicks) / float64(st.Impressions+1)
			if ctr > bestCTR {
				bestCTR = ctr
				bestID = id
			}
		}
		bannerID = bestID
	}

	// 3) Инкрементим показ
	if err := s.statDAO.IncrementView(ctx, slotID, bannerID, groupID); err != nil {
		return 0, err
	}
	return bannerID, nil
}

// RecordClick нужно вызывать при клике, чтобы увеличить счётчик кликов.
func (s *Service) RecordClick(
	ctx context.Context,
	slotID, bannerID, groupID int64,
) error {
	return s.statDAO.IncrementClick(ctx, slotID, bannerID, groupID)
}
