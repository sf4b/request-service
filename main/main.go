package main

import (
	"context"
	"flag"
	"log"
	"service/config"
	"service/infra"
)

func main() {
	configFile := flag.String("c", "config.template.yaml", "configuration file")
	envFile := flag.String("e", "", "env file")
	flag.Parse()

	var cfg config.Config
	if err := config.Parse(*configFile, *envFile, &cfg); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	builder := infra.NewBuilder(ctx, &cfg)

	info := builder.GetInfo()

	log.Printf("Service: %s  of version %s is started", info.Name, info.Version)

}
