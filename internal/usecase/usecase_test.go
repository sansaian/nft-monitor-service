package usecase

import (
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"reflect"
	"testing"
	"time"

	"nft-monitor-service/internal/entities"
	"nft-monitor-service/internal/usecase/mock"
	"nft-monitor-service/pkg/logger"
)

func TestNFTMonitoring_GetNFTDataFromBlock(t *testing.T) {
	type fields struct {
		log        *logrus.Logger
		blockchain Blockchain
		timeWait   time.Duration
		lastBlock  uint64
	}
	type args struct {
		numBlock uint64
	}

	logger := logger.New()
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockBlockchain := mock.NewMockBlockchain(ctl)

	var blockFromUser uint64 = 148725696
	expectedNFTMetadata := &entities.NFTTokenData{
		TokenID:     "testTokenID",
		BlockNumber: "148725696",
		Content:     entities.Content{},
	}
	mockBlockchain.EXPECT().ParseBlock(blockFromUser).Return([]*entities.NFTTokenData{expectedNFTMetadata}, nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*entities.NFTTokenData
		wantErr bool
	}{
		{
			name: "happy path",
			fields: fields{
				log:        logger,
				blockchain: mockBlockchain,
				timeWait:   1,
				lastBlock:  blockFromUser + 3,
			},
			args: args{
				numBlock: blockFromUser,
			},
			want:    []*entities.NFTTokenData{expectedNFTMetadata},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &NFTMonitoring{
				log:        tt.fields.log,
				blockchain: tt.fields.blockchain,
				timeWait:   tt.fields.timeWait,
				lastBlock:  tt.fields.lastBlock,
			}
			got, err := m.GetNFTDataFromBlock(tt.args.numBlock)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNFTDataFromBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNFTDataFromBlock() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNFTMonitoring_isFinalizedBlock(t *testing.T) {
	logger := logger.New()
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockBlockchain := mock.NewMockBlockchain(ctl)

	var blockFromUser uint64 = 148725696
	var lastFinalizedBlock uint64 = 148825696
	mockBlockchain.EXPECT().GetLastBlock().Return(lastFinalizedBlock, nil)

	type fields struct {
		log        *logrus.Logger
		blockchain Blockchain
		timeWait   time.Duration
		lastBlock  uint64
	}
	type args struct {
		numBlock uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "load last block",
			fields: fields{
				log:        logger,
				blockchain: mockBlockchain,
				timeWait:   1,
				lastBlock:  0,
			},
			args: args{
				numBlock: blockFromUser,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "not need load last block",
			fields: fields{
				log:        logger,
				blockchain: nil,
				timeWait:   1,
				lastBlock:  blockFromUser + 3,
			},
			args: args{
				numBlock: blockFromUser,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &NFTMonitoring{
				log:        tt.fields.log,
				blockchain: tt.fields.blockchain,
				timeWait:   tt.fields.timeWait,
				lastBlock:  tt.fields.lastBlock,
			}
			got, err := m.isFinalizedBlock(tt.args.numBlock)
			if (err != nil) != tt.wantErr {
				t.Errorf("isFinalizedBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("isFinalizedBlock() got = %v, want %v", got, tt.want)
			}
		})
	}
}
