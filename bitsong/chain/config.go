package chain

type Config struct {
	ChainID       string `json:"chain-id" yaml:"chain-id"`
	RPCAddr       string `json:"rpc-addr" yaml:"rpc-addr"`
	GRPCAddr      string `json:"grpc-addr" yaml:"grpc-addr"`
	GRPCInsecure  bool   `json:"grpc-insecure" yaml:"grpc-insecure"`
	AccountPrefix string `json:"account-prefix" yaml:"account-prefix"`
	Timeout       string `json:"timeout" yaml:"timeout"`
}

/*
rpc.fantest-1.bitsong.network
api.fantest-1.bitsong.network

rpc.osmo-test.bitsong.network
api.osmo-test.bitsong.network
*/

func GetBitsongConfig() *Config {
	/*return &Config{
		ChainID:       "bitsong-sinfonia-test-1",
		RPCAddr:       "https://rpc.testnet.bitsong.network:443",
		GRPCAddr:      "http://142.132.252.143:9090",
		GRPCInsecure:  true,
		AccountPrefix: "bitsong",
		Timeout:       "10s",
	}*/
	return &Config{
		ChainID:       "bitsong-2b",
		RPCAddr:       "https://rpc.explorebitsong.com:443",
		GRPCAddr:      "http://88.99.184.249:9090",
		GRPCInsecure:  true,
		AccountPrefix: "bitsong",
		Timeout:       "10s",
	}
}
