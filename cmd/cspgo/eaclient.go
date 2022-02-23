package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
)

var restClient = resty.New()
var unixSocket = "/var/run/edge-core/edge-core.sock"
var socketFound = false

var transport = http.Transport{
	Dial: func(_, _ string) (net.Conn, error) {
		return net.Dial("unix", unixSocket)
	},
}

type applicationList struct {
	Applications []string `json:"applications"`
}

func initClient() {
	if _, err := os.Stat(unixSocket); os.IsNotExist(err) {
		unixSocket = "/edge-core/edge-core.sock"
		if _, err := os.Stat(unixSocket); os.IsNotExist(err) {
			log.Println("Edge Core API socket does not exist!")
		} else {
			socketFound = true
		}
	} else {
		socketFound = true
	}
	if socketFound {
		restClient.SetTransport(&transport).SetScheme("http").SetHostURL(unixSocket)
	}
}

func restGet(urlToGet string) (*resty.Response, error) {
	resp, err := restClient.R().
		EnableTrace().
		Get(urlToGet)
	return resp, err
}

func restPost(urlToPost string, bodyToPost string) (*resty.Response, error) {
	resp, err := restClient.R().
		SetBody(bodyToPost).
		SetHeader("Content-Type", "application/json").
		Post(urlToPost)
	return resp, err
}

func queryApplications() ([]string, error) {
	getURL := "http://localhost/api/v1/applications"
	var jsonData = new(applicationList)
	resp, err := restClient.R().
		EnableTrace().
		Get(getURL)
	if err != nil {
		log.Println(err)
	} else {
		err = json.Unmarshal([]byte(resp.Body()), &jsonData)
		if err != nil {
			log.Println(err)
		}
	}
	return jsonData.Applications, err
}
