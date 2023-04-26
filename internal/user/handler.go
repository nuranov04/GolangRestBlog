package user

import (
	"github.com/julienschmidt/httprouter"
	"go.mod/internal/handlers"
	"go.mod/pkg/logging"
	"net/http"
)

const (
	usersUrl = "/users"
	userUrl  = "/users/:uuid"
)

type handler struct {
	logger logging.Logger
}

func NewUserHandler(logger logging.Logger) handlers.Handler {
	return &handler{
		logger: logger,
	}
}

func (h handler) Register(router *httprouter.Router) {
	router.GET(usersUrl, h.GetList)
	router.GET(userUrl, h.GetUser)
	router.POST(usersUrl, h.CreateUser)
	router.PUT(userUrl, h.UpdateUser)
	router.PATCH(userUrl, h.PartiallyUpdateUser)
	router.DELETE(userUrl, h.DeleteUser)
}

func (h handler) GetList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(200)
	h.logger.Info("Get User List")
	w.Write([]byte("This is user list"))

}

func (h handler) CreateUser(w http.ResponseWriter, request *http.Request, params httprouter.Params) {
	w.WriteHeader(200)
	w.Write([]byte("User Created"))
}

func (h handler) GetUser(w http.ResponseWriter, request *http.Request, params httprouter.Params) {
	w.WriteHeader(200)
	w.Write([]byte("Get user by uuid"))
}

func (h handler) UpdateUser(w http.ResponseWriter, request *http.Request, params httprouter.Params) {
	w.WriteHeader(200)
	w.Write([]byte("User Update"))
}

func (h handler) PartiallyUpdateUser(w http.ResponseWriter, request *http.Request, params httprouter.Params) {
	w.WriteHeader(200)
	w.Write([]byte("User Partially Update"))

}

func (h handler) DeleteUser(w http.ResponseWriter, request *http.Request, params httprouter.Params) {
	w.WriteHeader(200)
	w.Write([]byte("Delete User"))
}
