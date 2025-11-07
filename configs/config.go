package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `yaml:"app" mapstructure:"app"`
	Database DatabaseConfig `yaml:"database" mapstructure:"database"`
	JWT      JwtConfig      `yaml:"jwt" mapstructure:"jwt"`
}

type AppConfig struct {
	Port uint `mapstructure:"port"`
}

type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`
	Url    string `mapstructure:"url"`
}

type JwtConfig struct {
	SecretKey string `mapstructure:"secretKey"`
}

var Cfg *Config

func InitConfig() error {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("加载配置文件失败: %v", err)
	}

	Cfg = &Config{}

	err = viper.Unmarshal(Cfg)
	if err != nil {
		return fmt.Errorf("解析配置文件失败：%v", err)
	}

	return nil
}
