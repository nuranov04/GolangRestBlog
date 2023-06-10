package main

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"go.mod/internal/api"
	"go.mod/internal/apps/category"
	categorydb "go.mod/internal/apps/category/db"
	"go.mod/internal/apps/product"
	productdb "go.mod/internal/apps/product/db"
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

	logger.Info("Register Product api")
	productRepository := productdb.NewProductRepository(postgresClient, logger)
	productService := product.NewService(productRepository, logger)
	productHandler := api.NewPostHandler(logger, productService)
	productHandler.Register(router)

	logger.Info("Register Category api")
	categoryRepository := categorydb.NewCategoryRepository(postgresClient, logger)
	categoryService := category.NewService(categoryRepository, logger)
	categoryHandler := api.NewCategoryHandler(logger, categoryService)
	categoryHandler.Register(router)

	start(router, cfg, logger)
}

func start(router *httprouter.Router, cfg *config.Config, logger *logging.Logger) {

	logger.Info("start application")
	logger.Infof("server is listening port %s:%s", cfg.Listen.BindIp, cfg.Listen.Port)

	//log.Fatal(http.ListenAndServe(":8000", router))

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
	handler := cors.Default().Handler(router)
	server := http.Server{
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	log.Fatal(server.Serve(listener))
}
