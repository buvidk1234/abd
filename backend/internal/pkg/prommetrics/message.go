package prommetrics

import "github.com/prometheus/client_golang/prometheus"

var (
	MsgProcessSuccessCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "msg_process_success_total",
		Help: "The number of msg successful processed",
	})
	MsgProcessFailedCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "msg_process_failed_total",
		Help: "The number of msg failed processed",
	})
)

func RegistryMsg() {
	registry.MustRegister(
		MsgProcessSuccessCounter,
		MsgProcessFailedCounter,
	)
}
