package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/lib/random"
	"url-shortener/storage"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=URLSaver
type URLSaver interface {
	SaveURL(url string, alias string) error
}

//go:generate go run github.com/vektra/mockery/v3@latest --name=URLGetter
type URLGetter interface {
	GetURL(url string) (string, error)
}

type Request struct {
	URL   string `json:"url" validate:"required,url"` //validation here!!! on structure type
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

const aliasLength = 10

func New(log *slog.Logger, urlSaver URLSaver, urlGetter URLGetter) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handler.save.url.New"

		log = log.With(slog.String("op", op),
			slog.String("request_id",
				middleware.GetReqID(request.Context())),
		)

		var req Request

		err := render.DecodeJSON(request.Body, &req)
		if err != nil {
			log.Error("failed decode body", sl.Err(err))

			render.JSON(writer, request, response.Error("failed to decode message"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(writer, request, response.ValidationError(validateErr))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)

			_, err := urlGetter.GetURL(alias)
			for err == nil {
				_, err = urlGetter.GetURL(alias)
				alias = random.NewRandomString(aliasLength)
			}
		}

		err = urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(writer, request, response.Error("url already exists"))

			return
		}

		if err != nil {
			log.Error("failed save url", sl.Err(err))

			render.JSON(writer, request, response.Error("failed save url"))

			return
		}

		log.Info("url saved", slog.String("url", req.URL))

		responseOK(writer, request, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: response.OK(),
		Alias:    alias,
	})
}
