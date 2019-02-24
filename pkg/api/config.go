package api

import (
	"github.com/spf13/viper"
)

func loadConfig(filename string) error {
	viper.SetConfigName(filename)
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}
