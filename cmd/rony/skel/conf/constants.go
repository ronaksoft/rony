package c

import "time"

const (
	ServerID                 = "server.id"
	GatewayListen            = "gateway.listen"
	GatewayAdvertiseHostPort = "gateway.advertise.hostport"
	TunnelListen             = "tunnel.listen"
	TunnelAdvertiseHostPort  = "tunnel.advertise.hostport"
	IdleTime                 = "idle-time"
	ReplicaSet               = "replica-set"
	GossipPort               = "gossip.port"
	DataPath                 = "data.path"
	Bootstrap                = "bootstrap"
	Host                     = "host"
	Port                     = "post"
)

func GetReplicaSet() uint64 {
	return globalC.GetUint64(ReplicaSet)
}

func GetIdleTime() time.Duration {
	return globalC.GetDuration(IdleTime)
}

func GetGatewayListen() string {
	return globalC.GetString(GatewayListen)
}

func GetGatewayAdvertiseHostPort() []string {
	return globalC.GetStringSlice(GatewayAdvertiseHostPort)
}

func GetTunnelListen() string {
	return globalC.GetString(TunnelListen)
}

func GetTunnelAdvertiseHostPort() []string {
	return globalC.GetStringSlice(TunnelAdvertiseHostPort)
}

func GetServerID() string {
	return globalC.GetString(ServerID)
}

func GetGossipPort() int {
	return globalC.GetInt(GossipPort)
}

func GetDataPath() string {
	return globalC.GetString(DataPath)
}

func GetBootstrap() bool {
	return globalC.GetBool(Bootstrap)
}

func GetHost() string {
	return globalC.GetString(Host)
}

func GetPort() int {
	return globalC.GetInt(Port)
}
