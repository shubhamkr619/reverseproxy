package database

import (
	"fmt"
	"github.com/revproxy/src/config"
	"github.com/spf13/viper"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		panic(err)
	}
	time.Local = loc
}

func UseDB() *gorm.DB {
	config, err := ReadConfig()
	if err != nil {
		fmt.Errorf("error received while reading the configuration ", err)
	}
	if db == nil {
		db = initDB(config)
	}
	return db
}

func ReadConfig() (config.Configurations, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	var configuration config.Configurations
	if err := viper.ReadInConfig(); err != nil {
		return configuration, fmt.Errorf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}
	return configuration, nil
}

func initDB(configurations config.Configurations) *gorm.DB {
	config, err := ReadConfig()
	if err != nil {
		fmt.Println("Cant read the config ", err)
	}
	dbUser := config.Database.DBUser
	dbPass := config.Database.DBPass
	dbHost := config.Database.DbHost
	dbPort := config.Database.DbPort
	dbName := config.Database.DBName
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
