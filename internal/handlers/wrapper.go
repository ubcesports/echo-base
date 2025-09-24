package handlers

import (
	"errors"
	"log/slog"
	"net/http"
)

type HTTPHandlerWithErr func(http.ResponseWriter, *http.Request) error

func Wrap(h HTTPHandlerWithErr) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			var httpErr *HTTPError
			if errors.As(err, &httpErr) {
				http.Error(w, httpErr.Message, httpErr.Status)
				slog.Debug("HTTP error occurred", "status", httpErr.Status, "message", httpErr.Message, "err", err)
			} else {
				// If stuff went really wrong
				http.Error(w, "internal server error", http.StatusInternalServerError)
				slog.Error("Internal server error", "err", err)
			}
		}
	}
}
