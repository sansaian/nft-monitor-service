package usecase

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"

	"nft-monitor-service/internal/entities"
)

type NFTMonitoring struct {
	log        *logrus.Logger
	blockchain Blockchain
	timeWait   time.Duration
	lastBlock  uint64
}

func New(log *logrus.Logger, blockchain Blockchain, timeWait int) *NFTMonitoring {
	return &NFTMonitoring{
		log:        log,
		blockchain: blockchain,
		timeWait:   time.Duration(timeWait) * time.Second,
		lastBlock:  0,
	}
}

func (m *NFTMonitoring) GetNFTDataFromBlock(numBlock uint64) ([]*entities.NFTTokenData, error) {
	isFinalized, err := m.isFinalizedBlock(numBlock)
	if err != nil {
		if errors.Is(err, ErrConnectionProblem) {
			m.log.Debugf("num block %d. connection problem try to reconect", numBlock)
			m.wait()
			return nil, ErrConnectionProblem
		}
		return nil, fmt.Errorf("failed get last finalized block %w", err)
	}
	if !isFinalized {
		m.log.Debugf("num block %d not finalize.Need wait", numBlock)
		m.wait()
		return nil, ErrBlockNotFinalize
	}
	nftMetadata, err := m.blockchain.ParseBlock(numBlock)
	if err != nil {
		return nil, fmt.Errorf("failed parse block %w", err)
	}
	return nftMetadata, nil
}

func (m *NFTMonitoring) isFinalizedBlock(numBlock uint64) (bool, error) {
	var err error
	if m.lastBlock <= numBlock {
		m.lastBlock, err = m.blockchain.GetLastBlock()
		if err != nil {
			return false, fmt.Errorf("failed get last block %w", err)
		}
		return false, nil
	}

	return true, nil
}

func (m *NFTMonitoring) wait() {
	m.log.Debugf("need sleeping %s sec", m.timeWait.String())
	time.Sleep(m.timeWait)
}
