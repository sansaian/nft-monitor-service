package output

import (
	"encoding/json"
	"fmt"

	"nft-monitor-service/internal/entities"
)

type STDOut struct {
}

// New - create tools for output data about token
func New() *STDOut {
	return &STDOut{}
}

// Output - print to stdout data about nft
func (o *STDOut) Output(data []*entities.NFTTokenData) error {
	if len(data) == 0 {
		return nil
	}
	for _, nftMeta := range data {
		if nftMeta == nil {
			return fmt.Errorf("nft metadata is nil")
		}
		b, err := json.Marshal(nftMeta)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	}

	return nil
}
