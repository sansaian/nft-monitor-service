package usecase

import (
	"nft-monitor-service/internal/entities"
)

type DataOutput interface {
	Output(data []*entities.NFTTokenData) error
}

type Blockchain interface {
	GetLastBlock() (uint64, error)
	ParseBlock(numBlock uint64) ([]*entities.NFTTokenData, error)
}

type UseCase interface {
	GetNFTDataFromBlock(numBlock uint64) ([]*entities.NFTTokenData, error)
}
