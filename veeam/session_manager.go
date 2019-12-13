package veeam

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// GetToken ... Get Token from existing store or create a new one.
func GetToken(serverIP string, port int, username, password string) (string, error) {

	token, err := getTokenFromHost(serverIP, port, username, password)
	if err != nil {
		log.Printf("[Error] Cannot get Token : %s", err.Error())
		return "", err
	}
	return token, nil
}

func getTokenFromHost(serverIP string, port int, username, password string) (string, error) {
	req, err := http.NewRequest("POST", "http://"+serverIP+":"+strconv.Itoa(port)+"/api/sessionMngr/?v=v1_4", nil)
	if err != nil {
		log.Println("[ERROR] Error while requesting sessionID ", err)
		return "", err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	req.SetBasicAuth(username, password)
	req.Header.Add("Accept", "application/json")

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {

		return "", fmt.Errorf("[Error] While connecting to server .. Error in connection. Check Username, Password and Server IP")

	}
	if resp.StatusCode == 401 {
		log.Println("[ERROR] Error in connection. Check Username, Password and serverIP")
		return "", fmt.Errorf("[ERROR] Error in connection. Check Username, Password and Server IP")
	}

	val := resp.Header.Get("X-Restsvcsessionid")

	return val, nil
}
