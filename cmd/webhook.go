package main

import (
	"fmt"
	"net/http"

	"github.com/artem-shestakov/autofaq-webhook/internal/apperror"
	"github.com/artem-shestakov/autofaq-webhook/internal/handlers"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {
	errc := make(chan *apperror.Error)
	infoc := make(chan string)

	// New logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// New router
	router := mux.NewRouter()

	// Create and register handlers
	afHandler := handlers.NewAutoFAQHandler(logger, errc, infoc)
	afHandler.Register(router)

	// Create http server
	srv := http.Server{
		Addr:    ":8000",
		Handler: router,
	}
	// Logging handler
	go func() {
		for {
			select {
			case err := <-errc:
				logger.Errorf("Msg: %s. DevMsg: %s", err.Msg, err.DevMsg)
			case info := <-infoc:
				logger.Infoln(info)
			}
		}
	}()

	// Start server
	logger.Infoln(fmt.Sprintf("Server starting on %s", srv.Addr))
	err := srv.ListenAndServe()
	if err != nil {
		logger.Errorf("Error while server starting: %s", err.Error())
	}
}
