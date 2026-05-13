package observability

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
)

// SetupPrometheusMetrics регистрирует глобальный MeterProvider и возвращает HTTP-handler с /metrics в формате Prometheus.
// Вызовите до api.NewServer, чтобы сгенерированный ogen-сервер использовал этот провайдер.
func SetupPrometheusMetrics() (metrics http.Handler, shutdown func(context.Context) error, err error) {
	reg := prometheus.NewRegistry()
	exporter, err := otelprom.New(otelprom.WithRegisterer(reg))
	if err != nil {
		return nil, nil, fmt.Errorf("otel prometheus exporter: %w", err)
	}
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(exporter))
	otel.SetMeterProvider(provider)

	handler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{
		Registry: reg,
	})
	return handler, provider.Shutdown, nil
}
