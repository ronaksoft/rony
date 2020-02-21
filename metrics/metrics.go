package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// Counters
	CntAuthorizedRequests   = "requests_authorized"
	CntUnauthorizedRequests = "requests_unauthorized"
	CntWebsocketOpen        = "websocket_open"
	CntWebsocketClose       = "websocket_close"
	CntWebsocketWrite       = "websocket_writes"
	CntFunctions            = "functions"
	CntClients              = "clients"
	CntSlowFunctions        = "slow_functions"
	CntNotificationRpc      = "notification_rpc"
	CntFileUploads          = "file_uploads"
	CntFileDownloads        = "file_downloads"
	CntPushedMessages       = "pushed_messages"
	CntPushedUpdates        = "pushed_updates"
	CntBridgeMessage        = "bridge_message"
	CntBridgeContainer      = "bridge_container"
	CntBridgeUndelivered    = "bridge_undelivered"
	CntBridgeDelivered      = "bridge_delivered"
	CntBridgeNotify         = "bridge_notify"
	CntApplePush            = "notification_apple_push"
	CntFirebasePush         = "notification_firebase_push"

	// Gauges
	GaugeActiveConnections        = "active_connections"
	GaugeGatewayWebsocketInQueue  = "gateway_websocket_in_queue"
	GaugeGatewayWebsocketOutQueue = "gateway_websocket_out_queue"
	GaugeActiveUsers24h           = "active_users_24h"
	GaugeActiveUsersLast24h       = "active_users_last_24h"
	GaugeActiveUsers1w            = "active_users_1w"
	GaugeActiveUsersLast1w        = "active_users_last_1w"

	// Histograms
	HistRequestTimeMS        = "request_time"
	HistUploadSizeKB         = "file_uploads_size"
	HistDownloadSizeKB       = "file_download_size"
	HistDbQueryLatenciesMS   = "db_latency"
	HistWebsocketWriteTimeMS = "websocket_write_time"
	HistWebsocketReadTimeMS  = "websocket_read_time"
	HistBucketHandleTime     = "notification_bucket_handle_time"
)

var (
	_Metrics         *Prometheus
	SizeBucketKB     = []float64{1, 5, 10, 20, 40, 80, 160, 320, 640, 1280, 2560}
	TimeBucketMS     = []float64{0.1, 0.5, 1, 10, 20, 50, 100, 500, 1000, 2000, 3000, 5000}
	TimeBucketMicroS = []float64{1, 2, 5, 10, 20, 50, 100, 200, 300, 400, 500, 1000, 2000, 5000, 10000}
)

func Run(bundleID, instanceID string, port int) {
	if _Metrics != nil {
		return
	}
	_Metrics = NewPrometheus(bundleID, instanceID)
	initMetrics()
	go _Metrics.Run(port)
}

func initMetrics() {
	_Metrics.RegisterCounter(CntFileDownloads, "Number of download requests", nil)
	_Metrics.RegisterCounter(CntFileUploads, "Number of uploaded files", nil)
	_Metrics.RegisterCounter(CntPushedUpdates, "Number of messages pushed to the UPDATES redis cache queue", nil)
	_Metrics.RegisterCounter(CntPushedMessages, "Number of messages pushed to the MESSAGES redis cache queue", nil)
	_Metrics.RegisterCounterVec(CntSlowFunctions, "Number of slow functions", nil, []string{"fn"})
	_Metrics.RegisterCounterVec(CntFunctions, "Number of api calls", nil, []string{"fn"})
	_Metrics.RegisterCounterVec(CntClients, "Number of clients by version", nil, []string{"cn"})
	_Metrics.RegisterCounterVec(CntNotificationRpc, "Number of RPCs by name", nil, []string{"name"})
	_Metrics.RegisterCounter(CntAuthorizedRequests, "Number of Authorized requests received from websocket clients", nil)
	_Metrics.RegisterCounter(CntUnauthorizedRequests, "Number of Unauthorized requests received from websocket clients", nil)
	_Metrics.RegisterCounter(CntWebsocketWrite, "Number of messages written to the websocket", nil)
	_Metrics.RegisterCounter(CntWebsocketOpen, "Counts opened websocket connections", nil)
	_Metrics.RegisterCounter(CntWebsocketClose, "Counts closed websocket connections", nil)
	_Metrics.RegisterCounter(CntBridgeDelivered, "Number of delivered bridge messages", nil)
	_Metrics.RegisterCounter(CntBridgeNotify, "Number of notify messages received", nil)
	_Metrics.RegisterCounter(CntBridgeUndelivered, "Number of undelivered bridge messages", nil)
	_Metrics.RegisterCounter(CntBridgeMessage, "Number of received bridge messages", nil)
	_Metrics.RegisterCounter(CntBridgeContainer, "Number of received bridge containers", nil)
	_Metrics.RegisterCounter(CntApplePush, "Number of Apple push messages", nil)
	_Metrics.RegisterCounter(CntFirebasePush, "Number of Firebase push messages", nil)

	_Metrics.RegisterGauge(GaugeActiveConnections, "Number of currently connected websocket clients", nil)
	_Metrics.RegisterGauge(GaugeGatewayWebsocketInQueue, "Number of waiting incoming messages to read", nil)
	_Metrics.RegisterGauge(GaugeGatewayWebsocketOutQueue, "Number of waiting outgoing messages to be sent", nil)
	_Metrics.RegisterGauge(GaugeActiveUsers24h, "Number of unique users in the last 24h", nil)
	_Metrics.RegisterGauge(GaugeActiveUsers1w, "Number of unique users in the current week", nil)
	_Metrics.RegisterGauge(GaugeActiveUsersLast24h, "Number of unique users in two days ago", nil)
	_Metrics.RegisterGauge(GaugeActiveUsersLast1w, "Number of unique users in the last week", nil)

	_Metrics.RegisterHistogram(HistDownloadSizeKB, "Histogram of download size average", SizeBucketKB, nil)
	_Metrics.RegisterHistogram(HistUploadSizeKB, "Histogram of upload size average", SizeBucketKB, nil)
	_Metrics.RegisterHistogram(HistRequestTimeMS, "Histogram of request time", TimeBucketMS, nil)
	_Metrics.RegisterHistogram(HistDbQueryLatenciesMS, "Histogram of the db request latency", TimeBucketMS, nil)
	_Metrics.RegisterHistogram(HistWebsocketReadTimeMS, "Histogram of the readPump time per request", TimeBucketMS, nil)
	_Metrics.RegisterHistogram(HistWebsocketWriteTimeMS, "Histogram of the writePump time per request", TimeBucketMS, nil)
	_Metrics.RegisterHistogram(HistBucketHandleTime, "Processing time of each notification bucket in milliseconds", TimeBucketMS, nil)

}

func Counter(name string) prometheus.Counter {
	return _Metrics.Counter(name)
}

func Gauge(name string) prometheus.Gauge {
	return _Metrics.Gauge(name)
}

func Histogram(name string) prometheus.Histogram {
	return _Metrics.Histogram(name)
}

func CounterVec(name string) *prometheus.CounterVec {
	return _Metrics.CounterVec(name)
}
