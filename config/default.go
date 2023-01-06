package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	DSNS []string

	DB_USER string   `mapstructure:"DB_USER"`
	DB_PASS string   `mapstructure:"DB_PASS"`
	DB_NAME string   `mapstructure:"DB_NAME"`
	DB_HOST string   `mapstructure:"DB_HOST"`
	PORTS   []string `mapstructure:"PORTS"`
}

func LoadConfig(path string) (Config, error) {
	config := Config{}
	viper.AddConfigPath(path)
	//viper.SetConfigType("env")
	//viper.SetConfigName("")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}
	err = viper.Unmarshal(&config)
	config.create_dsns()
	return config, nil
}
func (c *Config) create_dsns() {
	dsns := make([]string, 0)
	for _, port := range c.PORTS {
		psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			c.DB_HOST, port, c.DB_USER, c.DB_PASS, c.DB_NAME)
		dsns = append(dsns, psqlInfo)
	}
	c.DSNS = dsns
}
