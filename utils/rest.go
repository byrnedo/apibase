package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type RestClient struct {
	Headers map[string]string

	client           *http.Client
	lastResponse     http.Response
	lastResponseBody []byte
}

func NewRestClient() *RestClient {
	timeout := time.Duration(5 * time.Second)
	return &RestClient{nil, &http.Client{Timeout: timeout}, http.Response{}, nil}
}

func (c *RestClient) Client() *http.Client {
	return c.Client()
}

func (c *RestClient) LastResponseBody() []byte {
	return c.lastResponseBody
}

func (c *RestClient) LastResponseStatus() int {
	return c.lastResponse.StatusCode
}

func doRequest(c *RestClient, req *http.Request) error {
	for name, value := range c.Headers {
		req.Header.Add(name, value)
	}
	response, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	c.lastResponse = *response
	c.lastResponseBody, err = ioutil.ReadAll(response.Body)
	return err
}

func (c *RestClient) GetWithBytes(url string, byteData []byte) error {
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(byteData))
	if err != nil {
		return err
	}

	return doRequest(c, req)
}

func (c *RestClient) Get(url string) error {
	return c.GetWithBytes(url, nil)
}

func (c *RestClient) Post(url string, data io.Reader) error {

	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		return err
	}

	return doRequest(c, req)
}

func (c *RestClient) PostString(url string, data string) error {
	return c.Post(url, bytes.NewBufferString(data))
}

func (c *RestClient) PostJson(url string, data interface{}) error {

	byteData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.Post(url, bytes.NewBuffer(byteData))
}

func (c *RestClient) Put(url string, data io.Reader) error {

	req, err := http.NewRequest("PUT", url, data)
	if err != nil {
		return err
	}

	return doRequest(c, req)
}

func (c *RestClient) PutString(url string, data string) error {
	return c.Put(url, bytes.NewBufferString(data))
}

func (c *RestClient) PutJson(url string, data interface{}) error {

	byteData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.Put(url, bytes.NewBuffer(byteData))
}
