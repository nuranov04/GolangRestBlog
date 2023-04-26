package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go.mod/internal/user"
	"go.mod/pkg/logging"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("Server is starting")

	router := httprouter.New()

	userHandler := user.NewUserHandler(*logger)
	userHandler.Register(router)

	start(router)

	fmt.Println("Server is started")
}

func start(router *httprouter.Router) {
	logger := logging.GetLogger()

	logger.Info("create listener")
	listener, err := net.Listen("tcp", ":1234")

	if err != nil {
		panic(err)
	}

	logger.Info("Create server")
	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	logger.Info("server is listening port 1234")
	log.Fatal(server.Serve(listener))

}
