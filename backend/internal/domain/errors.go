package domain

import "errors"

var (
	// ErrNotFound indica que o recurso solicitado não existe (HTTP 404)
	ErrNotFound = errors.New("not found")

	// ErrValidation indica dados inválidos ou violação de regra de negócio (HTTP 400)
	ErrValidation = errors.New("validation failed")

	// ErrConflict indica conflito com o estado atual do recurso (HTTP 409)
	ErrConflict = errors.New("conflict")

	// ErrUnauthorized indica falta de autenticação (HTTP 401)
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden indica falta de permissão de acesso (HTTP 403)
	ErrForbidden = errors.New("forbidden")
)

type domainError struct {
	msg  string
	base error
}

func (e *domainError) Error() string {
	return e.msg
}

func (e *domainError) Unwrap() error {
	return e.base
}

func NewError(msg string, base error) error {
	return &domainError{
		msg:  msg,
		base: base,
	}
}
