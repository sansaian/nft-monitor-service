package solana

import (
	"github.com/gagliardetto/solana-go"
	"strconv"

	"nft-monitor-service/internal/entities"
)

type MetaPlexMeta struct {
	Key             byte
	UpdateAuthority solana.PublicKey
	Mint            solana.PublicKey
	Data            MetaPlexData
}

type MetaPlexData struct {
	Name   string
	Symbol string
	Uri    string
}

type ResponseNFTMeta struct {
	Name                 string `json:"name"`
	URL                  string
	Symbol               string `json:"symbol"`
	Description          string `json:"description"`
	SellerFeeBasisPoints int    `json:"seller_fee_basis_points"`
	Image                string `json:"image"`
	ExternalURL          string `json:"external_url"`
	Properties           struct {
		Files []struct {
			URI  string `json:"uri"`
			Type string `json:"type"`
		} `json:"files"`
		Creators []struct {
			Address string `json:"address"`
			Share   int    `json:"share"`
		} `json:"creators"`
	} `json:"properties"`
}

func (resp *ResponseNFTMeta) ConvertToEntity(tokenID string, blockNumber uint64) *entities.NFTTokenData {
	blockNum := strconv.FormatUint(blockNumber, 10)
	var creators []entities.Creator
	for _, creator := range resp.Properties.Creators {
		creators = append(creators, creator)
	}
	return &entities.NFTTokenData{
		TokenID:     tokenID,
		BlockNumber: blockNum,
		Content: entities.Content{
			Data: entities.Data{
				Name:                 resp.Name,
				Symbol:               resp.Symbol,
				URI:                  resp.URL,
				SellerFeeBasisPoints: resp.SellerFeeBasisPoints,
				Creators:             creators,
			},
		},
	}
}
