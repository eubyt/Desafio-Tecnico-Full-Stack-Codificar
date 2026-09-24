package response_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adriancf/demo-sistema-chamado/backend/internal/domain"
	httpResponse "github.com/adriancf/demo-sistema-chamado/backend/internal/infra/http/response"
	"github.com/stretchr/testify/assert"
)

func TestRespondError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "should map ErrNotFound to 404",
			err:            domain.NewError("item not found", domain.ErrNotFound),
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"item not found"}`,
		},
		{
			name:           "should map ErrValidation to 400",
			err:            domain.NewError("invalid title length", domain.ErrValidation),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid title length"}`,
		},
		{
			name:           "should map ErrConflict to 409",
			err:            domain.NewError("already exists", domain.ErrConflict),
			expectedStatus: http.StatusConflict,
			expectedBody:   `{"error":"already exists"}`,
		},
		{
			name:           "should map ErrUnauthorized to 401",
			err:            domain.NewError("missing token", domain.ErrUnauthorized),
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"missing token"}`,
		},
		{
			name:           "should map ErrForbidden to 403",
			err:            domain.NewError("permission denied", domain.ErrForbidden),
			expectedStatus: http.StatusForbidden,
			expectedBody:   `{"error":"permission denied"}`,
		},
		{
			name:           "should map unknown/internal error to 500 without leaking details",
			err:            errors.New("raw postgres connection timeout"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			httpResponse.RespondError(rec, tt.err)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
		})
	}
}
