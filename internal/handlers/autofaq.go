package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/artem-shestakov/autofaq-webhook/internal/apperror"
	"github.com/artem-shestakov/autofaq-webhook/internal/models"
	"github.com/artem-shestakov/autofaq-webhook/internal/prom"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Handler for handle reqiest from AutoFAQ
type afHandler struct {
	log   *logrus.Logger
	errc  chan *apperror.Error
	infoc chan string
}

// NewAutoFAQHandler create new handler
func NewAutoFAQHandler(logger *logrus.Logger, errc chan *apperror.Error, infoc chan string) *afHandler {
	return &afHandler{
		log:   logger,
		errc:  errc,
		infoc: infoc,
	}
}

// Register URLs
func (h *afHandler) Register(router *mux.Router) {
	router.HandleFunc("/webhook", h.handleRequest)
}

// Handle request from AutoFAQ
func (h *afHandler) handleRequest(rw http.ResponseWriter, r *http.Request) {
	var messages models.Messages
	// Decode response from AutoFAQ
	err := json.NewDecoder(r.Body).Decode(&messages)
	if err != nil {
		h.errc <- apperror.NewError("Can't decode response from AutoFAQ", err.Error(), "0000", err)
		return
	}
	// Send messages from AutoFAQ to webhook client
	prom.ReceivedMessages.With(prometheus.Labels{}).Add(float64(len(messages.Messages)))
	for _, message := range messages.Messages {
		go h.sendCallback(message)
	}
}

func (h *afHandler) sendCallback(message models.Message) {
	// Marshal message
	msgJSON, _ := json.Marshal(message)
	// Create request to webhook client
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8001", bytes.NewBuffer(msgJSON))
	if err != nil {
		h.errc <- apperror.NewError("Create request to callback client error", err.Error(), "0000", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	// Make client and request
	// 3 attempt to send message
	client := http.Client{}
	for i := 1; i <= 3; i++ {
		resp, err := client.Do(req)
		if err != nil {
			if i <= 2 {
				h.errc <- apperror.NewError(fmt.Sprintf("Send request to callback client error. Attempt %d in 10 seconds", i+1), err.Error(), "0000", err)
				time.Sleep(10 * time.Second)
				continue
			}
			h.errc <- apperror.NewError(
				fmt.Sprintf("Message not sent: %v", message),
				err.Error(),
				"0000",
				err)
			break
		}
		if resp.StatusCode == 200 {
			h.infoc <- fmt.Sprintf("Message sent. ID: %s. Conversation: %s", message.Id, message.ConversationId)
			prom.SendMessages.With(prometheus.Labels{}).Inc()
			break
		} else {
			if i <= 2 {
				h.errc <- apperror.NewError(
					fmt.Sprintf("Send message error: response from server not 200. Attempt %d in 10 seconds", i+1),
					fmt.Sprintf("Bad status code from URL %s", req.URL),
					"0000",
					nil)
				time.Sleep(10 * time.Second)
				continue
			}
			h.errc <- apperror.NewError(
				fmt.Sprintf("Message not sent: %v", message),
				fmt.Sprintf("Bad status code from URL %s", req.URL),
				"0000",
				nil)
		}
	}
}
