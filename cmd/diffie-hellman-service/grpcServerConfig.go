package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type GrpcServerConfig struct {
	TlsEnable     bool   `json:"tlsEnable"`
	Port          uint32 `json:"port"`
	CertFile      string `json:"certFile"`
	KeyFile       string `json:"keyFile"`
	TmpfsLimitMB  int32  `json:"tmpfsLimitMB"`
	FileChunkSize int    `json:"fileChunkSize"`
}

type ServerConfig struct {
	AppName          string           `json:"appName"`
	GrpcServerConfig GrpcServerConfig `json:"gRPCServerConfig"`
}

var configData ServerConfig

func initConfig(configPath string) (*ServerConfig, error) {

	configFile, err := os.Open(configPath)

	if err != nil {
		log.Printf("Error: %v", err)

		configFile, err = os.Open("./config-server.json")

		if err != nil {
			log.Printf("Error: %v", err)
			return &configData, err
		}
	}
	log.Print("Successfully Opened config file:", configPath)

	defer configFile.Close()

	byteValue, _ := ioutil.ReadAll(configFile)

	err = json.Unmarshal([]byte(byteValue), &configData)

	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}

	log.Print("gRPCServer SSL enable:", configData.GrpcServerConfig.TlsEnable)
	log.Print("gRPCServer port: ", configData.GrpcServerConfig.Port)
	log.Print("gRPCServer cert file: ", configData.GrpcServerConfig.CertFile)
	log.Print("gRPCServer key file: ", configData.GrpcServerConfig.KeyFile)
	log.Print("gRPCServer cached file in tmpfs limit: ", configData.GrpcServerConfig.TmpfsLimitMB)
	log.Print("gRPCServer transfer file chunk size: ", configData.GrpcServerConfig.FileChunkSize)

	return &configData, err
}
