package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	CntGatewayIncomingHttpMessage      = "gateway_incoming_http_message"
	CntGatewayIncomingWebsocketMessage = "gateway_incoming_websocket_message"
	CntGatewayOutgoingHttpMessage      = "gateway_outgoing_http_message"
	CntGatewayOutgoingWebsocketMessage = "gateway_outgoing_websocket_message"
	CntTunnelIncomingMessage           = "tunnel_incoming_message"
	CntTunnelOutgoingMessage           = "tunnel_outgoing_message"
	CntStoreConflicts                  = "store_conflicts"
	GaugeActiveWebsocketConnections    = "gateway_active_websocket_conns"
	HistGatewayRequestTime             = "gateway_request_time"
	HistTunnelRequestTime              = "tunnel_request_time"
	HistTunnelRoundtripTime            = "tunnel_roundtrip_time"
)

var (
	_Prom            *Prometheus
	SizeBucketKB     = []float64{1, 5, 10, 20, 40, 80, 160, 320, 640, 1280, 2560}
	TimeBucketMS     = []float64{0.1, 0.5, 1, 10, 20, 50, 100, 500, 1000, 2000, 3000, 5000}
	TimeBucketMicroS = []float64{1, 2, 5, 10, 20, 50, 100, 200, 300, 400, 500, 1000, 2000, 5000, 10000}
)

func Init(constLabels map[string]string) {
	if _Prom != nil {
		return
	}
	_Prom = NewPrometheus("rony", constLabels)
	_Prom.RegisterCounter(CntGatewayIncomingHttpMessage, "number of incoming http messages", nil)
	_Prom.RegisterCounter(CntGatewayOutgoingHttpMessage, "number of outgoing http messages", nil)
	_Prom.RegisterCounter(CntGatewayIncomingWebsocketMessage, "number of incoming websocket messages", nil)
	_Prom.RegisterCounter(CntGatewayOutgoingWebsocketMessage, "number of outgoing websocket messages", nil)
	_Prom.RegisterCounter(CntTunnelIncomingMessage, "number of incoming messages", nil)
	_Prom.RegisterCounter(CntTunnelOutgoingMessage, "number of outgoing messages", nil)
	_Prom.RegisterCounter(CntStoreConflicts, "number of txn conflicts", nil)

	_Prom.RegisterGauge(GaugeActiveWebsocketConnections, "number of gateway active websocket connections", nil)

	_Prom.RegisterHistogram(HistGatewayRequestTime, "the amount of process time for gateway requests", TimeBucketMS, nil)
	_Prom.RegisterHistogram(HistTunnelRequestTime, "the amount of process time for tunnel requests", TimeBucketMS, nil)
	_Prom.RegisterHistogram(HistTunnelRoundtripTime, "the roundtrip of a execute remote command", TimeBucketMS, nil)
}

func Register(registerer prometheus.Registerer) {
	for _, m := range _Prom.counters {
		registerer.MustRegister(m)
	}
	for _, m := range _Prom.counterVectors {
		registerer.MustRegister(m)
	}
	for _, m := range _Prom.gauges {
		registerer.MustRegister(m)
	}
	for _, m := range _Prom.gaugeVectors {
		registerer.MustRegister(m)
	}
	for _, m := range _Prom.histograms {
		registerer.MustRegister(m)
	}
}

func IncCounter(name string) {
	if _Prom == nil {
		return
	}
	_Prom.Counter(name).Inc()
}

func InCounterVec(name string, labelValues ...string) {
	if _Prom == nil {
		return
	}
	_Prom.CounterVec(name).WithLabelValues(labelValues...).Inc()
}

func AddCounter(name string, v float64) {
	if _Prom == nil {
		return
	}
	_Prom.Counter(name).Add(v)
}

func AddCounterVec(name string, v float64, labelValues ...string) {
	if _Prom == nil {
		return
	}
	_Prom.CounterVec(name).WithLabelValues(labelValues...).Add(v)
}

func SetGauge(name string, v float64) {
	if _Prom == nil {
		return
	}
	_Prom.Gauge(name).Set(v)
}

func ObserveHistogram(name string, v float64) {
	if _Prom == nil {
		return
	}
	_Prom.Histogram(name).Observe(v)
}
