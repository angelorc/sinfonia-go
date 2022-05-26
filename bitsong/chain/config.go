package chain

type Config struct {
	ChainID       string `json:"chain-id" yaml:"chain-id"`
	RPCAddr       string `json:"rpc-addr" yaml:"rpc-addr"`
	GRPCAddr      string `json:"grpc-addr" yaml:"grpc-addr"`
	AccountPrefix string `json:"account-prefix" yaml:"account-prefix"`
	Timeout       string `json:"timeout" yaml:"timeout"`
}

func GetBitsongConfig() *Config {
	return &Config{
		ChainID:       "bitsong-sinfonia-test-1",
		RPCAddr:       "https://rpc.testnet.bitsong.network:443",
		GRPCAddr:      "http://142.132.252.143:9090",
		AccountPrefix: "bitsong",
		Timeout:       "10s",
	}
}
