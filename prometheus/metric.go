package prometheus

import (
	"context"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"os"
	"sync"
	"time"
)

type hub struct {
	reg prometheus.Registerer
	collectors map[string]prometheus.Collector
	exit chan bool

	options Options

	sync.RWMutex
}

type Options struct {
	Address string
	Job string
	Instance string
	SyncTime time.Duration
	Grouping map[string]string

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type Option func(*Options)

// defaultHub add process_collector and go_collector default
var defaultMetricsHub = NewMetricsHub()

func init()  {
	defaultMetricsHub.RegisterCollector("process_collector",prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	defaultMetricsHub.RegisterCollector("go_collector",prometheus.NewGoCollector())
}

func NewMetricsHub() *hub{
	instance, err := os.Hostname()
	if err != nil {
		instance = "unknown"
	}
	o := Options{
		Address:  "http://127.0.0.1:9091",
		Job:      "default-job",
		Instance: instance,
		SyncTime: 5*time.Second,
		Grouping: make(map[string]string),
	}
	return &hub{
		reg:        prometheus.NewRegistry(),
		collectors: make(map[string]prometheus.Collector),
		exit: make(chan bool),
		options: o,
	}
}

func RegisterCollector(key string, c prometheus.Collector) error {
	return defaultMetricsHub.RegisterCollector(key, c)
}

func Address(address string) Option {
	return func(options *Options) {
		options.Address = address
	}
}

func Job(job string) Option {
	return func(options *Options) {
		options.Job = job
	}
}

func SyncTime(duration time.Duration) Option {
	return func(options *Options) {
		options.SyncTime = duration
	}
}

func Instance(instance string) Option {
	return func(options *Options) {
		options.Instance = instance
	}
}

func Grouping(grouping map[string]string) Option {
	return func(options *Options) {
		for name, value := range grouping {
			options.Grouping[name] = value
		}
	}
}

func Context(ctx context.Context) Option {
	return func(options *Options) {
		options.Context = ctx
	}
}

func configure(h *hub, opts ...Option)  {
	for _, o := range opts {
		o(&h.options)
	}
}

func (h *hub) RegisterCollector(key string, c prometheus.Collector) error {
	h.Lock()
	defer h.Unlock()
	if _, ok := h.collectors[key]; ok {
		return errors.New("collector already register")
	}

	if err := h.reg.Register(c); err != nil {
		return err
	}

	h.collectors[key] = c
	return nil
}

func (h *hub) StartMetricsPush(opts ...Option) {
	configure(h, opts...)
	go h.loop()
}

func (h *hub) CloseMetrics()  {
	close(h.exit)
}

func Start(opts ...Option)  {
	defaultMetricsHub.StartMetricsPush(opts...)
}

func Close()  {
	defaultMetricsHub.CloseMetrics()
}

func (h *hub) loop()  {
	ticker := time.NewTicker(h.options.SyncTime)
	defer ticker.Stop()

	for  {
		select {
		case <-ticker.C:
			h.doPush()
		case <-h.exit:
			return
		}
	}
}

func (h *hub) doPush() {
	if len(h.collectors) == 0 {
		return
	}

	pusher := push.New(h.options.Address, h.options.Job).Grouping("instance", h.options.Instance)
	for name, value := range h.options.Grouping {
		pusher = pusher.Grouping(name, value)
	}

	h.RLock()
	defer h.RUnlock()

	for _, c := range h.collectors {
		pusher = pusher.Collector(c)
	}

	if err := pusher.Push(); err != nil {
		fmt.Printf("metric push err:%s \n", err.Error())
	}
}