package api

import (
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"go.mod/internal"
	"go.mod/internal/apperror"
	"go.mod/internal/apps/category"
	"go.mod/pkg/logging"
	"net/http"
	"strconv"
)

const (
	categoriesUrl    = "/categories/"
	categoryIdUrl    = "/categories/id"
	categoryTitleUrl = "/categories/title"
)

type categoryHandler struct {
	logger  *logging.Logger
	service category.Service
}

func NewCategoryHandler(logger *logging.Logger, s category.Service) internal.Handler {
	return &categoryHandler{
		logger:  logger,
		service: s,
	}
}

func (h categoryHandler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, categoriesUrl, apperror.Middleware(h.Create))
	router.HandlerFunc(http.MethodGet, categoriesUrl, apperror.Middleware(h.GetList))
	router.HandlerFunc(http.MethodGet, categoryIdUrl, apperror.Middleware(h.GetOneById))
	router.HandlerFunc(http.MethodGet, categoryTitleUrl, apperror.Middleware(h.GetOneByTitle))

}

func (h categoryHandler) GetList(writer http.ResponseWriter, request *http.Request) error {
	all, err := h.service.FindAll(context.TODO())
	if err != nil {
		return err
	}
	allBytes, err := json.Marshal(all)
	writer.Write(allBytes)
	writer.WriteHeader(http.StatusOK)
	return nil
}

func (h categoryHandler) GetOneById(writer http.ResponseWriter, request *http.Request) error {
	id := request.URL.Query().Get("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return apperror.IdQueryParamError
	}
	categoryObj, err := h.service.FindOneById(context.TODO(), idInt)
	if err != nil {
		return err
	}
	categoryObjBytes, err := json.Marshal(categoryObj)
	if err != nil {
		return err
	}
	writer.Write(categoryObjBytes)
	writer.WriteHeader(http.StatusOK)
	return nil
}

func (h categoryHandler) GetOneByTitle(writer http.ResponseWriter, request *http.Request) error {
	title := request.URL.Query().Get("title")
	byTitle, err := h.service.FindOneByTitle(context.TODO(), title)
	if err != nil {
		return err
	}
	byTitleBytes, err := json.Marshal(byTitle)
	if err != nil {
		return err
	}
	writer.Write(byTitleBytes)
	writer.WriteHeader(http.StatusOK)
	return nil
}

func (h categoryHandler) Create(writer http.ResponseWriter, request *http.Request) error {
	writer.Header().Set("Content-Type", "application/json")
	var categoryDTO category.CreateUpdateCategory
	if err := json.NewDecoder(request.Body).Decode(&categoryDTO); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return apperror.BadRequestError("can't decode")
	}
	categoryObj, err := h.service.Create(context.TODO(), categoryDTO)
	if err != nil {
		return err
	}
	categoryObjBytes, err := json.Marshal(categoryObj)
	if err != nil {
		return err
	}
	writer.Write(categoryObjBytes)
	writer.WriteHeader(http.StatusCreated)
	return nil
}
