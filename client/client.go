package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

func (c *Client) SendRequest(method string, path string, payload interface{}, statusCode int) (value string, err error) {
	url := c.url + path
	client := &http.Client{}

	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(payload)
	if err != nil {
		return "", err
	}

	if c.insecure == true {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	}

	req, err := http.NewRequest(method, url, b)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(c.username, c.password)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	resp.Body.Close()

	strbody := string(body)

	if statusCode != 0 {
		if resp.StatusCode != statusCode {

			return "", fmt.Errorf("[ERROR] unexpected status code got: %v expected: %v \n %v", resp.StatusCode, statusCode, strbody)
		}
	}

	return strbody, nil
}
