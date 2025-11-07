package app

import (
	"fmt"

	"github.com/blog/configs"
)

func InitAll(cfg *configs.Config) error {
	err := InitSqlite3(cfg)
	if err != nil {
		return fmt.Errorf("初始化sqlite3失败: %v", err)
	}

	return nil
}
