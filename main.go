package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/promalert")
	viper.SetConfigName("config")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	viper.AutomaticEnv()
	viper.SetEnvPrefix("promalert")

	r := gin.New()
	r.Use(gin.LoggerWithWriter(gin.DefaultWriter, "/healthz"))
	r.Use(gin.Recovery())

	r.GET("/healthz", healthz)
	r.POST("/webhook", webhook)

	err = r.Run(":" + viper.GetString("http_port"))
	panic(fmt.Errorf("Cant start web server: %s \n", err))
}
