package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fansqz/fancode-backend/common"
	config2 "github.com/fansqz/fancode-backend/common/config"
	"github.com/fansqz/fancode-backend/common/logger"
	ratelimiter "github.com/fansqz/fancode-backend/common/rate_limiter"
	"github.com/fansqz/fancode-backend/models/po"
	"github.com/gin-gonic/gin"
)

func newApp(engine *gin.Engine, config *config2.AppConfig) *http.Server {
	srv := &http.Server{
		Addr:    config.Port,
		Handler: engine,
	}
	return srv
}

func main() {
	//加载配置
	conf := config2.InitSetting()
	// 初始化日志
	logger.InitLogger(context.Background(), conf.LoggerConfig)

	//连接数据库
	common.InitMysql(conf.MySqlConfig)

	//连接redis
	if err := common.InitRedis(conf.RedisConfig); err != nil {
		fmt.Println("redis连接失败")
	}

	// 初始化限流器
	ratelimiter.InitSentinel()

	// 模型绑定
	err := common.Mysql.AutoMigrate(
		&po.SysApi{},
		&po.SysMenu{},
		&po.SysRole{},
		&po.SysUser{},
		&po.UserCode{},
		&po.VisualDocument{},
		&po.VisualDocumentCode{},
		&po.VisualDocumentBank{},
		&po.UserSavedCode{},
	)
	if err != nil {
		panic(err)
	}

	//注册路由
	srv, err := initApp(conf)
	if err != nil {
		logger.WithCtx(context.Background()).Errorf("init app error, err = %v", err)
		panic(err)
	}

	if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.WithCtx(context.Background()).Errorf("listen serve error, err = %v", err)
		panic(err)
	}
}
