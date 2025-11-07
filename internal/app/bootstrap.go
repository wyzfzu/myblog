package app

import (
	"fmt"

	"github.com/blog/configs"
	"github.com/blog/internal/api"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Start() {
	err := configs.InitConfig()
	if err != nil {
		log.Errorf("初始化配置异常：%v", err)
		return
	}

	err = InitAll(configs.Cfg)
	if err != nil {
		log.Errorf("模块初始化异常：%v", err)
		return
	}

	ctx := gin.New()
	api.SetupRoutes(ctx, DB, configs.Cfg)

	err = ctx.Run(fmt.Sprintf(":%d", configs.Cfg.App.Port))
	if err != nil {
		log.Errorf("服务启动异常：%v", err)
		return
	}

	log.Info("服务启动成功！")
}
