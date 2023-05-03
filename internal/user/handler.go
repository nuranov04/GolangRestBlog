package user

import (
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"go.mod/internal/apperror"
	"go.mod/internal/handlers"
	"go.mod/pkg/logging"
	"net/http"
)

const (
	usersUrl = "/users"
	userUrl  = "/users/:uuid"
)

type handler struct {
	logger  logging.Logger
	service Service
}

func NewUserHandler(logger logging.Logger, service Service) handlers.Handler {
	return &handler{
		logger:  logger,
		service: service,
	}
}

func (h handler) Register(router *httprouter.Router) {
	router.GET(usersUrl, h.GetList)
	router.GET(userUrl, h.GetUser)
	router.HandlerFunc(http.MethodPost, usersUrl, apperror.Middleware(h.CreateUser))
	router.PUT(userUrl, h.UpdateUser)
	router.PATCH(userUrl, h.PartiallyUpdateUser)
	router.DELETE(userUrl, h.DeleteUser)
}

func (h handler) GetList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userList, err := h.service.FindAll(context.TODO())
	userListBytes, err := json.Marshal(userList)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(userListBytes)
	return
}

func (h handler) CreateUser(w http.ResponseWriter, request *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	var CreateUser User
	if err := json.NewDecoder(request.Body).Decode(&CreateUser); err != nil {
		return apperror.BadRequestError("can't decode")
	}

	err := h.service.Create(context.TODO(), CreateUser)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	return nil
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
