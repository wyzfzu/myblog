package app

import (
	"fmt"

	"github.com/blog/configs"
	"github.com/blog/internal/models"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitSqlite3(cfg *configs.Config) error {
	db, err := gorm.Open(sqlite.Open(cfg.Database.Url), &gorm.Config{})
	if err != nil {
		log.Errorf("初始化sqlite失败: %v", err)
		return err
	}

	sqlDB, err := db.DB()
	if err == nil {
		if pingErr := sqlDB.Ping(); pingErr != nil {
			log.Errorf("数据库连接测试异常: %v", pingErr)
			return fmt.Errorf("数据库连接测试异常: %v", pingErr)
		}
	}
	DB = db

	initTable(db)

	return nil
}

func initTable(db *gorm.DB) {
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Post{})
	db.AutoMigrate(&models.Comment{})
}
