/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gomodule/redigo/redis"
)

func pingRedis(conn redis.Conn) {
	_, err := conn.Do("PING")
	if err != nil {
		log.Println("Can't connect to the Redis database")
	}
	log.Println("PING ok")
}

func connectToRedis(hostname string) (redis.Conn, error) {
	conn, err := redis.Dial("tcp", hostname+":6379")
	return conn, err
}

func main() {
	// Get the host names
	redishost := os.Getenv("REDIS_HOST")
	log.Println("REDIS_HOST", redishost)
	miniohost := os.Getenv("MINIO_HOST")
	log.Println("MINIO_HOST", miniohost)
	mqtthost := os.Getenv("MQTT_HOST")
	log.Println("MQTT_HOST", mqtthost)

	// Setup stop signal handling
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		log.Printf("caught sig: %+v", sig)
		os.Exit(0)
	}()

	// Connect to Redis and PING it if all is well
	redisconn, err := connectToRedis(redishost)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Connected to Redis")
		pingRedis(redisconn)
		defer redisconn.Close()
	}

	// Initialise Edge Agent Client
	initClient()

	// Get list of applications from Edge Agent
	var applications []string
	applications, err = queryApplications()
	if err == nil {
		for i := 0; i < len(applications); i++ {
			log.Println(applications[i])
		}
	}

	// Hang around until stopped
	for {
		time.Sleep(1 * time.Second)
	}
}
