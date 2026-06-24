package corehttpresponse

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"path/filepath"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	"go.uber.org/zap"
)

type Handler struct {
	log *corezaplogger.Logger
	rw  http.ResponseWriter
}

func New(
	log *corezaplogger.Logger,
	rw http.ResponseWriter,
) *Handler {
	return &Handler{
		log: log,
		rw:  rw,
	}
}

func From(
	ctx context.Context,
	rw http.ResponseWriter,
) *Handler {
	log := corezaplogger.FromContext(ctx)

	return &Handler{
		log: log,
		rw:  rw,
	}
}

func (h *Handler) JSONResponse(
	responseBody any,
	statusCode int,
) {
	h.rw.Header().Set("Content-Type", "application/json")
	h.rw.WriteHeader(statusCode)

	if err := json.NewEncoder(h.rw).Encode(responseBody); err != nil {
		h.log.Error("write HTTP response", zap.Error(err))
	}
}

func (h *Handler) NoContentResponse() {
	h.rw.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ErrorResponse(
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

	case errors.Is(err, coreerrors.ErrUnauthorized):
		statusCode = http.StatusUnauthorized
		logFunc = h.log.Warn

	case errors.Is(err, coreerrors.ErrForbidden):
		statusCode = http.StatusForbidden
		logFunc = h.log.Debug

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

func (h *Handler) FileResponse(
	path string,
	rc io.ReadCloser,
) {
	h.rw.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(path)))

	if _, err := io.Copy(h.rw, rc); err != nil {
		h.log.Error("copy to RW", zap.Error(err))
	}
}
