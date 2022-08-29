package solana

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/near/borsh-go"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/metaplex/tokenmeta"
	"github.com/sirupsen/logrus"
	"io"
	"strings"
	"time"

	"nft-monitor-service/config"
	"nft-monitor-service/internal/entities"
	"nft-monitor-service/internal/usecase"
)

type ClientSolana struct {
	cfg       *config.Config
	log       *logrus.Logger
	client    *client.Client
	programID solana.PublicKey
	timeout   time.Duration
}

func New(log *logrus.Logger, cfg *config.Config) *ClientSolana {
	url := fmt.Sprintf("%s/%s", cfg.Solana.URL, cfg.Solana.ApiKey)
	c := client.NewClient(url)
	return &ClientSolana{
		cfg:       cfg,
		log:       log,
		client:    c,
		programID: cfg.Solana.ProgramID,
		timeout:   time.Duration(cfg.Solana.Timeout) * time.Second}
}

func (sol *ClientSolana) GetLastBlock() (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sol.timeout)
	defer cancel()
	height, err := sol.client.RpcClient.GetSlot(ctx)
	if err != nil {
		if strings.Contains(err.Error(), errorMsgInternal) || strings.Contains(err.Error(), errorMsgToManyRequest) {
			return 0, usecase.ErrConnectionProblem
		}
		return 0, err
	}
	return height.Result, nil
}

func (sol *ClientSolana) ParseBlock(numBlock uint64) ([]*entities.NFTTokenData, error) {
	sol.log.Debugf("start parsing block %d", numBlock)
	defer sol.log.Debugf("finished parsing block %d", numBlock)
	ctx, cancel := context.WithTimeout(context.Background(), sol.timeout)
	defer cancel()
	block, err := sol.client.GetBlock(ctx, numBlock)
	if err != nil {
		if strings.Contains(err.Error(), errorMsgSkippedBlock) {
			sol.log.Debugf("Slot %d was skipped, or missing in long-term storage", numBlock)
			return nil, nil
		}
		if strings.Contains(err.Error(), errorMsgInternal) || strings.Contains(err.Error(), errorMsgToManyRequest) {
			return nil, usecase.ErrConnectionProblem
		}
		return nil, fmt.Errorf("failed get block from solana %w", err)
	}
	var nftMetas []*entities.NFTTokenData
	for _, value := range block.Transactions {
		for _, account := range value.Transaction.Message.Accounts {
			if sol.programID.String() == account.String() {
				for _, address := range value.Meta.PostTokenBalances {
					metadata, err := sol.getMetadata(address.Mint)
					if err != nil {
						return nil, err
					}
					if metadata == nil {
						continue
					}
					nftMetas = append(nftMetas, metadata.ConvertToEntity(address.Mint, numBlock))
				}
			}

		}
	}

	return nftMetas, nil

}

func (sol *ClientSolana) getMetadata(addressNFT string) (*ResponseNFTMeta, error) {

	metadataAccount, err := tokenmeta.GetTokenMetaPubkey(common.PublicKeyFromString(addressNFT))
	if err != nil {
		return nil, fmt.Errorf("failed get token meta pubkey %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), sol.timeout)
	defer cancel()
	accountInfo, err := sol.client.GetAccountInfo(ctx, metadataAccount.ToBase58())
	if err != nil {
		if strings.Contains(err.Error(), errorMsgInternal) || strings.Contains(err.Error(), errorMsgToManyRequest) {
			return nil, usecase.ErrConnectionProblem
		}
		return nil, fmt.Errorf("failed to get accountInfo, err: %w", err)
	}
	url, err := parseURI(&accountInfo)
	if err != nil {
		return nil, err
	}
	metadata, err := loadFromURI(url)
	if err != nil {

		return nil, fmt.Errorf("failed load metadata from uri %w", err)
	}
	return metadata, nil
}

// parseURI-method parses the uri from AccountInfo and results in a human-readable form.
func parseURI(accountInfo *client.AccountInfo) (string, error) {
	var sanitizedURI string
	if accountInfo.Data != nil {
		mm := new(MetaPlexMeta)
		err := borsh.Deserialize(mm, accountInfo.Data)
		if err != nil {
			return "", fmt.Errorf("failed deserialize url from accountInfo %w", err)
		}
		uri := mm.Data.Uri
		sanitizedURI = strings.Replace(uri, "\u0000", "", -1)
	}
	return sanitizedURI, nil
}

// loadFromURI json-metadata from nft-url
func loadFromURI(uri string) (*ResponseNFTMeta, error) {
	if uri == "" {
		return nil, nil
	}
	resp, err := retryablehttp.Get(uri)
	if err != nil {
		return nil, fmt.Errorf("failed request to url %s ,error- %w", uri, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read json body %w", err)
	}
	nftMeta := &ResponseNFTMeta{}
	err = json.Unmarshal(body, nftMeta)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshal metadata from url %s ,error- %w", uri, err)
	}
	nftMeta.URL = uri
	return nftMeta, nil
}
