package response

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/adriancf/demo-sistema-chamado/backend/internal/domain"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func RespondError(w http.ResponseWriter, err error) {
	var status int

	switch {
	case errors.Is(err, domain.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, domain.ErrValidation):
		status = http.StatusBadRequest
	case errors.Is(err, domain.ErrConflict):
		status = http.StatusConflict
	case errors.Is(err, domain.ErrUnauthorized):
		status = http.StatusUnauthorized
	case errors.Is(err, domain.ErrForbidden):
		status = http.StatusForbidden
	default:
		log.Printf("[ERROR] Internal server error: %v", err)
		JSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}

	JSON(w, status, ErrorResponse{Error: err.Error()})
}
