package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/BESTSELLER/terraform-provider-harbor/models"
)

type Client struct {
	url        string
	username   string
	password   string
	insecure   bool
	httpClient *http.Client
}

// NewClient creates common settings
func NewClient(url string, username string, password string, insecure bool) *Client {

	return &Client{
		url:        url,
		username:   username,
		password:   password,
		insecure:   insecure,
		httpClient: &http.Client{},
	}
}

// SendRequest send a http request
func (c *Client) SendRequest(method string, path string, payload interface{}, statusCode int) (value string, respheaders string, err error) {
	url := c.url + path
	client := &http.Client{}

	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(payload)
	if err != nil {
		return "", "", err
	}

	if c.insecure == true {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	}

	req, err := http.NewRequest(method, url, b)
	if err != nil {
		return "", "", err
	}

	req.SetBasicAuth(c.username, c.password)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	resp.Body.Close()

	strbody := string(body)

	respHeaders := resp.Header
	headers, err := json.Marshal(respHeaders)
	if err != nil {
		return "", "", err
	}

	if statusCode != 0 {
		if resp.StatusCode != statusCode {

			return "", "", fmt.Errorf("[ERROR] unexpected status code got: %v expected: %v \n %v", resp.StatusCode, statusCode, strbody)
		}
	}

	return strbody, string(headers), nil
}

// GetID gets the resource id from locaation resposne header
func GetID(body string) (id string, err error) {
	var jsonData models.ResponseHeaders
	err = json.Unmarshal([]byte(body), &jsonData)
	if err != nil {
		return "", err
	}

	localation := jsonData.Location[0]
	localation = strings.Replace(localation, "/api", "", -1)

	// removes /v2.0 from string if using api version 2
	localation = strings.Replace(localation, "/v2.0", "", -1)

	return localation, nil
}
