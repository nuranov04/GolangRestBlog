package main

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"go.mod/internal/api"
	"go.mod/internal/apps/post"
	db2 "go.mod/internal/apps/post/db"
	"go.mod/internal/apps/user"
	"go.mod/internal/apps/user/db"
	"go.mod/internal/config"
	"go.mod/pkg/cache/freecache"
	"go.mod/pkg/client/postgresql"
	"go.mod/pkg/jwt"
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

	router.ServeFiles("/swagger/*filepath", http.Dir("docs"))

	corsHandler := cors.Default().Handler(router)

	http.ListenAndServe("8000", corsHandler)
	postgresClient, err := postgresql.NewClient(context.Background(), 3, cfg)
	if err != nil {
		logger.Fatal(err)
	}
	userRepository := db.NewUserRepository(postgresClient, logger)

	refreshTokenCache := freecache.NewCacheRepo(104857600) // 100MB
	logger.Info("Register User api")
	jwtHelper := jwt.NewHelper(refreshTokenCache, logger)
	userService := user.NewUserService(userRepository, logger)
	userHandler := api.NewUserHandler(*logger, userService, jwtHelper)
	userHandler.Register(router)

	logger.Info("Register Post api")
	postRepository := db2.NewPostRepository(postgresClient, logger)
	postService := post.NewPostService(postRepository, logger)
	postHandler := api.NewPostHandler(*logger, postService)
	postHandler.Register(router)
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
