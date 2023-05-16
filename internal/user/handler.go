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
	loginUrl        = "/login/"
)

type handler struct {
	logger  logging.Logger
	service service
}

func NewUserHandler(logger logging.Logger, service service) handlers.Handler {
	return &handler{
		logger:  logger,
		service: service,
	}
}

func (h handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, usersUrl, apperror.Middleware(h.GetList))
	router.HandlerFunc(http.MethodGet, userUrlId, apperror.Middleware(h.GetUserById))
	router.HandlerFunc(http.MethodGet, userUrlEmail, apperror.Middleware(h.GetUserByEmail))
	router.HandlerFunc(http.MethodGet, userUrlUsername, apperror.Middleware(h.GetUserByUsername))
	router.HandlerFunc(http.MethodPost, usersUrl, apperror.Middleware(h.CreateUser))
	router.HandlerFunc(http.MethodPut, userUrlId, apperror.Middleware(h.UpdateUser))
	router.HandlerFunc(http.MethodDelete, userUrlId, apperror.Middleware(h.DeleteUser))
	router.HandlerFunc(http.MethodPost, loginUrl, apperror.Middleware(h.Login))
}

func (h handler) GetList(w http.ResponseWriter, request *http.Request) error {
	userList, err := h.service.FindAll(context.TODO())
	userListBytes, err := json.Marshal(userList)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	w.WriteHeader(http.StatusOK)
	w.Write(userListBytes)
	return nil
}

func (h handler) CreateUser(w http.ResponseWriter, request *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	var CreateUser CreateUserDTO
	if err := json.NewDecoder(request.Body).Decode(&CreateUser); err != nil {
		return apperror.BadRequestError("can't decode")
	}

	userObj, err := h.service.Create(context.TODO(), CreateUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	userObjBytes, err := json.Marshal(*userObj)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	w.Write(userObjBytes)
	w.WriteHeader(http.StatusCreated)
	return nil
}

func (h handler) Login(w http.ResponseWriter, request *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	var userLoginDTO LoginDTO
	if err := json.NewDecoder(request.Body).Decode(&userLoginDTO); err != nil {
		return apperror.BadRequestError("can't decode")
	}
	return nil
}

func (h handler) UpdateUser(w http.ResponseWriter, request *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	userId := request.URL.Query().Get("id")
	UserIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return apperror.IdQueryParamError
	}
	userObj, err := h.service.FindOneById(context.TODO(), UserIdInt)
	if err != nil {
		return apperror.ErrorNotFound
	}
	var updateUser UpdateUserDTO
	if err := json.NewDecoder(request.Body).Decode(&updateUser); err != nil {
		return apperror.BadRequestError("can't decode")
	}
	updatedUserObj, err := h.service.userUpdate(context.TODO(), *userObj, updateUser)
	if err != nil {
		return err
	}
	updatedUserObjBytes, err := json.Marshal(updatedUserObj)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	w.WriteHeader(http.StatusOK)
	w.Write(updatedUserObjBytes)
	return nil
}

func (h handler) GetUserById(w http.ResponseWriter, request *http.Request) error {
	userId := request.URL.Query().Get("id")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return apperror.IdQueryParamError
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

func (h handler) GetUserByEmail(w http.ResponseWriter, request *http.Request) error {
	email := request.URL.Query().Get("email")
	userObj, err := h.service.FindOneByEmail(context.TODO(), email)
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

func (h handler) DeleteUser(w http.ResponseWriter, request *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	userId := request.URL.Query().Get("id")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return err
	}
	err = h.service.Delete(context.TODO(), userIdInt)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}
