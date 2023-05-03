package post

import (
	"github.com/julienschmidt/httprouter"
	"go.mod/internal/handlers"
	"go.mod/pkg/logging"
	"net/http"
)

const (
	postsUrl = "/posts"
	postUrl  = "/posts/:id"
)

type handler struct {
	logger logging.Logger
}

func NewPostHandler(logger logging.Logger) handlers.Handler {
	return &handler{
		logger: logger,
	}
}

func (h handler) Register(router *httprouter.Router) {
	router.GET(postsUrl, h.GetList)
	router.GET(postUrl, h.GetPost)
	router.POST(postsUrl, h.CreatePost)
	router.PUT(postUrl, h.UpdatePost)
	router.PATCH(postUrl, h.PartiallyUpdatePost)
	router.DELETE(postUrl, h.DeletePost)

}

func (h handler) GetList(w http.ResponseWriter, request *http.Request, params httprouter.Params) {
	w.WriteHeader(200)
	h.logger.Info("Get post List")
	w.Write([]byte("This is post list"))
}

func (h handler) GetPost(w http.ResponseWriter, request *http.Request, params httprouter.Params) {
	w.WriteHeader(200)
	h.logger.Info("Get post")
	w.Write([]byte("This is get single post"))
}

func (h handler) CreatePost(w http.ResponseWriter, request *http.Request, params httprouter.Params) {
	w.WriteHeader(200)
	h.logger.Info("Create post")
	w.Write([]byte("This is create post"))
}

func (h handler) UpdatePost(w http.ResponseWriter, request *http.Request, params httprouter.Params) {
	w.WriteHeader(200)
	w.Write([]byte("Post Update"))
}

func (h handler) PartiallyUpdatePost(w http.ResponseWriter, request *http.Request, params httprouter.Params) {
	w.WriteHeader(200)
	w.Write([]byte("Post Partially Update"))
}

func (h handler) DeletePost(w http.ResponseWriter, request *http.Request, params httprouter.Params) {
	w.WriteHeader(200)
	w.Write([]byte("Delete Post"))
}
