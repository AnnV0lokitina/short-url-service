package main

import (
	"encoding/json"
	"flag"
	"github.com/caarlos0/env/v6"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"  envDefault:"localhost:8080" json:"server_address,omitempty"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080" json:"base_url,omitempty"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"" json:"file_storage_path,omitempty"`
	DataBaseDSN     string `env:"DATABASE_DSN" envDefault:"" json:"database_dsn,omitempty"`
	EnableHTTPS     bool   `env:"ENABLE_HTTPS" json:"enable_https,omitempty"`
	Config          string `env:"CONFIG" envDefault:"" json:"-"`
	TrustedSubnet   string `env:"TRUSTED_SUBNET" envDefault:"" json:"trusted_subnet,omitempty"`
}

type paramsConfig struct {
	ServerAddress   *string
	BaseURL         *string
	FileStoragePath *string
	DataBaseDSN     *string
	EnableHTTPS     *bool
	Config          *string
	TrustedSubnet   *string
}

func InitParams() *paramsConfig {
	cfg := paramsConfig{}

	flag.Func("a", "Server address", func(flagValue string) error {
		cfg.ServerAddress = &flagValue
		return nil
	})
	flag.Func("b", "Base URL", func(flagValue string) error {
		cfg.BaseURL = &flagValue
		return nil
	})
	flag.Func("f", "File storage path", func(flagValue string) error {
		cfg.FileStoragePath = &flagValue
		return nil
	})
	flag.Func("d", "DB connect string", func(flagValue string) error {
		cfg.DataBaseDSN = &flagValue
		return nil
	})
	flag.Func("s", "Enable HTTPS", func(flagValue string) error {
		v, err := strconv.ParseBool(flagValue)
		if err != nil {
			log.Println("error while parse enable HTTPS flag")
		}
		cfg.EnableHTTPS = &v
		return nil
	})
	flag.Func("c", "config path string", func(flagValue string) error {
		cfg.Config = &flagValue
		return nil
	})
	flag.Func("t", "trusted subnet", func(flagValue string) error {
		cfg.TrustedSubnet = &flagValue
		return nil
	})
	flag.Parse()
	return &cfg
}

func initParamsWithConfig(params *paramsConfig) *config {
	var cfg config

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if params.BaseURL != nil {
		cfg.BaseURL = *params.BaseURL
	}
	if params.ServerAddress != nil {
		cfg.ServerAddress = *params.ServerAddress
	}
	if params.FileStoragePath != nil {
		cfg.FileStoragePath = *params.FileStoragePath
	}
	if params.DataBaseDSN != nil {
		cfg.DataBaseDSN = *params.DataBaseDSN
	}
	if params.EnableHTTPS != nil {
		cfg.EnableHTTPS = *params.EnableHTTPS
	}
	if params.Config != nil {
		cfg.Config = *params.Config
	}
	if params.TrustedSubnet != nil {
		cfg.TrustedSubnet = *params.TrustedSubnet
	}
	return &cfg
}

func InitConfig(params *paramsConfig) *config {
	cfgNoFile := initParamsWithConfig(params)
	if cfgNoFile.Config == "" {
		return cfgNoFile
		//cfgNoFile.Config = "defaults/defaults_run_test.json"
	}
	err := setEnvFromJSON(cfgNoFile.Config)
	if err != nil {
		log.Print(err)
		return cfgNoFile
	}
	return initParamsWithConfig(params)
}

func setEnvFromJSON(path string) error {
	fContent, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	var config = config{}
	if err = json.Unmarshal(fContent, &config); err != nil {
		return err
	}
	if _, hasEnv := os.LookupEnv("SERVER_ADDRESS"); !hasEnv {
		err = os.Setenv("SERVER_ADDRESS", config.ServerAddress)
	}
	if _, hasEnv := os.LookupEnv("BASE_URL"); !hasEnv {
		err = os.Setenv("BASE_URL", config.BaseURL)
	}
	if _, hasEnv := os.LookupEnv("FILE_STORAGE_PATH"); !hasEnv {
		err = os.Setenv("FILE_STORAGE_PATH", config.FileStoragePath)
	}
	if _, hasEnv := os.LookupEnv("DATABASE_DSN"); !hasEnv {
		err = os.Setenv("DATABASE_DSN", config.DataBaseDSN)
	}
	if _, hasEnv := os.LookupEnv("ENABLE_HTTPS"); !hasEnv {
		err = os.Setenv("ENABLE_HTTPS", strconv.FormatBool(config.EnableHTTPS))
	}
	if _, hasEnv := os.LookupEnv("TRUSTED_SUBNET"); !hasEnv {
		err = os.Setenv("TRUSTED_SUBNET", config.TrustedSubnet)
	}
	return err
}
