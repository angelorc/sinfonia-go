package chain

type Config struct {
	ChainID       string `json:"chain-id" yaml:"chain-id"`
	RPCAddr       string `json:"rpc-addr" yaml:"rpc-addr"`
	GRPCAddr      string `json:"grpc-addr" yaml:"grpc-addr"`
	AccountPrefix string `json:"account-prefix" yaml:"account-prefix"`
	Timeout       string `json:"timeout" yaml:"timeout"`
}

func GetOsmosisConfig() *Config {
	return &Config{
		ChainID:       "osmosis-sinfonia-test-1",
		RPCAddr:       "https://rpc.osmosis.devnet.bitsong.network:443",
		GRPCAddr:      "http://142.132.252.143:10090",
		AccountPrefix: "osmo",
		Timeout:       "10s",
	}
}
