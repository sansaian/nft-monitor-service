package config

import (
	"github.com/gagliardetto/solana-go"
	"github.com/spf13/viper"
)

type Config struct {
	Loglevel      string
	StartBlock    uint64 // start monitoring blockchain from this block
	WaitTimeBlock int    // how long to wait for a new final block, time in seconds
	Solana        Solana
}

type Solana struct {
	Timeout   int // time in seconds
	URL       string
	ApiKey    string
	ProgramID solana.PublicKey
}

// NewFromENV load env params from environment variables
func NewFromENV() *Config {
	viper.AutomaticEnv()
	viper.SetDefault("LOG_LEVEL", "info")

	return &Config{
		Loglevel:      viper.GetString("LOG_LEVEL"),
		StartBlock:    viper.GetUint64("START_BLOCK"),
		WaitTimeBlock: viper.GetInt("WAIT_TIME_BLOCK"),
		Solana: Solana{
			Timeout:   viper.GetInt("TIMEOUT"),
			URL:       viper.GetString("URL"),
			ApiKey:    viper.GetString("API_KEY"),
			ProgramID: solana.MustPublicKeyFromBase58(viper.GetString("PROGRAM_ID")),
		},
	}
}

// IsValid - validate required params
func (cfg *Config) IsValid() bool {
	if cfg.StartBlock > 0 && cfg.Solana.Timeout > 0 &&
		len(cfg.Solana.ApiKey) > 0 && len(cfg.Solana.URL) > 0 &&
		cfg.WaitTimeBlock > 0 {
		return true
	}
	return false
}
