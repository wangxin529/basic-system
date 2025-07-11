package mysql

import (
	"elevate-hub/db/conf"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/sharding"
	"log"
	"os"
	"time"
)

func InitMysql(config conf.Mysql, middleware *sharding.Sharding) *gorm.DB {

	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", //&timeout=%s",
		config.User, config.Passwd, config.Host, config.Port, config.Database)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	mysqlConfig := &gorm.Config{}
	if config.Logger == true {
		mysqlConfig.Logger = newLogger
	}

	db, err := gorm.Open(mysql.Open(dns), mysqlConfig)

	if middleware != nil && err == nil {
		err = db.Use(middleware)
	}
	if err != nil {
		panic(err)
	}
	fmt.Println("mysql init success.....")
	return db
}
