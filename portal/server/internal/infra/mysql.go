package infra

import (
	"fmt"
	"portal/internal/config"
	"portal/internal/pkg/log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// MysqlDB 全局mysql连接对象
var MysqlDB *gorm.DB

// InitMysql 获取mysql连接对象
func InitMysql(cfg config.Config, mysqlConf config.MysqlConfig) (*gorm.DB, error) {
	var err error

	MysqlDB, err = gorm.Open(mysql.Open(mysqlConf.Dsn), &gorm.Config{
		Logger: log.NewGormLogger(cfg.Debug),
	})
	if err != nil {
		return nil, fmt.Errorf("连接mysql失败: %w", err)
	}

	sqlDB, err := MysqlDB.DB()
	if err != nil {
		return nil, fmt.Errorf("打开mysql连接失败: %w", err)
	}
	sqlDB.SetMaxIdleConns(mysqlConf.MaxIdleConns)
	sqlDB.SetMaxOpenConns(mysqlConf.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(mysqlConf.ConnMaxLifetime)

	//MysqlDB = MysqlDB.Unscoped() // 关闭软删除
	return MysqlDB, nil
}
