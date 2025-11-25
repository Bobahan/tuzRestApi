package url

import (
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"github.com/go-playground/validator/v10"
)

type Request struct {
	Url   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

const aliasLength = 6

type UrlSaver interface {
	SaveUrl(urlToSave string, alias string) (int64, error) // сигнатура этого метода должна дублировать сигнатуру метода в storage
}

func New(logger *slog.Logger, urlSaver UrlSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		logger = logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		// Unmarshal запроса
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			logger.Error("failed to decode json body", sl.Err(err))
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		logger.Info("request body decoded", slog.Any("request", req))

		validate := validator.New()

		if err := validate.Struct(req); err != nil {
			logger.Error("invalid request", sl.Err(err))
			render.JSON(w, r, response.Error("invalid request"))
			return
		}

		// alias не валидируем, потому что мы даем возможность клиенту не отправлять его (по желанию), а генерировать у себя самостоятельно
		alias := req.Alias

		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveUrl(req.Url, alias)

		if errors.Is(err, storage.ErrUrlExists) {
			logger.Info("url already exists", sl.Err(err))
			render.JSON(w, r, response.Error("url already exists"))
			return
		}

		if err != nil {
			logger.Info("failed to save url", sl.Err(err))
			render.JSON(w, r, response.Error("failed to save url"))
			return
		}

		logger.Info("url saved", slog.Any("id", id))

		render.JSON(w, r, Response{
			Response: response.OK(),
			Alias:    alias,
		})
	}
}
