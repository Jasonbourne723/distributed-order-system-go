package global

import (
	"distributed-order-system-go/pkg/config"
	"distributed-order-system-go/pkg/zookeeper"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var App = new(Application)

type Application struct {
	ConfigViper       *viper.Viper
	Config            config.Configuration
	DB                *gorm.DB
	Redis             *redis.Client
	DistributedLocker *zookeeper.ZookeeperLock
}
