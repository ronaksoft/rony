package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"sync"
)

/*
   Creation Time: 2019 - Dec - 06
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// Prometheus
type Prometheus struct {
	l              sync.RWMutex
	tags           map[string]string
	registry       *prometheus.Registry
	labels         map[string]string
	counters       map[string]prometheus.Counter
	counterVecs    map[string]*prometheus.CounterVec
	cachedCounters map[string]prometheus.Counter
	gauges         map[string]prometheus.Gauge
	gaugeVecs      map[string]*prometheus.GaugeVec
	histograms     map[string]prometheus.Histogram
}

func NewPrometheus(bundleID, instanceID string) *Prometheus {
	m := new(Prometheus)
	m.tags = make(map[string]string)
	m.counters = make(map[string]prometheus.Counter)
	m.counterVecs = make(map[string]*prometheus.CounterVec)
	m.gauges = make(map[string]prometheus.Gauge)
	m.gaugeVecs = make(map[string]*prometheus.GaugeVec)
	m.histograms = make(map[string]prometheus.Histogram)
	m.registry = prometheus.NewRegistry()

	m.labels = map[string]string{
		"bundleID":   bundleID,
		"instanceID": instanceID,
	}

	return m
}

func (p *Prometheus) RegisterGoCollector() {
	p.registry.MustRegister(prometheus.NewGoCollector())
}

func (p *Prometheus) SetConstLabels(m map[string]string) {
	appendMap(&p.labels, &m)
}

func (p *Prometheus) Run(port int) error {
	http.Handle("/metrics", promhttp.HandlerFor(p.registry, promhttp.HandlerOpts{}))
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func (p *Prometheus) Counter(name string) prometheus.Counter {
	return p.counters[name]
}

func (p *Prometheus) CounterVec(name string) *prometheus.CounterVec {
	return p.counterVecs[name]
}

func (p *Prometheus) CounterVecWithLabel(name string, labels ...string) prometheus.Counter {
	return p.counterVecs[name].WithLabelValues(labels...)
}

func (p *Prometheus) RegisterCounter(name, help string, labels map[string]string) {
	appendMap(&labels, &p.labels)
	p.counters[name] = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   "cnt",
		Name:        name,
		Help:        help,
		ConstLabels: labels,
	})
	p.registry.MustRegister(p.counters[name])
}

func (p *Prometheus) RegisterCounterVec(name, help string, constLabels map[string]string, varLabels []string) {
	appendMap(&constLabels, &p.labels)
	p.counterVecs[name] = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   "cnt",
		Name:        name,
		Help:        help,
		ConstLabels: constLabels,
	}, varLabels)
	p.registry.MustRegister(p.counterVecs[name])
}

func (p *Prometheus) Gauge(name string) prometheus.Gauge {
	return p.gauges[name]
}

func (p *Prometheus) GaugeVec(name string) *prometheus.GaugeVec {
	return p.gaugeVecs[name]
}

func (p *Prometheus) RegisterGauge(name, help string, labels map[string]string) {
	appendMap(&labels, &p.labels)
	p.gauges[name] = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   "gauge",
		Name:        name,
		Help:        help,
		ConstLabels: labels,
	})
	p.registry.MustRegister(p.gauges[name])

}

func (p *Prometheus) RegisterGaugeVec(name, help string, constLabels map[string]string, varLabels []string) {
	appendMap(&constLabels, &p.labels)
	p.gaugeVecs[name] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   "gauge",
		Name:        name,
		Help:        help,
		ConstLabels: constLabels,
	}, varLabels)
	p.registry.MustRegister(p.gaugeVecs[name])
}

func (p *Prometheus) Histogram(name string) prometheus.Histogram {
	return p.histograms[name]
}

func (p *Prometheus) RegisterHistogram(name, help string, buckets []float64, labels map[string]string) {
	appendMap(&labels, &p.labels)
	p.histograms[name] = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace:   "hist",
		Name:        name,
		Help:        help,
		ConstLabels: labels,
		Buckets:     buckets,
	})
	p.registry.MustRegister(p.histograms[name])

}

func appendMap(m1, m2 *map[string]string) {
	if *m1 == nil {
		*m1 = make(map[string]string)
	}
	for k, v := range *m2 {
		(*m1)[k] = v
	}
}
