package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
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

	http.HandleFunc("/healthz", healthz)
	http.HandleFunc("/webhook", webhook)

	listenAddress := ":" + viper.GetString("http_port")
	log.Printf("listening on: %v", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
