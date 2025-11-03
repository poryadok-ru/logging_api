package errors

import "errors"

var (
	ErrNotFound      = errors.New("ресурс не найден")
	ErrAlreadyExists = errors.New("ресурс уже существует")
	ErrUnauthorized  = errors.New("не авторизован")
	ErrForbidden     = errors.New("доступ запрещён")
)

func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}
