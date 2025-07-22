package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	apimw "github.com/Sucsz/banner-rotator/internal/http/middleware"
)

// NewRouter создаёт и возвращает HTTP-маршрутизатор с middleware.
func NewRouter(api *API) http.Handler {
	r := chi.NewRouter()

	// ─ Middleware ─
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(apimw.RequestLogger)

	// ─ Routes ─
	r.Route("/slots/{slot_id}", func(r chi.Router) {
		r.Post("/banners", api.AddBanner)
		r.Delete("/banners/{banner_id}", api.RemoveBanner)
		r.Post("/show", api.ShowBanner)
		r.Post("/click", api.ClickBanner)
	})

	return r
}
