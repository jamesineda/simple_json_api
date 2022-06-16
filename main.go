package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jamesineda/simple_json_api/app"
	"github.com/jamesineda/simple_json_api/app/models"
	"github.com/jamesineda/simple_json_api/app/service"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	PORT = "PORT"
)

func BindCommandLineArgs() {
	flag.String(PORT, "8080", "port number to listen on")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
}

/*
	TODO
		- Add a database/ cache storage client
		- Store requests that have been processed
		- Add the sref to the Gin logs for debug purposes
		- Add authentication
		- Go tests
*/

func ListenAndServe(portNumber string, service *service.Service) (srv *http.Server, err error) {
	router := gin.Default()
	app.CreateV1Routes(router, service)

	srv = &http.Server{
		Addr:    fmt.Sprintf(":%s", portNumber),
		Handler: router,
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			return
		}
	}()
	return
}

func main() {
	BindCommandLineArgs()
	portNumber := viper.GetString(PORT)

	// Go routine channels
	photoChannel := make(chan models.Photos)
	shutdownChan := make(chan bool)
	quit := make(chan os.Signal)

	service := service.NewService(photoChannel)

	// Creating webservice
	srv, err := ListenAndServe(portNumber, service)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to start webserver on port :%s", portNumber))

	} else {
		// Start the process that'll print the parsed photo MIME types
		service.StartPhotoProcessorRoutine(photoChannel, shutdownChan)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutting down server...")
		shutdownChan <- true

		// The context is used to inform the server it has 5 seconds to finish
		// the request it is currently handling
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server forced to shutdown:", err)
		}

		log.Println("Server exiting")
	}
}
