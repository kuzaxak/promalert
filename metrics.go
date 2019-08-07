package main

// Metrics fetches data from Prometheus.
import (
	"context"
	"fmt"
	prometheus "github.com/prometheus/client_golang/api"
	prometheusApi "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"time"
)

func Metrics(server, query string, queryTime time.Time, duration, step time.Duration) (model.Matrix, error) {
	client, err := prometheus.NewClient(prometheus.Config{Address: server})
	if err != nil {
		return nil, fmt.Errorf("failed to create Prometheus client: %v", err)
	}

	api := prometheusApi.NewAPI(client)
	value, _, err := api.QueryRange(context.Background(), query, prometheusApi.Range{
		Start: queryTime.Add(-duration),
		End:   queryTime,
		Step:  duration / step,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query Prometheus: %v", err)
	}

	metrics, ok := value.(model.Matrix)
	if !ok {
		return nil, fmt.Errorf("unsupported result format: %s", value.Type().String())
	}

	return metrics, nil
}
