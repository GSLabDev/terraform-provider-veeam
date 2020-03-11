package veeam

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Config ... configuration for VEEAM
type Config struct {
	ServerIP string
	Port     int
	Username string
	Password string
}

// GetResponse ... get Response according to request
func (c *Config) GetResponse(request *http.Request) ([]byte, error) {

	token, err := GetToken(c.ServerIP, c.Port, c.Username, c.Password)
	if err != nil {
		log.Println("[ERROR] Error in getting token")
		return nil, fmt.Errorf("[Error] .\n Error: %s", err.Error())
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	var tempURL *url.URL
	tempURL, err = url.Parse("http://" + c.ServerIP + ":" + strconv.Itoa(c.Port) + "/api/" + request.URL.String())
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
		data, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("[Error]  Error %d: %s", resp.StatusCode, strings.ToValidUTF8(string(data), ""))
	}

	return ioutil.ReadAll(resp.Body)
}
