package api

import (
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"go.mod/internal"
	"go.mod/internal/apperror"
	"go.mod/internal/apps/product"
	"go.mod/pkg/logging"
	"net/http"
	"strconv"
)

const (
	postsUrl = "/products/"
	postUrl  = "/products/id/"
)

type postHandler struct {
	logger  *logging.Logger
	service product.Service
}

func NewPostHandler(logger *logging.Logger, s product.Service) internal.Handler {
	return &postHandler{
		logger:  logger,
		service: s,
	}
}

func (h postHandler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, postsUrl, apperror.Middleware(h.GetList))
	router.HandlerFunc(http.MethodGet, postUrl, apperror.Middleware(h.Get))
	router.HandlerFunc(http.MethodPost, postsUrl, apperror.Middleware(h.Create))
	router.HandlerFunc(http.MethodPut, postUrl, apperror.Middleware(h.Update))
	router.HandlerFunc(http.MethodDelete, postUrl, apperror.Middleware(h.Delete))
}

func (h postHandler) GetList(w http.ResponseWriter, request *http.Request) error {
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

func (h postHandler) Get(w http.ResponseWriter, request *http.Request) error {
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

func (h postHandler) Create(w http.ResponseWriter, request *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	var CreatePostDTO product.CreateProductDTO
	if err := json.NewDecoder(request.Body).Decode(&CreatePostDTO); err != nil {
		return apperror.BadRequestError("can't decode")
	}
	createdPost, err := h.service.Create(context.TODO(), CreatePostDTO)
	if err != nil {
		return err
	}
	createdPostBytes, err := json.Marshal(*createdPost)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(createdPostBytes)
	return nil
}

func (h postHandler) Update(w http.ResponseWriter, request *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	postId := request.URL.Query().Get("id")
	postIdInt, err := strconv.Atoi(postId)
	if err != nil {
		return apperror.IdQueryParamError
	}
	postObj, err := h.service.FindOneById(context.TODO(), postIdInt)
	if err != nil {
		return err
	}
	var updatePost product.UpdateProductDTO
	if err := json.NewDecoder(request.Body).Decode(&updatePost); err != nil {
		h.logger.Debug(err)
		return apperror.BadRequestError("can't decode")
	}
	updatedPostObj, err := h.service.Update(context.TODO(), postObj, updatePost)
	if err != nil {
		return err
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

func (h postHandler) Delete(w http.ResponseWriter, request *http.Request) error {
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
	w.Write([]byte("status: Deleted"))
	w.WriteHeader(http.StatusOK)
	return nil
}
