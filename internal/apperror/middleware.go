package apperror

import (
	"errors"
	"net/http"
)

type appHandler func(http.ResponseWriter, *http.Request) error

func Middleware(h appHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var appErr *AppError
		err := h(w, r)
		if err != nil {
			if errors.As(err, &appErr) {
				if errors.Is(err, ErrorNotFound) {
					w.WriteHeader(http.StatusNotFound)
					w.Write(ErrorNotFound.Marshal())
					return
				}
				if errors.Is(err, UserAlreadyExist) {
					w.WriteHeader(http.StatusBadRequest)
					w.Write(UserAlreadyExist.Marshal())
					return
				}
				err := err.(*AppError)
				w.WriteHeader(http.StatusBadRequest)
				w.Write(err.Marshal())
				return
			}
			w.WriteHeader(http.StatusTeapot)
			w.Write(systemError(err.Error()).Marshal())
		}
	}
}
