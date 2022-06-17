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
	MODE = "MODE"
)

func BindCommandLineArgs() {
	flag.String(PORT, "8080", "port number to listen on")

	// See https://github.com/gin-gonic/gin for modes
	flag.String(MODE, gin.DebugMode, "what mode to start webserver in")
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
	router := gin.New()
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
	gin.SetMode(viper.GetString(MODE))

	// Go routine channels
	photoChannel := make(chan models.Photos)
	shutdownChan := make(chan bool)
	quit := make(chan os.Signal)

	appService := service.NewService(photoChannel)

	// Creating webservice
	webserver, err := ListenAndServe(portNumber, appService)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to start webserver on port :%s", portNumber))

	} else {
		// Start the process that'll print the parsed photo MIME types
		appService.StartPhotoProcessorRoutine(photoChannel, shutdownChan)

		// wait for shutdown signal
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// gracefully shutdown go routines and webserver
		<-quit
		log.Println("Shutting down server...")
		shutdownChan <- true

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := webserver.Shutdown(ctx); err != nil {
			log.Fatal("Webserver forced to shutdown:", err)
		}

		log.Println("Server exiting")
	}
}
