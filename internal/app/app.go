package app

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"

	"nft-monitor-service/config"
	"nft-monitor-service/internal/usecase"
)

type App struct {
	ctx        context.Context
	log        *logrus.Logger
	config     *config.Config
	dataOutput usecase.DataOutput
	useCase    usecase.UseCase
}

func New(config *config.Config, log *logrus.Logger, dataOutput usecase.DataOutput, useCase usecase.UseCase) *App {
	return &App{
		config:     config,
		log:        log,
		dataOutput: dataOutput,
		useCase:    useCase,
	}
}

// Run -run application for find nft in blockchain
func (app *App) Run(ctx context.Context) {
	block := app.config.StartBlock
	for {
		select {
		case <-ctx.Done():
			app.log.Infof("application shutdown. Last process block =%d ", block)
			return
		default:
			nftMetadata, err := app.useCase.GetNFTDataFromBlock(block)
			if err != nil {
				if errors.Is(err, usecase.ErrBlockNotFinalize) || errors.Is(err, usecase.ErrConnectionProblem) {
					continue
				}
				logrus.WithError(err).Fatalf("failed get nft data from solana. stop on %d block", block)
			}
			err = app.dataOutput.Output(nftMetadata)
			if err != nil {
				logrus.WithError(err).Fatalf("failed output data information.stop on %d block", block)
			}
			block++
		}
	}
}
