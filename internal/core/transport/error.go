package transport

import (
	coreerrors "Miji/internal/core/errors"
	"errors"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

type ErrorResponse struct {
	Message    string `json:"message"`
	StatusCode string `json:"statusCode"`
	Error      string `json:"error"`
}

type ValidationErrorResponse struct {
	Message  string            `json:"message"`
	Problems map[string]string `json:"problems"`
}

var ErrValidation = errors.New("validation error")

func Error(w http.ResponseWriter, _ *http.Request, log *zap.Logger, err error, msg string) {
	statusCode := statusFromError(err)

	log = log.With(zap.Error(err))
	switch {
	case errors.Is(err, coreerrors.ErrNotFound):
		log.Debug(msg)
	case errors.Is(err, coreerrors.ErrInvalidArgument), errors.Is(err, coreerrors.ErrConflict):
		log.Warn(msg)
	default:
		log.Error(msg)
	}

	_ = Encode(w, nil, statusCode, ErrorResponse{
		Message:    msg,
		StatusCode: strconv.Itoa(statusCode),
		Error:      err.Error(),
	})
}

func ValidationError(w http.ResponseWriter, _ *http.Request, log *zap.Logger, problems map[string]string) {
	log.Warn("validation failed", zap.Any("problems", problems))

	_ = Encode(w, nil, http.StatusBadRequest, ValidationErrorResponse{
		Message:  "validation failed",
		Problems: problems,
	})
}

func statusFromError(err error) int {
	switch {
	case errors.Is(err, coreerrors.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, coreerrors.ErrInvalidArgument):
		return http.StatusBadRequest
	case errors.Is(err, coreerrors.ErrConflict):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
