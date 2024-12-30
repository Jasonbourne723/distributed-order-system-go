package global

import (
	"distributed-order-system-go/pkg/config"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var App = new(Application)

type Application struct {
	ConfigViper *viper.Viper
	Config      config.Configuration
	DB          *gorm.DB
}
