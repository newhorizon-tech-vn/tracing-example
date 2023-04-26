package setting

import (
	"fmt"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/newhorizon-tech-vn/tracing-example/pkg/util"
	"github.com/spf13/viper"
)

type RedisConfig struct {
	Address           string
	MaxRetries        int
	Password          string
	DefaultExpireTime time.Duration
}

type ServiceConfig struct {
	Redis *RedisConfig
}

var Setting *ServiceConfig

func InitSetting() error {
	cfgFile := os.Getenv("CONFIG_PATH")
	if len(cfgFile) == 0 {
		cfgFile = "./conf/config.dev.toml"
	}

	folder, fileName, ext, err := util.ExtractFilePath(cfgFile)
	if err != nil {
		fmt.Printf("ERROR: Extract config file failed file: %s error: %s \n", cfgFile, err.Error())
		return err
	}
	// Setting
	viper.AddConfigPath(folder)
	viper.SetConfigName(fileName)
	viper.AutomaticEnv()
	viper.SetConfigType(ext)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("ERROR: viper using config file failed file: %s error: %s \n", viper.ConfigFileUsed(), err.Error())
		return err
	}

	//watch on config change
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("WARN: Config file changed event: %s \n", e.Name)
	})

	viper.SetDefault("redis.maxRetry", 5)
	Setting = &ServiceConfig{
		Redis: &RedisConfig{
			Address:           viper.GetString("redis.address"),
			MaxRetries:        viper.GetInt("redis.maxRetry"),
			Password:          viper.GetString("redis.password"),
			DefaultExpireTime: time.Second * (3600 * 24), // one day
		},
	}

	fmt.Printf("INFO: read config successful \n")
	return nil
}
