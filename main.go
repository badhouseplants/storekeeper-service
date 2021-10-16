package main

import (
	"fmt"

	"github.com/droplez/droplez-uploader/pkg/server"
	"github.com/droplez/droplez-uploader/tools/constants"
	"github.com/spf13/viper"
)

func init() {
	// Set default values
	viper.SetDefault(constants.DroplezUploaderPort, "9090")
	viper.SetDefault(constants.DroplezUploaderMode, "dev")

	viper.SetDefault(constants.MinioBucket, "droplez-dev")
	viper.SetDefault(constants.MinioEndpoint, "localhost:9000")
	viper.SetDefault(constants.MinioAccessKeyID, "minio")
	viper.SetDefault(constants.MinioSecretAccessKey, "minio123")

	// Read environment
	viper.AutomaticEnv()
}

func main() {
	fmt.Println(viper.GetString(constants.DroplezUploaderPort))
	fmt.Println(viper.GetString(constants.MinioBucket))
	if err := server.Serve(); err != nil {
		fmt.Println(err)
	}
}
