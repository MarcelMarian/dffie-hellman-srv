package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var DefaultConfigFile = "/config/config-adapter.json"
var configDataPtr *ServerConfig
var seqNo int32 = 0
var err error

func initSignalHandle() {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		// Run Cleanup
		fmt.Println("Receive: get exit signal, exit now.")
		os.Exit(1)
	}()
}

func main() {

	initSignalHandle()
	fmt.Println("Diffie-Hellman gRPC server application")

	// Initialize configuration
	for {
		configDataPtr, err = initConfig(DefaultConfigFile)
		if err == nil {
			break
		}
		time.Sleep(3 * time.Second)
	}

	// Initialize gRPC interface
	startGRPCSever()
	fmt.Println("Diffie-Hellman gRPC server application DONE")

}
