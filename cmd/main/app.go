package main

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go.mod/internal/config"
	"go.mod/internal/user"
	"go.mod/internal/user/db"
	"go.mod/pkg/client/mongodb"
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
	cfgMongo := cfg.MongoDb
	mongoDBClient, err := mongodb.NewClient(
		context.Background(),
		cfgMongo.Host,
		cfgMongo.Port,
		cfgMongo.Username,
		cfgMongo.Password,
		cfgMongo.Database,
		cfgMongo.AuthDB,
	)
	if err != nil {
		panic(err)
	}
	storage := db.NewStorage(mongoDBClient, cfgMongo.Collection, logger)

	u := user.User{
		ID:           "",
		Username:     "admin",
		PasswordHash: "admin",
		Email:        "admin@gmail.com",
	}
	user1, err := storage.Create(context.Background(), u)
	if err != nil {
		panic(err)
	}
	logger.Info(user1)

	logger.Info("Register User handler")
	userHandler := user.NewUserHandler(*logger)
	userHandler.Register(router)

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
