package setting

import (
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/newhorizon-tech-vn/tracing-example/pkg/log"
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
	log.Info("Read config file", cfgFile)

	folder, fileName, ext, err := util.ExtractFilePath(cfgFile)
	if err != nil {
		log.Error("Extract config file failed", "file", cfgFile, "error", err)
		return err
	}
	// Setting
	viper.AddConfigPath(folder)
	viper.SetConfigName(fileName)
	viper.AutomaticEnv()
	viper.SetConfigType(ext)

	if err := viper.ReadInConfig(); err != nil {
		log.Error("EViper using config file failed", "file", viper.ConfigFileUsed(), "error", err)
		return err
	}

	//watch on config change
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Warn("Config file changed", "event", e.Name)
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

	log.Info("Read config successful")
	return nil
}
