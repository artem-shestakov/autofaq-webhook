package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/artem-shestakov/autofaq-webhook/internal/apperror"
	"github.com/artem-shestakov/autofaq-webhook/internal/config"
	"github.com/artem-shestakov/autofaq-webhook/internal/handlers"
	"github.com/artem-shestakov/autofaq-webhook/internal/server"

	"github.com/sirupsen/logrus"
)

var confPath = flag.String("config", "./", "Path to the configuration file")

func main() {
	flag.Parse()
	ctx := context.Background()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	// Logging channels
	errc := make(chan *apperror.Error)
	infoc := make(chan string)
	warnc := make(chan string)

	// New logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Logging handler
	go func() {
		for {
			select {
			case err := <-errc:
				logger.Errorf("Msg: %s. DevMsg: %s", err.Msg, err.DevMsg)
			case warn := <-warnc:
				logger.Warnln(warn)
			case info := <-infoc:
				logger.Infoln(info)
			}
		}
	}()

	conf := config.LoadConfig(*confPath, errc, warnc)

	// Create http server
	srv, router := server.NewServer(conf.Server.Address+":"+conf.Server.Port, errc, infoc)

	// Create and register handlers
	afHandler := handlers.NewAutoFAQHandler(logger, errc, infoc)
	afHandler.Register(router)

	// Start server
	go srv.Run()

	// Stop server
	signal := <-stop
	logger.Infoln(fmt.Sprintf("Server is stoping. Get signal %v", signal))
	srv.Stop(ctx)
}
