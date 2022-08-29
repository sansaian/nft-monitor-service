package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"nft-monitor-service/config"
	"nft-monitor-service/internal/app"
	"nft-monitor-service/internal/output"
	"nft-monitor-service/internal/usecase"
	"nft-monitor-service/pkg/logger"
	"nft-monitor-service/pkg/solana"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	//chanel for interrupt
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	// read config from env
	cfg := config.NewFromENV()
	if !cfg.IsValid() {
		log.Fatalf("config has empty required fields")
	}

	// initialize logger
	log, err := logger.NewWithConfig(cfg.Loglevel)
	if err != nil {
		log.Fatalf("failed initialize logger")
	}

	dataOutput := output.New()
	solanaClient := solana.New(log, cfg)
	nftMonitoring := usecase.New(log, solanaClient, cfg.WaitTimeBlock)
	// initialize app
	appMonitoring := app.New(cfg, log, dataOutput, nftMonitoring)

	go appMonitoring.Run(ctx)
	// get interrupt syscall call and stop graceful stop
	<-interrupt
	cancel()
}
