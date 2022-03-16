package model

import (
	"fmt"
	_config "goblog/pkg/config"
	"goblog/pkg/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
	var err error
	_DNS := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&loc=Local",
		_config.GetString("database.mysql.username"),
		_config.GetString("database.mysql.password"),
		_config.GetString("database.mysql.host"),
		_config.GetString("database.mysql.port"),
		_config.GetString("database.mysql.database"),
		_config.GetString("database.mysql.charset"))
	config := mysql.New(mysql.Config{
		DSN: _DNS,
	})
	//准备数据库连接池

	var level gormlogger.LogLevel
	if _config.GetBooler("app.debug") {
		level = gormlogger.Warn
	} else {
		level = gormlogger.Error
	}
	DB, err = gorm.Open(config, &gorm.Config{
		Logger: gormlogger.Default.LogMode(level),
	})
	logger.LogError(err)
	return DB

}
