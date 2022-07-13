package config

import (
	"flag"
	"fmt"
	"github.com/angelorc/sinfonia-go/utility"
	"gopkg.in/yaml.v2"
	"os"
)

type GraphQL struct {
	Address        string `yaml:"address" validate:"required"`
	Port           string `yaml:"port" validate:"required"`
	Endpoint       string `yaml:"endpoint" validate:"required"`
	PlaygroundPass string `yaml:"playground_pass"`
}

type Mongo struct {
	Uri    string `yaml:"uri" validate:"required"`
	DbName string `yaml:"dbname" validate:"required"`
	Retry  bool   `yaml:"retry" validate:"required"`
}

type ChainConfig struct {
	ChainID       string `yaml:"chain-id" validate:"required"`
	RPCAddr       string `yaml:"rpc-addr" validate:"required"`
	GRPCAddr      string `yaml:"grpc-addr" validate:"required"`
	GRPCInsecure  bool   `yaml:"grpc-insecure" validate:"required"`
	AccountPrefix string `yaml:"account-prefix" validate:"required"`
	Timeout       string `yaml:"timeout" validate:"required"`
}

type CloudflareConfig struct {
	Account string `yaml:"account" validate:"required"`
	Images  string `yaml:"images" validate:"required"`
}

type Config struct {
	GraphQL    GraphQL          `yaml:"graphql" validate:"required"`
	Mongo      Mongo            `yaml:"mongo" validate:"required"`
	Cloudflare CloudflareConfig `yaml:"cloudflare" validate:"required"`
	Bitsong    ChainConfig      `yaml:"bitsong" validate:"required"`
	Osmosis    ChainConfig      `yaml:"osmosis" validate:"required"`
}

func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	if err := utility.ValidateStruct(config); err != nil {
		return nil, err
	}

	return config, nil
}

func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

func ParseFlags() (string, error) {
	// String that contains the configured configuration path
	var configPath string

	// Set up a CLI flag called "-config" to allow users
	// to supply the configuration file
	flag.StringVar(&configPath, "config", "./config.yml", "path to config file")

	// Actually parse the flags
	flag.Parse()

	// Validate the path first
	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}

	// Return the configuration path
	return configPath, nil
}
