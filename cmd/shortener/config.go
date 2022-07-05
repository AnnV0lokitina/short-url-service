package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/caarlos0/env/v6"
)

type config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"  envDefault:"localhost:8080" json:"server_address"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080" json:"base_url"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"" json:"file_storage_path"`
	DataBaseDSN     string `env:"DATABASE_DSN" envDefault:"" json:"database_dsn"`
	EnableHTTPS     bool   `env:"ENABLE_HTTPS" json:"enable_https"`
	Config          string `env:"CONFIG" envDefault:"" json:"-"`
}

type paramsConfig struct {
	ServerAddress   *string
	BaseURL         *string
	FileStoragePath *string
	DataBaseDSN     *string
	EnableHTTPS     *bool
	Config          *string
}

//type fileConfig struct {
//	ServerAddress   string `json:"server_address"`
//	BaseURL         string `json:"base_url"`
//	FileStoragePath string `json:"file_storage_path"`
//	DataBaseDSN     string `json:"database_dsn"`
//	EnableHTTPS     bool   `json:"enable_https"`
//}

func initParams() *paramsConfig {
	cfg := paramsConfig{}
	args := strings.Join(os.Args[1:], ",")
	if strings.Contains(args, "-a,") {
		cfg.ServerAddress = flag.String("a", "", "Server address")
	}
	if strings.Contains(args, "-b,") {
		cfg.BaseURL = flag.String("b", "", "Base URL")
	}
	if strings.Contains(args, "-f,") {
		cfg.FileStoragePath = flag.String("f", "", "File storage path")
	}
	if strings.Contains(args, "-d,") {
		cfg.DataBaseDSN = flag.String("d", "", "DB connect string")
	}
	if strings.Contains(args, "-s,") {
		cfg.EnableHTTPS = flag.Bool("s", false, "Enable HTTPS")
	}
	if strings.Contains(args, "-c,") {
		cfg.Config = flag.String("c", "", "config path string")
	}
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
	return &cfg
}

func InitConfig() *config {
	params := initParams()
	cfgNoFile := initParamsWithConfig(params)
	if cfgNoFile.Config == "" {
		return cfgNoFile
	}
	err := setEnvFromJSON(cfgNoFile.Config)
	if err != nil {
		log.Print(err)
		return cfgNoFile
	}
	return initParamsWithConfig(params)
}

func setEnvFromJSON(path string) error {
	if path == "" {
		return errors.New("no json env file")
	}
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
	return err
}
