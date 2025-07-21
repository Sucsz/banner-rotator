package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/Sucsz/banner-rotator/internal/kafka"
	"github.com/Sucsz/banner-rotator/internal/log"
)

// AddBanner — POST /slots/{slot_id}/banners
func (a *API) AddBanner(w http.ResponseWriter, r *http.Request) {
	logger := log.WithComponent("api.AddBanner")

	slotID, err := strconv.ParseInt(chi.URLParam(r, "slot_id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid slot_id", http.StatusBadRequest)
		return
	}

	var body struct {
		BannerID int64 `json:"banner_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.BannerSlotDAO.AddBannerToSlot(r.Context(), body.BannerID, slotID); err != nil {
		logger.Error().Err(err).Msg("AddBannerToSlot failed.")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RemoveBanner — DELETE /slots/{slot_id}/banners/{banner_id}
func (a *API) RemoveBanner(w http.ResponseWriter, r *http.Request) {
	logger := log.WithComponent("api.RemoveBanner")

	slotID, err := strconv.ParseInt(chi.URLParam(r, "slot_id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid slot_id", http.StatusBadRequest)
		return
	}
	bannerID, err := strconv.ParseInt(chi.URLParam(r, "banner_id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid banner_id", http.StatusBadRequest)
		return
	}

	if err := a.BannerSlotDAO.RemoveBannerFromSlot(r.Context(), bannerID, slotID); err != nil {
		logger.Error().Err(err).Msg("RemoveBannerFromSlot failed.")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ShowBanner — POST /slots/{slot_id}/show
func (a *API) ShowBanner(w http.ResponseWriter, r *http.Request) {
	logger := log.WithComponent("api.ShowBanner")

	slotID, err := strconv.ParseInt(chi.URLParam(r, "slot_id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid slot_id", http.StatusBadRequest)
		return
	}

	var body struct {
		GroupID int64 `json:"group_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 1) Выбрать баннер
	bannerID, err := a.Selector.Select(r.Context(), slotID, body.GroupID)
	if err != nil {
		logger.Error().Err(err).Msg("selector.Select failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// 2) Отправить событие показа
	event := kafka.BannerEvent{
		Type:        "impression",
		SlotID:      slotID,
		BannerID:    bannerID,
		UserGroupID: body.GroupID,
		Timestamp:   time.Now(),
	}
	if err := a.Producer.Send(r.Context(), event); err != nil {
		logger.Error().Err(err).Msg("producer.Send impression failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// 3) Ответ клиенту
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]int64{"banner_id": bannerID}); err != nil {
		logger.Error().Err(err).Msg("json.Encode failed")
	}
}

// ClickBanner — POST /slots/{slot_id}/click
func (a *API) ClickBanner(w http.ResponseWriter, r *http.Request) {
	logger := log.WithComponent("api.ClickBanner")

	slotID, err := strconv.ParseInt(chi.URLParam(r, "slot_id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid slot_id", http.StatusBadRequest)
		return
	}

	var body struct {
		BannerID int64 `json:"banner_id"`
		GroupID  int64 `json:"group_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 1) Засчитать клик
	if err := a.Selector.RecordClick(r.Context(), slotID, body.BannerID, body.GroupID); err != nil {
		logger.Error().Err(err).Msg("selector.RecordClick failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// 2) Отправить событие клика
	event := kafka.BannerEvent{
		Type:        "click",
		SlotID:      slotID,
		BannerID:    body.BannerID,
		UserGroupID: body.GroupID,
		Timestamp:   time.Now(),
	}
	if err := a.Producer.Send(r.Context(), event); err != nil {
		logger.Error().Err(err).Msg("producer.Send click failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
