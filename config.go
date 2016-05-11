package main

import (
	"log"
	"runtime"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	LISTEN_PORT    = 8080
	LISTEN_ADDRESS = "0.0.0.0"
)

func init() {
	viper.SetConfigName("slackhub")
	viper.AddConfigPath("/etc/slackhub")
	viper.AddConfigPath("/etc/default")
	viper.AddConfigPath(".")
	viper.SetConfigType("json")

	viper.SetDefault("max_workers", runtime.NumCPU())
	viper.SetDefault("listen_port", LISTEN_PORT)
	viper.SetDefault("listen_address", LISTEN_ADDRESS)
	viper.SetDefault("webhook_url", "")

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("[WARN] Config file: %s", err)
	}

	viper.SetEnvPrefix("SLACKHUB")
	viper.AutomaticEnv()

	pflag.Int("max_workers", runtime.NumCPU(), "Maximum number of workers to use")
	pflag.Int("listen_port", LISTEN_PORT, "Port to listen on")
	pflag.String("listen_address", LISTEN_ADDRESS, "Interface address to listen on")
	pflag.String("webhook_url", "", "Slack WebHook URL")

	viper.BindPFlag("max_workers", pflag.Lookup("max_workers"))
	viper.BindPFlag("listen_port", pflag.Lookup("listen_port"))
	viper.BindPFlag("listen_address", pflag.Lookup("listen_address"))
	viper.BindPFlag("webhook_url", pflag.Lookup("webhook_url"))

	pflag.Parse()
}

func MaxWorkers() int {
	mw := viper.GetInt("max_workers")
	if mw <= 0 {
		log.Fatal("No workers available")
	}
	return mw
}

func ListenPort() int {
	lp := viper.GetInt("listen_port")
	if lp <= 0 {
		log.Fatal("Invalid port")
	}
	return lp
}

func ListenAddress() string {
	la := viper.GetString("listen_address")
	if la == "" {
		log.Fatal("Invalid address")
	}
	return la
}

func WebhookUrl() string {
	wh := viper.GetString("webhook_url")
	if wh == "" {
		log.Fatal("Invalid webhook URL")
	}
	return wh
}

func DisplayAllProperties() {
	log.Printf("Max workers: %d", MaxWorkers())
	log.Printf("Listen address: %s", ListenAddress())
	log.Printf("Listen port: %d", ListenPort())
	log.Printf("Slack webhook URL: %s", WebhookUrl())
}
