package prometheus

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"math"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

var (
	normMean = 0.00001
	normDomain = 0.0002
	uniformDomain = 0.0002
	oscillationPeriod = 1*time.Minute

	oscillationFactor = func(start time.Time) float64 {
		return 2 + math.Sin(math.Sin(2*math.Pi*float64(time.Since(start))/float64(oscillationPeriod)))
	}
)

var (
	// Create a summary to track fictional interservice RPC latencies for three
	// distinct services with different latency distributions. These services are
	// differentiated via a "service" label.
	rpcDurations = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "rpc_durations_seconds",
			Help:       "RPC latency distributions.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"service"},
	)
	// The same as above, but now as a histogram, and only for the normal
	// distribution. The buckets are targeted to the parameters of the
	// normal distribution, with 20 buckets centered on the mean, each
	// half-sigma wide.
	rpcDurationsHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "rpc_durations_histogram_seconds",
		Help:    "RPC latency distributions.",
		Buckets: prometheus.LinearBuckets(normMean-5*normDomain, .5*normDomain, 20),
	})

	httpReqTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_request_total",
		Help: "Request total",
	}, []string{"method", "path", "status"})

	goGoroutines = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "go_goroutine",
		Help: "current goroutine num",
	}, []string{"service"})

)

func startReqTotalCollect()  {
	for {
		code := rand.Intn(10)
		httpReqTotal.WithLabelValues("Get", string(code+65), "200").Inc()
		time.Sleep(1*time.Second)
	}
}

func startGoGoroutines() {
	for {
		goGoroutines.WithLabelValues("test_server").Set(float64(runtime.NumGoroutine()))
		time.Sleep(5*time.Second)
	}
}

func startRpcDurationsCollect()  {
	for {
		v := rand.Float64() * normDomain
		rpcDurations.WithLabelValues("normal").Observe(v)
		time.Sleep(time.Duration(100*oscillationFactor(time.Now())) * time.Millisecond)
	}
}

func startRpcDurationsHistogramCollect()  {
	for {
		v := (rand.NormFloat64() * normDomain) + normMean
		rpcDurations.WithLabelValues("normal").Observe(v)
		// Demonstrate exemplar support with a dummy ID. This
		// would be something like a trace ID in a real
		// application.  Note the necessary type assertion. We
		// already know that rpcDurationsHistogram implements
		// the ExemplarObserver interface and thus don't need to
		// check the outcome of the type assertion.
		rpcDurationsHistogram.(prometheus.ExemplarObserver).ObserveWithExemplar(
			v, prometheus.Labels{"dummyID": fmt.Sprint(rand.Intn(100000))},
		)
		time.Sleep(time.Duration(75*oscillationFactor(time.Now())) * time.Millisecond)
	}
}

func TestStart(t *testing.T) {
	if err := RegisterCollector("http_request_total", httpReqTotal); err != nil {
		t.Error(err)
	}
	if err := RegisterCollector("go_goroutine", goGoroutines); err != nil {
		t.Error(err)
	}
	if err := RegisterCollector("rpc_durations_seconds", rpcDurations); err != nil {
		t.Error(err)
	}
	if err := RegisterCollector("rpc_durations_histogram_seconds", rpcDurationsHistogram); err != nil {
		t.Error(err)
	}

	// start collect
	go startReqTotalCollect()
	go startGoGoroutines()
	go startRpcDurationsCollect()
	go startRpcDurationsHistogramCollect()

	// start push
	Start(Address("http://124.71.103.18:9091"), Job("test-job-defaultHub"))

	time.Sleep(2*time.Minute)

}

func TestNewMetricsHub(t *testing.T) {
	h := NewMetricsHub()
	if err := h.RegisterCollector("http_request_total", httpReqTotal); err != nil {
		t.Error(err)
	}
	if err := h.RegisterCollector("go_goroutine", goGoroutines); err != nil {
		t.Error(err)
	}
	if err := h.RegisterCollector("rpc_durations_seconds", rpcDurations); err != nil {
		t.Error(err)
	}
	if err := h.RegisterCollector("rpc_durations_histogram_seconds", rpcDurationsHistogram); err != nil {
		t.Error(err)
	}

	// start collect
	go startReqTotalCollect()
	go startGoGoroutines()
	go startRpcDurationsCollect()
	go startRpcDurationsHistogramCollect()

	// start push
	h.StartMetricsPush(Address("http://124.71.103.18:9091"), Job("test-job-newtHub"))

	time.Sleep(20*time.Minute)

}
