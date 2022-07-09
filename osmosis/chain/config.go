package chain

type Config struct {
	ChainID       string `json:"chain-id" yaml:"chain-id"`
	RPCAddr       string `json:"rpc-addr" yaml:"rpc-addr"`
	GRPCAddr      string `json:"grpc-addr" yaml:"grpc-addr"`
	GRPCInsecure  bool   `json:"grpc-insecure" yaml:"grpc-insecure"`
	AccountPrefix string `json:"account-prefix" yaml:"account-prefix"`
	Timeout       string `json:"timeout" yaml:"timeout"`
}

func GetOsmosisConfig() *Config {
	return &Config{
		ChainID:       "osmosis-sinfonia-test-1",
		RPCAddr:       "https://rpc.osmo-test.bitsong.network:443",
		GRPCAddr:      "http://157.90.168.95:9090",
		GRPCInsecure:  true,
		AccountPrefix: "osmo",
		Timeout:       "10s",
	}
}
