package main

import "flag"

func initParams(cfg config) {
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "Server address")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "Base URL")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "File storage path")
	flag.StringVar(&cfg.BataBaseDSN, "d", cfg.BataBaseDSN, "DB connect string")
	flag.Parse()
}
