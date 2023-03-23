package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/margostino/gover/common"
	"net/http"
)

type HttpClient struct {
	client *http.Client
}

func New() *HttpClient {
	return &HttpClient{
		client: &http.Client{},
	}
}

func (c *HttpClient) Post(url string, data interface{}, accessToken string, response interface{}) {
	request := buildRequest(http.MethodPost, url, data, accessToken)
	c.call(request, response)
}

func (c *HttpClient) Patch(url string, data interface{}, accessToken string, response interface{}) {
	request := buildRequest(http.MethodPatch, url, data, accessToken)
	c.call(request, response)
}

func buildRequest(method string, url string, data interface{}, accessToken string) *http.Request {
	body, err := json.Marshal(data)
	common.Check(err)
	request, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	common.Check(err)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	request.Header.Set("Content-Type", "application/json")
	return request
}

func (c *HttpClient) call(request *http.Request, response interface{}) {
	res, err := c.client.Do(request)
	common.Check(err)
	err = json.NewDecoder(res.Body).Decode(&response)
	common.Check(err)

	if res.StatusCode != 200 && res.StatusCode != 201 {
		common.Fatal(fmt.Sprintf("No successful status code from %s: %d", request.URL, res.StatusCode))
	}
}
