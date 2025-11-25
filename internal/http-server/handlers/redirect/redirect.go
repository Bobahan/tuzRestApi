package redirect

import (
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type UrlGetter interface {
	GetUrl(alias string) (string, error) // сигнатура этого метода должна дублировать сигнатуру метода в storage
}

func New(logger *slog.Logger, urlGetter UrlGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.New"

		logger = logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			logger.Info("alias is empty")
			render.JSON(w, r, response.Error("invalid request"))
			return
		}

		url, err := urlGetter.GetUrl(alias)

		if errors.Is(err, storage.ErrUrlNotFound) {
			logger.Info("url not found", "alias", alias)
			render.JSON(w, r, response.Error("not found"))
			return
		}

		if err != nil {
			logger.Error("failed to get url", sl.Err(err))
			render.JSON(w, r, response.Error("internal error"))
			return
		}

		logger.Info("got url", slog.String("url", url))

		http.Redirect(w, r, url, http.StatusFound)
	}
}
