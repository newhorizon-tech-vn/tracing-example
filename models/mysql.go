package models

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DBConnection *gorm.DB
)

// Setup initializes the database instance
func InitMySQL() error {
	var err error
	// dsn := "root:@tcp(127.0.0.1:3306)/hr_auth?charset=utf8mb4&parseTime=True&loc=UTC"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s",
		viper.GetString("database.User"),     // setting.DatabaseSetting.User,
		viper.GetString("database.Password"), // setting.DatabaseSetting.Password,
		viper.GetString("database.Host"),     // setting.DatabaseSetting.Host,
		viper.GetString("database.Port"),     // setting.DatabaseSetting.Port,
		viper.GetString("database.Name"),     // setting.DatabaseSetting.Name,
		viper.GetString("database.Timezone"), // setting.DatabaseSetting.Timezone,
	)

	DBConnection, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	if err := DBConnection.Use(otelgorm.NewPlugin()); err != nil {
		return err
	}

	sqlDB, err := DBConnection.DB()
	if err != nil {
		return err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(viper.GetInt("database.MaxIdleConns"))

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(viper.GetInt("database.MaxOpenConns"))

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Duration(viper.GetInt("database.ConnMaxLifetime")) * time.Second)
	return nil
}
