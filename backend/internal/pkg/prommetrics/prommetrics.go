package prommetrics

import (
	"errors"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var registry = &prometheusRegistry{prometheus.NewRegistry()}

// prometheusRegistry 封装 prometheus.Registry，重复注册时不 panic
type prometheusRegistry struct {
	*prometheus.Registry
}

func (r *prometheusRegistry) MustRegister(cs ...prometheus.Collector) {
	for _, c := range cs {
		if err := r.Registry.Register(c); err != nil {
			// 已注册则跳过，避免重复注册 panic
			if errors.As(err, &prometheus.AlreadyRegisteredError{}) {
				continue
			}
			panic(err)
		}
	}
}

func init() {
	// 自动注册 Go 运行时指标（goroutine 数、GC 等）和进程指标（CPU、内存等）
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
}

// RegistryAll 注册所有业务指标
func RegistryAll() {
	RegistryMsg()
	RegistryMsgGateway()
}

// Start 启动 Prometheus HTTP 端点，暴露 /metrics
func Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	return http.Serve(listener, mux)
}
