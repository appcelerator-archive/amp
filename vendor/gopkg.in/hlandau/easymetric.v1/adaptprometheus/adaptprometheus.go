// Package adaptprometheus adapts measurables to be exported by Prometheus.
package adaptprometheus

import "gopkg.in/hlandau/measurable.v1"
import "github.com/prometheus/client_golang/prometheus"
import "sync"
import "errors"
import "regexp"
import "net/http"

type metric struct {
	Measurable measurable.Measurable
	Metric     prometheus.Metric
}

var errNotSupported = errors.New("measurable type not supported")

func (m *metric) init() error {
	name := m.Measurable.MsName()
	mangledName := mangleName(name)

	opts := prometheus.Opts{
		Name: mangledName,
		Help: name,
	}

	switch m.Measurable.MsType() {
	case measurable.CounterType:
		mi, ok := m.Measurable.(interface {
			MsInt64() int64
		})
		if !ok {
			return errNotSupported
		}

		m.Metric = prometheus.NewCounterFunc(prometheus.CounterOpts(opts), func() float64 {
			return float64(mi.MsInt64())
		})

	case measurable.GaugeType:
		mi, ok := m.Measurable.(interface {
			MsInt64() int64
		})
		if !ok {
			return errNotSupported
		}

		m.Metric = prometheus.NewGaugeFunc(prometheus.GaugeOpts(opts), func() float64 {
			return float64(mi.MsInt64())
		})

	default:
		return errNotSupported
	}

	return nil
}

type collector struct{}

var metricsMutex sync.RWMutex
var metrics = map[string]*metric{}

func (c *collector) Describe(descChan chan<- *prometheus.Desc) {
	metricsMutex.RLock()
	defer metricsMutex.RUnlock()

	for _, m := range metrics {
		descChan <- m.Metric.Desc()
	}
}

func (c *collector) Collect(metricChan chan<- prometheus.Metric) {
	metricsMutex.RLock()
	defer metricsMutex.RUnlock()

	for _, m := range metrics {
		metricChan <- m.Metric
	}
}

var col collector

var re_mangler = regexp.MustCompilePOSIX(`[^a-zA-Z0-9_:]`)

func mangleName(metricName string) string {
	return re_mangler.ReplaceAllString(metricName, "_")
}

var once, handlerOnce sync.Once
var hookKey int

func hook(m measurable.Measurable, event measurable.HookEvent) {
	name := m.MsName()

	switch event {
	case measurable.RegisterEvent, measurable.RegisterCatchupEvent:
		mi := &metric{
			Measurable: m,
		}
		err := mi.init()
		if err != nil {
			return
		}

		metricsMutex.Lock()
		defer metricsMutex.Unlock()
		metrics[name] = mi

	case measurable.UnregisterEvent:
		metricsMutex.Lock()
		defer metricsMutex.Unlock()

		delete(metrics, name)
	}
}

func RegisterNoNexus() {
	once.Do(func() {
		measurable.RegisterHook(&hookKey, hook)
		prometheus.Register(&col)
	})
}

func Register() {
	RegisterNoNexus()
	handlerOnce.Do(func() {
		http.Handle("/metrics", prometheus.Handler())
	})
}
