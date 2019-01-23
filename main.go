package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promauto"
        "github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	goapp_requests_count = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "goapp_requests_count",
			Help: "Number of requests served",
		},
		[]string{"type"},
	)
)

func init() {
	prometheus.MustRegister(goapp_requests_count)
}

func handler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["a"]
	if !ok || len(keys[0]) < 1 {
		goapp_requests_count.With(prometheus.Labels{"type": "missinga"}).Inc()
		log.Println("Url Param 'a' is missing")
		return
	}

	a, err := strconv.Atoi(keys[0])

	if err != nil {
		goapp_requests_count.With(prometheus.Labels{"type": "bada"}).Inc()
		return
	}

	keys, ok = r.URL.Query()["b"]
	if !ok || len(keys[0]) < 1 {
		goapp_requests_count.With(prometheus.Labels{"type": "missingb"}).Inc()
		log.Println("Url Param 'b' is missing")
		return
	}
	b, err := strconv.Atoi(keys[0])
	if err != nil {
		goapp_requests_count.With(prometheus.Labels{"type": "badb"}).Inc()
		log.Println("Url Param 'b' is missing")
		log.Println("Bad number!")
		return
	}

	goapp_requests_count.With(prometheus.Labels{"type": "ok"}).Inc()
	fmt.Fprintf(w, "%d + %d = %d", a, b, sum(a, a))
}

func main() {
	http.HandleFunc("/", handler)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func sum(a, b int) int {
	t := a + b
	return t
}
