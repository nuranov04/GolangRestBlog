package post

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
	postsUrl = "/posts"
	postUrl  = "/posts/:id"
)

type handler struct {
	logger  logging.Logger
	service service
}

func NewPostHandler(logger logging.Logger, s service) handlers.Handler {
	return &handler{
		logger:  logger,
		service: s,
	}
}

func (h handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, postsUrl, apperror.Middleware(h.GetList))
	router.HandlerFunc(http.MethodGet, postUrl, apperror.Middleware(h.GetPost))
	router.HandlerFunc(http.MethodPost, postsUrl, apperror.Middleware(h.CreatePost))
	router.HandlerFunc(http.MethodPut, postUrl, apperror.Middleware(h.UpdatePost))
	router.HandlerFunc(http.MethodDelete, postUrl, apperror.Middleware(h.DeletePost))
}

func (h handler) GetList(w http.ResponseWriter, request *http.Request) error {
	posts, err := h.service.FindAll(context.TODO())
	if err != nil {
		return err
	}
	postsBytes, err := json.Marshal(posts)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	w.WriteHeader(http.StatusOK)
	w.Write(postsBytes)
	return nil
}

func (h handler) GetPost(w http.ResponseWriter, request *http.Request) error {
	id := request.URL.Query().Get("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	post, err := h.service.FindOneById(context.TODO(), idInt)
	if err != nil {
		return err
	}
	postBytes, err := json.Marshal(post)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	w.Write(postBytes)
	return nil
}

func (h handler) CreatePost(w http.ResponseWriter, request *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	var CreatePostDTO CreatePostDTO
	if err := json.NewDecoder(request.Body).Decode(&CreatePostDTO); err != nil {
		return apperror.BadRequestError("can't decode")
	}
	createdPost, err := h.service.Create(context.TODO(), CreatePostDTO)
	if err != nil {
		return err
	}
	createdPostBytes, err := json.Marshal(createdPost)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(createdPostBytes)
	return nil
}

func (h handler) UpdatePost(w http.ResponseWriter, request *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	postId := request.URL.Query().Get("id")
	postIdInt, err := strconv.Atoi(postId)
	if err != nil {
		return apperror.IdQueryParamError
	}
	postObj, err := h.service.FindOneById(context.TODO(), postIdInt)

	var updatePost UpdatePostDTO
	if err := json.NewDecoder(request.Body).Decode(&updatePost); err != nil {
		return apperror.BadRequestError("can't decode")
	}
	updatedPostObj, err := h.service.Update(context.TODO(), postObj, updatePost)
	if err != nil {
		return apperror.BadRequestError("can't update user")
	}
	updatedPostObjBytes, err := json.Marshal(updatedPostObj)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	w.WriteHeader(http.StatusOK)
	w.Write(updatedPostObjBytes)
	return nil
}

func (h handler) DeletePost(w http.ResponseWriter, request *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	postId := request.URL.Query().Get("id")
	postIdInt, err := strconv.Atoi(postId)
	if err != nil {
		return apperror.IdQueryParamError
	}
	err = h.service.Delete(context.TODO(), postIdInt)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}
