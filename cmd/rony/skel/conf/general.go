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
	GossipSeedHostPort       = "gossip.seed.hostport"
	DataPath                 = "data.path"
	Bootstrap                = "bootstrap"
	Host                     = "host"
	Port                     = "port"
	Config                   = "config"
)

func GetConfig() string {
	return globalC.GetString(Config)
}

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

func GetGossipSeedHostPort() string {
	return globalC.GetString(GossipSeedHostPort)
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
