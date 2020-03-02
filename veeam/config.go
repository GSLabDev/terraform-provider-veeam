package veeam

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// Config ... configuration for VEEAM
type Config struct {
	ServerIP string
	Port     int
	Username string
	Password string
	Scheme   string
}

// GetResponse ... get Response according to request
func (c *Config) GetResponse(request *http.Request) ([]byte, error) {

	token, err := GetToken(c.ServerIP, c.Port, c.Username, c.Password, c.Scheme)
	if err != nil {
		log.Println("[ERROR] Error in getting token")
		return nil, fmt.Errorf("[Error] .\n Error: %s", err.Error())
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	var tempURL *url.URL
	tempURL, err = url.Parse(c.Scheme + "://" + c.ServerIP + ":" + strconv.Itoa(c.Port) + "/api/" + request.URL.String())
	if err != nil {
		log.Println("[Error] URL is not in correct format")
		return nil, fmt.Errorf("[Error]  Error: %s", err.Error())
	}
	request.URL = tempURL

	log.Println(request.URL)
	log.Println(token)
	fmt.Println(request.Method)
	request.Header.Set("X-RestSvcSessionId", token)
	if request.Method == "POST" {
		request.Header.Set("Accept", "application/xml")
		request.Header.Set("Content-Type", "application/xml")

	} else {
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Content-Type", "application/json")

	}

	client := &http.Client{Transport: tr}
	resp, err := client.Do(request)
	if err != nil {
		log.Println(" [ERROR] Do: ", err)
		return nil, fmt.Errorf("[Error]  Error: %s", err.Error())

	}

	if resp.StatusCode >= 200 && resp.StatusCode <= 204 {
		log.Println("success..")

	} else {
		return nil, fmt.Errorf("[Error]  Error: %s", err.Error())
	}

	return ioutil.ReadAll(resp.Body)
}
