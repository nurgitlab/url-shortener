package deleteURL

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
)

type Request struct {
	Alias string `json:"alias" validate:"required"`
}

type Response struct {
	response.Response
}

type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handler.delete.url.New"

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

		err = urlDeleter.DeleteURL(alias)

		fmt.Println("err", err)

		if err != nil {
			log.Error("failed delete url", sl.Err(err))
			render.JSON(writer, request, response.Error("failed to delete"))

			return
		}

		log.Info("url deleted", sl.Err(err))

		responseOK(writer, request)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Response: response.OK(),
	})
}
