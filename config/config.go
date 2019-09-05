package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

var initOnce sync.Once
var gcf atomic.Value

type AppConfig struct {
	Web struct {
		Mode          string `mapstructure:"mode"`
		ListenAddress string `mapstructure:"listen_address"`
	} `mapstructure:"web"`
	PieceSize int64 `mapstructure:"piece_size"`
	SessionExpire time.Duration `mapstructure:"session_expire"`
	SessionCheckInterval time.Duration `mapstructure:"session_check_interval"`
	DownloadDir string `mapstructure:"download_dir"`
}

func InitConfig(role string) error {
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/p2p-sharer")

	fname := fmt.Sprintf("config.%s", role)
	log.Println("config file name", fname)
	viper.SetConfigName(fname)

	initOnce.Do(func() {
		err := viper.ReadInConfig()
		if err == nil{
			log.Println("Using config file:", viper.ConfigFileUsed())
		} else {
			log.Fatal("can't find config file", err)
		}

		var cf = AppConfig{}
		viper.Unmarshal(&cf)
		gcf.Store(&cf)
	})

	log.Println("SessionExpire", GetOpts().SessionExpire)
	log.Println("SessionCheckInterval", GetOpts().SessionCheckInterval)
	log.Println("PieceSize", GetOpts().PieceSize)
	log.Println("Web", GetOpts().Web)
	return nil
}

func GetOpts() *AppConfig {
	return gcf.Load().(*AppConfig)
}
