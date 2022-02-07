package prom

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ReceivedMessages = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "receive_messages_total",
			Help: "The total number messages are received from AutoFAQ",
		},
		[]string{})

	SendMessages = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "send_messages_total",
			Help: "The total number messages are send callback client",
		},
		[]string{})

	Errors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "errors_total",
			Help: "The total number of errors",
		},
		[]string{"code"})
)
