//nolint:revive
package egreedy_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"

	"github.com/Sucsz/banner-rotator/internal/db/model"
	"github.com/Sucsz/banner-rotator/internal/service/egreedy"
)

// fakeSlotDAO реализует dao.BannerSlotDAO полностью.
type fakeSlotDAO struct {
	banners []int64
}

func (f *fakeSlotDAO) AddBannerToSlot(ctx context.Context, bannerID, slotID int64) error {
	// noop для тестов
	return nil
}

func (f *fakeSlotDAO) RemoveBannerFromSlot(ctx context.Context, bannerID, slotID int64) error {
	// noop для тестов
	return nil
}

func (f *fakeSlotDAO) IsBannerInSlot(ctx context.Context, bannerID, slotID int64) (bool, error) {
	// считаем, что всегда true
	return true, nil
}

func (f *fakeSlotDAO) GetBannersBySlot(ctx context.Context, slotID int64) ([]int64, error) {
	return f.banners, nil
}

// fakeStatDAO реализует dao.StatDAO, собирая вызовы IncrementView.
type fakeStatDAO struct {
	stats     map[[3]int64]*model.BannerStat
	viewCalls []struct{ SlotID, BannerID, GroupID int64 }
}

func (f *fakeStatDAO) Get(ctx context.Context, slotID, bannerID, groupID int64) (*model.BannerStat, error) {
	key := [3]int64{slotID, bannerID, groupID}
	if st, ok := f.stats[key]; ok {
		return st, nil
	}
	// по умолчанию нулевая статистика
	return &model.BannerStat{
		BannerID:    bannerID,
		SlotID:      slotID,
		UserGroupID: groupID,
		Impressions: 0,
		Clicks:      0,
	}, nil
}

func (f *fakeStatDAO) IncrementView(ctx context.Context, slotID, bannerID, groupID int64) error {
	f.viewCalls = append(f.viewCalls, struct{ SlotID, BannerID, GroupID int64 }{slotID, bannerID, groupID})
	return nil
}

func (f *fakeStatDAO) IncrementClick(ctx context.Context, slotID, bannerID, groupID int64) error {
	// noop
	return nil
}

//nolint:gosec
func TestSelect_Explore(t *testing.T) {
	ctx := context.Background()
	banners := []int64{1, 2, 3}

	slotDAO := &fakeSlotDAO{banners: banners}
	statDAO := &fakeStatDAO{stats: make(map[[3]int64]*model.BannerStat)}
	// rnd с фиксированным сидом, чтобы тесты были детерминированными
	rnd := rand.New(rand.NewSource(17))

	svc := egreedy.NewEpsilonGreedyWithRND(1.0, statDAO, slotDAO, rnd)

	seen := make(map[int64]bool)
	const tries = 100
	for i := 0; i < tries; i++ {
		id, err := svc.Select(ctx /* slotID */, 100 /* groupID */, 200)
		assert.NoError(t, err)
		assert.Contains(t, banners, id)
		seen[id] = true
	}

	// При explore (ε=1) из 100 попыток должны попасть хотя бы в 2 разных баннера
	// при условии, что за 100 раз мы случайно не выбрали один и тот же баннер
	assert.GreaterOrEqual(t, len(seen), 2)

	// И каждый раз должно инкрементироваться view
	assert.Len(t, statDAO.viewCalls, tries)
}

//nolint:gosec
func TestSelect_Exploit(t *testing.T) {
	ctx := context.Background()
	banners := []int64{10, 20, 30}
	slotDAO := &fakeSlotDAO{banners: banners}

	// Подготавливаем искусственные stats: у баннера 20 самый высокий CTR = 50/101
	stats := map[[3]int64]*model.BannerStat{
		{1, 10, 2}: {BannerID: 10, SlotID: 1, UserGroupID: 2, Impressions: 100, Clicks: 10},
		{1, 20, 2}: {BannerID: 20, SlotID: 1, UserGroupID: 2, Impressions: 100, Clicks: 50},
		{1, 30, 2}: {BannerID: 30, SlotID: 1, UserGroupID: 2, Impressions: 100, Clicks: 20},
	}
	statDAO := &fakeStatDAO{stats: stats}

	// Любой rnd, ε=0 → всегда exploit
	svc := egreedy.NewEpsilonGreedyWithRND(0.0, statDAO, slotDAO, rand.New(rand.NewSource(17)))

	id, err := svc.Select(ctx /* slotID */, 1 /* groupID */, 2)
	assert.NoError(t, err)
	// Ожидаем баннер 20, у него наибольший CTR
	assert.Equal(t, int64(20), id)

	// И один вызов IncrementView именно для этого баннера
	if assert.Len(t, statDAO.viewCalls, 1) {
		call := statDAO.viewCalls[0]
		assert.Equal(t, int64(1), call.SlotID)
		assert.Equal(t, int64(20), call.BannerID)
		assert.Equal(t, int64(2), call.GroupID)
	}
}
