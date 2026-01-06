package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/goharbor/terraform-provider-harbor/models"
)

type Client struct {
	url         string
	username    string
	password    string
	sessionId   string
	bearerToken string
	insecure    bool
	httpClient  *http.Client
	robotPrefix string
}

// NewClient creates common settings
func NewClient(url string, username string, password string, sessionId string, bearerToken string, insecure bool, robotPrefix string) *Client {
	return &Client{
		url:         url,
		username:    username,
		password:    password,
		sessionId:   sessionId,
		bearerToken: bearerToken,
		insecure:    insecure,
		httpClient:  &http.Client{},
		robotPrefix: robotPrefix,
	}
}

// sendRequestWithHeaders sends a http request with specified additional headers
func (c *Client) sendRequestWithHeaders(method string, path string, payload interface{}, extraHeaders http.Header) (resp *http.Response, err error) {
	url := c.url + path
	client := &http.Client{}

	b := new(bytes.Buffer)
	if payload != nil {
		err = json.NewEncoder(b).Encode(payload)
		if err != nil {
			return nil, err
		}
	}

	if c.insecure {
		tr := &http.Transport{
			Proxy:           http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	}

	req, err := http.NewRequest(method, url, b)
	if err != nil {
		return nil, err
	}

	// Use access token authentification if bearer Token is specified
	if c.sessionId != "" {
		req.Header.Add("Cookie", "sid="+c.sessionId)
	} else if c.bearerToken != "" {
		req.Header.Add("Authorization", "Bearer "+c.bearerToken)
	} else {
		req.SetBasicAuth(c.username, c.password)
	}

	req.Header.Add("Content-Type", "application/json")
	for header, values := range extraHeaders {
		for _, value := range values {
			req.Header.Add(header, value)
		}
	}

	resp, err = client.Do(req)
	if err != nil {
		if resp != nil {
			return resp, err
		} else {
			return nil, err
		}
	}

	return resp, nil
}

// SendRequest send a http request
func (c *Client) SendRequest(method string, path string, payload interface{}, statusCode int) (value string, respheaders string, respCode int, err error) {
	resp, err := c.sendRequestWithHeaders(method, path, payload, map[string][]string{})
	if err != nil {
		if resp != nil {
			return "", "", resp.StatusCode, err
		} else {
			return "", "", http.StatusBadGateway, err
		}
	}

	if resp.StatusCode == http.StatusForbidden {
		if c.sessionId != "" {
			csrfheaders, err := extractCsrfHeaders(resp.Header)
			if err != nil {
				return "", "", resp.StatusCode, err
			}
			resp, err = c.sendRequestWithHeaders(method, path, payload, csrfheaders)
			if err != nil {
				if resp != nil {
					return "", "", resp.StatusCode, err
				} else {
					return "", "", http.StatusBadGateway, err
				}
			}
		}
	}

	if resp.StatusCode == http.StatusForbidden {
		return "", "", resp.StatusCode, fmt.Errorf("[ERROR] forbidden: status=%s, code=%d \nIf you are using a robot account, this is likely due to RBAC limitations. See: https://github.com/goharbor/community/blob/main/proposals/new/Robot-Account-Expand.md", resp.Status, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", resp.StatusCode, err
	}
	resp.Body.Close()

	strbody := string(body)

	respHeaders := resp.Header
	headers, err := json.Marshal(respHeaders)
	if err != nil {
		return "", "", resp.StatusCode, err
	}

	if statusCode != 0 && resp.StatusCode != statusCode {
		return "", "", resp.StatusCode, fmt.Errorf("[ERROR] unexpected status code got: %v expected: %v \n %v", resp.StatusCode, statusCode, strbody)
	}

	return strbody, string(headers), resp.StatusCode, nil
}

func extractCsrfHeaders(respHeader http.Header) (headers http.Header, err error) {
	harborcsrfs := respHeader.Values("X-Harbor-Csrf-Token")
	if len(harborcsrfs) > 1 {
		return nil, fmt.Errorf("[ERROR] more than one X-Harbor-Csrf-Token present to retry")
	}
	if len(harborcsrfs) < 1 {
		return nil, fmt.Errorf("[ERROR] no X-Harbor-Csrf-Token present to retry")
	}
	harborcsrf := harborcsrfs[0]

	gorillacsrf := ""
	for _, cookiestr := range respHeader.Values("Set-Cookie") {
		if !strings.Contains(cookiestr, ";") {
			continue
		}
		cookiekv := strings.Split(cookiestr, ";")[0]
		if !strings.Contains(cookiekv, "=") {
			continue
		}
		cookie := strings.Split(cookiekv, "=")
		if cookie[0] == "_gorilla_csrf" {
			gorillacsrf = cookie[1]
			break
		}
	}
	if gorillacsrf == "" {
		return nil, fmt.Errorf("[ERROR] no valid _gorilla_csrf token present to retry")
	}

	headers = map[string][]string{
		"Cookie":              {fmt.Sprintf("_gorilla_csrf=%s", gorillacsrf)},
		"X-Harbor-Csrf-Token": {harborcsrf},
	}

	return headers, nil
}

// GetID gets the resource id from location response header
func GetID(body string) (id string, err error) {
	var jsonData models.ResponseHeaders
	err = json.Unmarshal([]byte(body), &jsonData)
	if err != nil {
		return "", err
	}

	location := jsonData.Location[0]
	location = strings.Replace(location, "/api", "", -1)

	// removes /v2.0 from string if using api version 2
	location = strings.Replace(location, "/v2.0", "", -1)

	return location, nil
}

// get robotPrefix
func (c *Client) GetRobotPrefix() string {
	return c.robotPrefix
}
