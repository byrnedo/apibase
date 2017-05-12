package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
	"strings"
)

type Response struct {
	Response     http.Response
}

func (this *Response) GetBody() []byte{
	defer this.Response.Body.Close()
	r, _ := ioutil.ReadAll(this.Response.Body)
	return r
}

// Must be called if not using AsJson or GetBody
func (this *Response) Close() {
	this.Response.Body.Close()
}

func (this *Response) AsJson(result interface{}) error {
	return json.NewDecoder(bytes.NewReader(this.GetBody())).Decode(result)
}

type RestClient struct {
	Headers map[string]string
	client           *http.Client
}

type ClientOptsFunc func(*http.Client)

func NewRestClient(optFuncs ...ClientOptsFunc) *RestClient {
	timeout := time.Duration(5 * time.Second)
	c := &http.Client{Timeout: timeout}
	for _, f := range optFuncs {
		f(c)
	}
	return &RestClient{nil, c}
}

func (c *RestClient) Client() *http.Client {
	return c.Client()
}

func doRequest(c *RestClient, req *http.Request) (*Response, error) {
	for name, value := range c.Headers {
		req.Header.Add(name, value)
	}
	response, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	res := &Response{
		Response: *response,
	}
	return res, err
}

func (c *RestClient) GetWithBytes(url string, byteData []byte) (*Response,  error) {
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(byteData))
	if err != nil {
		return nil, err
	}

	return doRequest(c, req)
}

func (c *RestClient) Get(url string) (*Response, error) {
	return c.GetWithBytes(url, nil)
}

func (c *RestClient) Post(url string, data io.Reader) (*Response, error) {

	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		return nil, err
	}

	return doRequest(c, req)
}

func (c *RestClient) PostString(url string, data string) (*Response, error) {
	return c.Post(url, bytes.NewBufferString(data))
}

func (c *RestClient) PostJson(url string, data interface{}) (*Response, error) {

	byteData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return c.Post(url, bytes.NewBuffer(byteData))
}

func (c *RestClient) Put(url string, data io.Reader) (*Response, error) {

	req, err := http.NewRequest("PUT", url, data)
	if err != nil {
		return nil, err
	}

	return doRequest(c, req)
}

func (c *RestClient) PutString(url string, data string) (*Response, error) {
	return c.Put(url, bytes.NewBufferString(data))
}

func (c *RestClient) PutJson(url string, data interface{}) (*Response, error) {

	byteData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return c.Put(url, bytes.NewBuffer(byteData))
}

func (c *RestClient) Patch(url string, data io.Reader) (*Response, error) {

	req, err := http.NewRequest("PATCH", url, data)
	if err != nil {
		return nil, err
	}

	return doRequest(c, req)
}

func (c *RestClient) PatchString(url string, data string) (*Response, error) {
	return c.Patch(url, bytes.NewBufferString(data))
}

func (c *RestClient) PatchJson(url string, data interface{}) (*Response, error) {

	byteData, err := json.Marshal(data)
	if err != nil {
		return nil, err

	}
	return c.Patch(url, bytes.NewBuffer(byteData))
}

func (c *RestClient) Delete(url string, data io.Reader) (*Response, error) {
	req, err := http.NewRequest("DELETE", url, data)
	if err != nil {
		return nil, err
	}

	return doRequest(c, req)
}

func (c *RestClient) DeleteString(url string, data string) (*Response, error) {
	return c.Put(url, bytes.NewBufferString(data))
}

func (c *RestClient) DeleteJson(url string, data interface{}) (*Response, error) {

	byteData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return c.Put(url, bytes.NewBuffer(byteData))
}

func  (c *RestClient) Do(method string, url string, data io.Reader) (*Response, error) {
	switch strings.ToLower(method) {
	case "get":
		return c.Get(url)
	case "post":
		return c.Post(url, data)
	case "put":
		return c.Put(url, data)
	case "patch":
		return c.Patch(url, data)
	case "delete":
		return c.Delete(url, data)
	default:
		panic("invalid method " + method)
	}
}

func  (c *RestClient) DoString(method string, url string, data string) (*Response, error) {
	switch strings.ToLower(method) {
	case "get":
		return c.Get(url)
	case "post":
		return c.PostString(url, data)
	case "put":
		return c.PutString(url, data)
	case "patch":
		return c.PatchString(url, data)
	case "delete":
		return c.DeleteString(url, data)
	default:
		panic("invalid method " + method)
	}
}

func  (c *RestClient) DoJson(method string, url string, data interface{}) (*Response, error) {
	switch strings.ToLower(method) {
	case "get":
		return c.Get(url)
	case "post":
		return c.PostJson(url, data)
	case "put":
		return c.PutJson(url, data)
	case "patch":
		return c.PatchJson(url, data)
	case "delete":
		return c.DeleteJson(url, data)
	default:
		panic("invalid method " + method)
	}
}
