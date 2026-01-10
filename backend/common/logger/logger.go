package logger

import (
	"context"
	"fmt"
	"github.com/Graylog2/go-gelf/gelf"
	"github.com/bwmarrin/snowflake"
	"github.com/fansqz/fancode-backend/common/config"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
)

const (
	LOG_ID_KEY   = "log_id"
	USER_ID_KEY  = "user_id"
	DEBUG_ID_KEY = "debug_id"
)

var snowflakeNode *snowflake.Node

func InitLogger(ctx context.Context, config *config.LoggerConfig) {
	var err error
	// 初始化雪花算法，用于生成日志id
	if snowflakeNode, err = snowflake.NewNode(1); err != nil {
		WithCtx(ctx).Errorf("[InitLogger] init logger error err = %v", err)
		log.Fatalf("snowflake.NewNode: %s", err)
		return
	}

	// 初始化日志收集
	if config.Type == "graylog" {
		initGraylog(config)
	}
}

// initGraylog 初始化graylog
func initGraylog(config *config.LoggerConfig) {
	gelfWriter, err := gelf.NewWriter(fmt.Sprintf("%s:%s", config.Host, config.Port))
	if err != nil {
		log.Fatalf("gelf.NewWriter: %s", err)
		return
	}
	// 设置日志收集
	logrus.SetOutput(io.MultiWriter(os.Stderr, gelfWriter))
}

func WithCtx(ctx context.Context) *logrus.Entry {
	r := logrus.WithContext(ctx)
	if logID, ok := ctx.Value(LOG_ID_KEY).(string); ok {
		r = r.WithField(LOG_ID_KEY, logID)
	}
	if userID, ok := ctx.Value(USER_ID_KEY).(string); ok {
		r = r.WithField(USER_ID_KEY, userID)
	}
	if debugID, ok := ctx.Value(DEBUG_ID_KEY).(string); ok {
		r = r.WithField(DEBUG_ID_KEY, debugID)
	}
	return r
}

func SetLogID(ctx context.Context, logID string) context.Context {
	return context.WithValue(ctx, LOG_ID_KEY, logID)
}

func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, USER_ID_KEY, userID)
}

// SetDebugID 设置调试id
func SetDebugID(ctx context.Context, debugID string) context.Context {
	return context.WithValue(ctx, DEBUG_ID_KEY, debugID)
}

// GenerateLogID 随机生成logID
func GenerateLogID() string {
	// Generate a snowflake ID.
	id := snowflakeNode.Generate()
	return id.String()
}
