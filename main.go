package main

import (
	"log"

	"github.com/droplez/droplez-storekeeper-service/pkg/constants"
	"github.com/droplez/droplez-storekeeper-service/pkg/server"
	"github.com/spf13/viper"
)

func init() {
	// Set default values
	viper.SetDefault(constants.ConstDroplezUploaderPort, "9090")
	viper.SetDefault(constants.ConstDroplezUploaderMode, "dev")
	viper.SetDefault(constants.ConstMinioBucket, "droplez-dev")
	viper.SetDefault(constants.ConstMinioEndpoint, "localhost:9000")
	viper.SetDefault(constants.ConstMinioAccessKeyID, "minio")
	viper.SetDefault(constants.ConstMinioSecretAccessKey, "minio123")

	// Read environment
	viper.AutomaticEnv()
}

func main() {
	if err := server.Serve(); err != nil {
		log.Panic(err)
	}
}
