package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/Go-Short/api/service"
)

const (
	defaultConfigPath = "./conf/conf.toml"
)

func main() {
	err := readConfig(defaultConfigPath)
	logger := service.GetLogger()

	if err != nil {
		fmt.Printf("read config error: %s", err)
		return
	}

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	ginPort := fmt.Sprintf(":%d", viper.GetInt64("server.port"))
	logger.Info("Starting server")

	r.Run(ginPort)
}

func readConfig(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	viper.SetConfigType("toml")

	viper.SetDefault("server.host", "127.0.0.1")
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("redis.port", "6379")
	return viper.ReadConfig(f)
}
