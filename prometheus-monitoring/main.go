package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strconv"
	"time"
)

var (
	// How often our /latest/volume request durations fall into one of the defined buckets.
	// We can use default buckets or set ones we are interested in.
	volumeRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "volume_request_duration_seconds",
			Help:    "Histogram of the /latest/volumes request duration.",
			Buckets: []float64{0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"code", "method"},
	)
	// Counter vector to which we can attach labels. That creates many key-value
	// label combinations. So in our case we count requests by status code separetly.
	volumeRequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "volume_requests_total",
			Help: "Total number of /latest/volumes requests.",
		},
		[]string{"code", "method"},
	)
)

// init registers Prometheus metrics.
func init() {
	prometheus.MustRegister(volumeRequestDuration)
	prometheus.MustRegister(volumeRequestCounter)
}

func handler(w http.ResponseWriter, r *http.Request) {
	var code int
	defer func(begun time.Time) {
		volumeRequestDuration.WithLabelValues(strconv.Itoa(code), r.Method).Observe(time.Since(begun).Seconds())

		// volume_requests_total{status="200"}
		volumeRequestCounter.WithLabelValues(strconv.Itoa(code), r.Method).Inc()
	}(time.Now())
	code = 200
	w.WriteHeader(code)

	fmt.Fprintf(w, "Hi there, the end point is :  %s !", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/latest/volumes", handler)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
