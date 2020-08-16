package sqlpad

import (
	"fmt"
	"net/http"
)

// Error расширенная ошибка
type Error struct {
	StatusCode int
	Response   *http.Response
	Body       []byte
}

func (e *Error) Error() string {
	return fmt.Sprintf("[sqlpad_client] request failed with code: %d", e.StatusCode)
}

// ErrorBody возвращает тело ошибки
func ErrorBody(err error) string {
	if err == nil {
		return ""
	}
	if e, ok := err.(*Error); ok {
		return string(e.Body)
	}
	return err.Error()
}
