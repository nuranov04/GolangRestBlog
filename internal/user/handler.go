package user

import (
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"go.mod/internal/apperror"
	"go.mod/internal/handlers"
	"go.mod/pkg/logging"
	"net/http"
	"strconv"
)

const (
	usersUrl        = "/users"
	userUrlId       = "/users/id/"
	userUrlEmail    = "/users/email/"
	userUrlUsername = "/users/username/"
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
	router.HandlerFunc(http.MethodGet, userUrlId, apperror.Middleware(h.GetUserById))
	router.HandlerFunc(http.MethodGet, userUrlEmail, apperror.Middleware(h.GetUserByEmail))
	router.HandlerFunc(http.MethodGet, userUrlUsername, apperror.Middleware(h.GetUserByUsername))
	router.HandlerFunc(http.MethodPost, usersUrl, apperror.Middleware(h.CreateUser))
	router.PUT(userUrlId, h.UpdateUser)
	router.PATCH(userUrlId, h.PartiallyUpdateUser)
	router.DELETE(userUrlId, h.DeleteUser)
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
	h.logger.Info("ERROR!!!!: ", err)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}

func (h handler) GetUserById(w http.ResponseWriter, request *http.Request) error {
	userId := request.URL.Query().Get("id")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return apperror.UserIdQueryParamError
	}
	userObj, err := h.service.FindOneById(context.TODO(), userIdInt)
	if err != nil {
		return err
	}
	userObjBytes, err := json.Marshal(userObj)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	w.WriteHeader(http.StatusOK)
	w.Write(userObjBytes)
	return nil
}

func (h handler) GetUserByUsername(w http.ResponseWriter, request *http.Request) error {
	username := request.URL.Query().Get("username")
	userObj, err := h.service.FindOneByUsername(context.TODO(), username)
	//if err != nil {
	//	return err
	//}
	userObjBytes, err := json.Marshal(userObj)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	w.WriteHeader(http.StatusOK)
	w.Write(userObjBytes)
	return nil
}

func (h handler) GetUserByEmail(w http.ResponseWriter, request *http.Request) error {
	email := request.URL.Query().Get("email")
	userObj, err := h.service.FindOneByEmail(context.TODO(), email)

	userObjBytes, err := json.Marshal(userObj)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	w.WriteHeader(http.StatusOK)
	w.Write(userObjBytes)
	return nil
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
