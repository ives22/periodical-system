package conf

import (
	"github.com/spf13/viper"
)

var (
	config *Config = DefaultConfig()
)

// C 返回全局对象config
func C() *Config {
	if config == nil {
		panic("请提前加载配置...")
	}
	return config
}

// LoadConfigFromYaml 从yaml配置文件中加载配置
func LoadConfigFromYaml(confPath string) error {
	viper.SetConfigFile(confPath)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(config); err != nil {
		return err
	}

	return nil
}
