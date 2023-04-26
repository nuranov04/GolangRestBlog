package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go.mod/internal/config"
	"go.mod/internal/user"
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
		socketPath := path.Join(appDir, "app.sock")
		logger.Debugf("socker path: %s", socketPath)
		listener, ListenError = net.Listen("unix", socketPath)
	} else {

		logger.Info("Create server")
		listener, ListenError = net.Listen("tcp", ":1234")
	}

	if ListenError != nil {
		panic(ListenError)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	logger.Infof("server is listening port %s:%s", cfg.Listen.BindIp, cfg.Listen.Port)
	log.Fatal(server.Serve(listener))

}
