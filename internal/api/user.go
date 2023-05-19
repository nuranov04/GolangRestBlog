package api

import (
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"go.mod/internal"
	"go.mod/internal/apperror"
	"go.mod/internal/apps/user"
	"go.mod/pkg/jwt"
	"go.mod/pkg/logging"
	"net/http"
	"strconv"
)

const (
	usersUrl        = "/users"
	userUrlId       = "/users/id/"
	userUrlEmail    = "/users/email/"
	userUrlUsername = "/users/username/"
	loginUrl        = "/users/login/"
)

type userHandler struct {
	logger    logging.Logger
	service   user.Service
	JWTHelper jwt.Helper
}

func NewUserHandler(logger logging.Logger, service user.Service, jwtHelper jwt.Helper) internal.Handler {
	return &userHandler{
		logger:    logger,
		service:   service,
		JWTHelper: jwtHelper,
	}
}

func (h userHandler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, usersUrl, apperror.Middleware(h.GetList))
	router.HandlerFunc(http.MethodGet, userUrlId, apperror.Middleware(h.GetUserById))
	router.HandlerFunc(http.MethodGet, userUrlEmail, apperror.Middleware(h.GetUserByEmail))
	router.HandlerFunc(http.MethodGet, userUrlUsername, apperror.Middleware(h.GetUserByUsername))
	router.HandlerFunc(http.MethodPost, usersUrl, apperror.Middleware(h.CreateUser))
	router.HandlerFunc(http.MethodPut, userUrlId, apperror.Middleware(h.UpdateUser))
	router.HandlerFunc(http.MethodDelete, userUrlId, apperror.Middleware(h.DeleteUser))
	router.HandlerFunc(http.MethodPost, loginUrl, apperror.Middleware(h.Login))
	router.HandlerFunc(http.MethodPut, loginUrl, apperror.Middleware(h.Login))
}

func (h userHandler) Login(w http.ResponseWriter, request *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	var (
		token []byte
		err   error
	)
	switch request.Method {
	case http.MethodPost:
		username := request.URL.Query().Get("username")
		password := request.URL.Query().Get("password")
		if username == "" || password == "" {
			w.WriteHeader(http.StatusBadRequest)
			return apperror.BadRequestError("invalid query parameters email or password")
		}
		u, err := h.service.FindUserByUsernameAndPassword(context.TODO(), username, password)
		if err != nil {
			return err
		}
		token, err = h.JWTHelper.GenerateAccessToken(u)
		if err != nil {
			return err
		}
	case http.MethodPut:
		var rt jwt.RT
		if err := json.NewDecoder(request.Body).Decode(&rt); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return apperror.BadRequestError("failed to decode data")
		}
		token, err = h.JWTHelper.UpdateRefreshToken(rt)
		if err != nil {
			return err
		}

	}
	w.WriteHeader(http.StatusCreated)
	w.Write(token)
	return nil
}

func (h userHandler) GetList(w http.ResponseWriter, request *http.Request) error {
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

func (h userHandler) CreateUser(w http.ResponseWriter, request *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	var CreateUser user.CreateUserDTO
	if err := json.NewDecoder(request.Body).Decode(&CreateUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return apperror.BadRequestError("can't decode")
	}

	userObj, err := h.service.Create(context.TODO(), CreateUser)
	h.logger.Info("error after create\n", err)
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

func (h userHandler) UpdateUser(w http.ResponseWriter, request *http.Request) error {
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
	var updateUser user.UpdateUserDTO
	if err := json.NewDecoder(request.Body).Decode(&updateUser); err != nil {
		return apperror.BadRequestError("can't decode")
	}
	updatedUserObj, err := h.service.UserUpdate(context.TODO(), *userObj, updateUser)
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

func (h userHandler) GetUserById(w http.ResponseWriter, request *http.Request) error {
	userId := request.URL.Query().Get("id")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return apperror.IdQueryParamError
	}
	userObj, err := h.service.FindOneById(context.TODO(), userIdInt)
	if err != nil {
		return apperror.ErrorNotFound
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

func (h userHandler) GetUserByUsername(w http.ResponseWriter, request *http.Request) error {
	username := request.URL.Query().Get("username")
	userObj, err := h.service.FindOneByUsername(context.TODO(), username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return apperror.ErrorNotFound
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

func (h userHandler) GetUserByEmail(w http.ResponseWriter, request *http.Request) error {
	email := request.URL.Query().Get("email")
	userObj, err := h.service.FindOneByEmail(context.TODO(), email)
	if err != nil {
		return apperror.ErrorNotFound
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

func (h userHandler) DeleteUser(w http.ResponseWriter, request *http.Request) error {
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
