package common

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var (
	Mysql *gorm.DB
	Redis *redis.Client
)
