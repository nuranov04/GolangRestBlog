package main

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go.mod/internal/config"
	"go.mod/internal/post"
	db2 "go.mod/internal/post/db"
	"go.mod/internal/user"
	"go.mod/internal/user/db"
	"go.mod/pkg/client/postgresql"
	"go.mod/pkg/logging"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

func main() {

	logger := logging.GetLogger()
	logger.Info("Server is starting")

	router := httprouter.New()

	cfg := config.GetConfig()
	logger.Info("read application configs")

	postgresClient, err := postgresql.NewClient(context.Background(), 3, cfg.Storage)
	if err != nil {
		logger.Fatal(err)
	}
	userRepository := db.NewRepository(postgresClient, logger)

	logger.Info("Register User handler")
	userService := user.NewService(userRepository, logger)
	userHandler := user.NewUserHandler(*logger, *userService)
	userHandler.Register(router)

	logger.Info("Register Post handler")
	postRepository := db2.NewRepository(postgresClient, logger)
	postService := post.NewService(postRepository, logger)
	postHandler := post.NewPostHandler(*logger, postService)

	start(router, cfg, logger)
	fmt.Println("Server is started")
}

func start(router *httprouter.Router, cfg *config.Config, logger *logging.Logger) {

	logger.Info("start application")

	var listener net.Listener
	var ListenError error

	if cfg.Listen.Type == "sock" {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		logger.Info(appDir)
		if err != nil {
			panic(err)
		}
		logger.Info("create socket")

		socketPath := path.Join(appDir, "app.sock")

		listener, ListenError = net.Listen("unix", socketPath)

		logger.Infof("Server is listening unix socket: %s", socketPath)
	} else {
		logger.Info("listen tcp")
		listener, ListenError = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIp, cfg.Listen.Port))
		logger.Infof("server is listening port %s:%s", cfg.Listen.BindIp, cfg.Listen.Port)
	}

	if ListenError != nil {
		panic(ListenError)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(server.Serve(listener))

}
