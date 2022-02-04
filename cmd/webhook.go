package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/artem-shestakov/autofaq-webhook/internal/apperror"
	"github.com/artem-shestakov/autofaq-webhook/internal/handlers"
	"github.com/artem-shestakov/autofaq-webhook/internal/server"

	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	// Logging channels
	errc := make(chan *apperror.Error)
	infoc := make(chan string)
	// warnc := make(chan string)

	// New logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// config.LoadConfig("./", errc, infoc)

	// Create http server
	srv, router := server.NewServer(":8000", errc, infoc)

	// Create and register handlers
	afHandler := handlers.NewAutoFAQHandler(logger, errc, infoc)
	afHandler.Register(router)

	// Logging handler
	go func() {
		for {
			select {
			case err := <-errc:
				logger.Errorf("Msg: %s. DevMsg: %s", err.Msg, err.DevMsg)
			// case warn := <-warnc:
			// 	// logger.Warnln(warn)
			// 	fmt.Println(warn)
			case info := <-infoc:
				logger.Infoln(info)
			}
		}
	}()

	// Start server
	go srv.Run()

	// Stop server
	signal := <-stop
	logger.Infoln(fmt.Sprintf("Server is stoping. Get signal %v", signal))
	srv.Stop(ctx)
}
