package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

var (
	// How often our /hello request durations fall into one of the defined buckets.
	// We can use default buckets or set ones we are interested in.
	duration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "hello_request_duration_seconds",
		Help:    "Histogram of the /hello request duration.",
		Buckets: []float64{0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
	})
	// Counter vector to which we can attach labels. That creates many key-value
	// label combinations. So in our case we count requests by status code separetly.
	counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hello_requests_total",
			Help: "Total number of /hello requests.",
		},
		[]string{"status"},
	)
)

// init registers Prometheus metrics.
func init() {
	prometheus.MustRegister(duration)
	prometheus.MustRegister(counter)
}

func handler(w http.ResponseWriter, r *http.Request) {
	var status int

	defer func(begun time.Time) {
		duration.Observe(time.Since(begun).Seconds())

		// hello_requests_total{status="200"} 2385
		counter.With(prometheus.Labels{
			"status": fmt.Sprint(status),
		}).Inc()
	}(time.Now())

	status = 200
	w.WriteHeader(status)

	fmt.Fprintf(w, "Hi there, the end point is :  %s !", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
