package stableibkr

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Auth_Status struct {
	Authenticated bool     `json:"authenticated"`
	Connected     bool     `json:"connected"`
	Competing     bool     `json:"Competing"`
	Fail          string   `json:"fail"`
	Message       string   `json:"message"`
	Prompts       []string `json:"prompts"`
}

// This verifies that the GW has been logged in
// if it has not, it will block
func Client() http.Client {
	// ignore Cert errors
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// create the client
	client := http.Client{
		Transport: tr,
		Jar:       nil,
		Timeout:   0,
	}

	for !AuthStatus(client) {
		fmt.Println("Gateway is not Auth'ed. Please log in and press Enter")
		fmt.Scanln() // wait for user input
	}
	fmt.Println("Gateway is authenticated, have fun!")
	return client
}

func RestGet(client http.Client, apiCall string) http.Response {
	req, _ := http.NewRequest("GET", "https://192.168.1.2:5000"+apiCall, nil)
	req.Header.Add("accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	return *resp
}

func AuthStatus(client http.Client) bool {
	//API
	apiCall := "/v1/api/iserver/auth/status"
	loggedIn := false

	resp := RestGet(client, apiCall)

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	var auth_status Auth_Status
	err = json.Unmarshal(data, &auth_status)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Logged in: ", auth_status.Authenticated)
	fmt.Println("Connected: ", auth_status.Connected)
	if auth_status.Authenticated {
		loggedIn = true
	}

	return loggedIn
}

func Tickle(client http.Client) {
	apiCall := "/v1/api/tickle"

	resp := RestGet(client, apiCall)
	if resp.StatusCode == 200 {
		fmt.Println("Gateway login refreshed")
	}
}
