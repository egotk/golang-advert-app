package corehttpresponse

import (
	"encoding/json"
	"errors"
	"net/http"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	"go.uber.org/zap"
)

type ResponseHandler struct {
	log *corezaplogger.Logger
	rw  http.ResponseWriter
}

func NewResponseHandler(
	log *corezaplogger.Logger,
	rw http.ResponseWriter,
) *ResponseHandler {
	return &ResponseHandler{
		log: log,
		rw:  rw,
	}
}

func (h *ResponseHandler) JSONResponse(
	responseBody any,
	statusCode int,
) {
	h.rw.Header().Set("Content-Type", "application/json")
	h.rw.WriteHeader(statusCode)

	if err := json.NewEncoder(h.rw).Encode(responseBody); err != nil {
		h.log.Error("write HTTP response", zap.Error(err))
	}
}

func (h *ResponseHandler) ErrorResponse(
	err error,
	msg string,
) {
	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
	)

	switch {
	case errors.Is(err, coreerrors.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn

	case errors.Is(err, coreerrors.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = h.log.Debug

	case errors.Is(err, coreerrors.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = h.log.Warn

	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}

	logFunc(msg, zap.Error(err))

	response := map[string]string{
		"message": msg,
		"error":   err.Error(),
	}

	h.JSONResponse(
		response,
		statusCode,
	)
}
