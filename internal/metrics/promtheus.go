package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

/*
   Creation Time: 2019 - Dec - 06
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

type Prometheus struct {
	ns             string
	labels         map[string]string
	counters       map[string]prometheus.Counter
	counterVectors map[string]*prometheus.CounterVec
	gauges         map[string]prometheus.Gauge
	gaugeVectors   map[string]*prometheus.GaugeVec
	histograms     map[string]prometheus.Histogram
}

func NewPrometheus(ns string, constLabels map[string]string) *Prometheus {
	m := &Prometheus{
		ns:             ns,
		counters:       make(map[string]prometheus.Counter),
		counterVectors: make(map[string]*prometheus.CounterVec),
		gauges:         make(map[string]prometheus.Gauge),
		gaugeVectors:   make(map[string]*prometheus.GaugeVec),
		histograms:     make(map[string]prometheus.Histogram),
		labels:         constLabels,
	}

	return m
}

func (p *Prometheus) SetConstLabels(m map[string]string) {
	appendMap(&p.labels, &m)
}

func (p *Prometheus) Counter(name string) prometheus.Counter {
	return p.counters[name]
}

func (p *Prometheus) CounterVec(name string) *prometheus.CounterVec {
	return p.counterVectors[name]
}

func (p *Prometheus) CounterVecWithLabelValues(name string, labelValues ...string) prometheus.Counter {
	return p.counterVectors[name].WithLabelValues(labelValues...)
}

func (p *Prometheus) CounterVecWithLabels(name string, labels prometheus.Labels) prometheus.Counter {
	return p.counterVectors[name].With(labels)
}

func (p *Prometheus) RegisterCounter(name, help string, labels map[string]string) {
	appendMap(&labels, &p.labels)
	p.counters[name] = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   p.ns,
		Name:        name,
		Help:        help,
		ConstLabels: labels,
	})
}

func (p *Prometheus) RegisterCounterVec(name, help string, constLabels map[string]string, varLabels []string) {
	appendMap(&constLabels, &p.labels)
	p.counterVectors[name] = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   p.ns,
		Name:        name,
		Help:        help,
		ConstLabels: constLabels,
	}, varLabels)
}

func (p *Prometheus) Gauge(name string) prometheus.Gauge {
	return p.gauges[name]
}

func (p *Prometheus) GaugeVec(name string) *prometheus.GaugeVec {
	return p.gaugeVectors[name]
}

func (p *Prometheus) RegisterGauge(name, help string, labels map[string]string) {
	appendMap(&labels, &p.labels)
	p.gauges[name] = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   p.ns,
		Name:        name,
		Help:        help,
		ConstLabels: labels,
	})
}

func (p *Prometheus) RegisterGaugeVec(name, help string, constLabels map[string]string, varLabels []string) {
	appendMap(&constLabels, &p.labels)
	p.gaugeVectors[name] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   p.ns,
		Name:        name,
		Help:        help,
		ConstLabels: constLabels,
	}, varLabels)
}

func (p *Prometheus) Histogram(name string) prometheus.Histogram {
	return p.histograms[name]
}

func (p *Prometheus) RegisterHistogram(name, help string, buckets []float64, labels map[string]string) {
	appendMap(&labels, &p.labels)
	p.histograms[name] = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace:   p.ns,
		Name:        name,
		Help:        help,
		ConstLabels: labels,
		Buckets:     buckets,
	})
}

func appendMap(m1, m2 *map[string]string) {
	if *m1 == nil {
		*m1 = make(map[string]string)
	}
	for k, v := range *m2 {
		(*m1)[k] = v
	}
}
