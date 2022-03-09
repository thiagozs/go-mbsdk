package client

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type HttpClientPort interface {
	Get(addrs string) ([]byte, error)
	Post(addrs string, body []byte) ([]byte, error)
	Delete(addrs string, payload []byte) ([]byte, error)

	GetWithResponse(addrs string) (*http.Response, error)
	PostWithResponse(addrs string, payload []byte) (*http.Response, error)
	DeleteWithResponse(addrs string, payload []byte) (*http.Response, error)
	PostFormWithResponse(addrs string) (*http.Response, error)

	GetFreePort() (int, error)

	SetHeader(method, key, value string)
	SetFormValue(method, key, value string)

	DeleteHeader(method, key string)
	DeleteFormValue(method, key string)

	GetHeaders(method string) map[string]string
	GetFormValue(method string) map[string]string

	SetMaxRetry(val int)
	SetMaxRetryWaitMin(val int)
	SetRetryWaitMax(val int)

	DisableLogLevel()
	EnableLogLevel()
}

type HttpClient struct {
	sync.Mutex
	client       *retryablehttp.Client
	MaxRetry     int
	RetryWaitMin int
	RetryWaitMax int
	headers      map[string]map[string]string
	forms        map[string]map[string]string
}

func NewHttpClient(retryWaitMinSec, retryWaitMaxSec, retryMax int) HttpClientPort {
	client := retryablehttp.NewClient()
	client.RetryWaitMin = time.Duration(retryWaitMinSec) * time.Second
	client.RetryWaitMax = time.Duration(retryWaitMaxSec) * time.Second
	client.RetryMax = retryMax
	return &HttpClient{
		client:       client,
		MaxRetry:     retryMax,
		RetryWaitMin: retryWaitMinSec,
		RetryWaitMax: retryWaitMaxSec,
		headers:      make(map[string]map[string]string),
		forms:        make(map[string]map[string]string),
	}
}

func (c *HttpClient) Get(addrs string) ([]byte, error) {
	req, err := retryablehttp.NewRequest(http.MethodGet, addrs, nil)
	if err != nil {
		return []byte{}, err
	}
	if len(c.headers[http.MethodGet]) > 0 {
		for k, v := range c.headers[http.MethodGet] {
			req.Header.Set(k, v)
		}
	}

	rr, err := c.client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer rr.Body.Close()
	return ioutil.ReadAll(rr.Body)
}

func (c *HttpClient) Post(addrs string, payload []byte) ([]byte, error) {

	req, err := retryablehttp.NewRequest(http.MethodPost, addrs, payload)
	if err != nil {
		return []byte{}, err
	}
	if len(c.headers[http.MethodPost]) > 0 {
		for k, v := range c.headers[http.MethodPost] {
			req.Header.Set(k, v)
		}
	}

	rr, err := c.client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	defer rr.Body.Close()
	return ioutil.ReadAll(rr.Body)
}

func (c *HttpClient) Delete(addrs string, payload []byte) ([]byte, error) {

	req, err := retryablehttp.NewRequest(http.MethodDelete, addrs, payload)
	if err != nil {
		return []byte{}, err
	}
	if len(c.headers[http.MethodDelete]) > 0 {
		for k, v := range c.headers[http.MethodDelete] {
			req.Header.Set(k, v)
		}
	}

	rr, err := c.client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	defer rr.Body.Close()
	return ioutil.ReadAll(rr.Body)
}

func (c *HttpClient) SetHeader(method, key, value string) {
	c.Lock()
	defer c.Unlock()
	v, ok := c.headers[strings.ToUpper(method)]
	if ok {
		_, ook := v[key]
		if !ook {
			v[key] = value
		}
	} else {
		c.headers[strings.ToUpper(method)] = make(map[string]string)
		c.headers[strings.ToUpper(method)][key] = value
	}
}

func (c *HttpClient) DeleteHeader(method, key string) {
	c.Lock()
	defer c.Unlock()
	delete(c.headers[strings.ToUpper(method)], key)
}

func (c *HttpClient) SetFormValue(method, key, value string) {
	c.Lock()
	defer c.Unlock()
	v, ok := c.forms[strings.ToUpper(method)]
	if ok {
		_, ook := v[key]
		if !ook {
			v[key] = value
		}
	} else {
		c.forms[strings.ToUpper(method)] = make(map[string]string)
		c.forms[strings.ToUpper(method)][key] = value
	}
}

func (c *HttpClient) DeleteFormValue(method, key string) {
	c.Lock()
	defer c.Unlock()
	delete(c.forms[strings.ToUpper(method)], key)
}

func (c *HttpClient) SetMaxRetry(val int) {
	c.MaxRetry = val
}

func (c *HttpClient) SetMaxRetryWaitMin(val int) {
	c.RetryWaitMin = val
}

func (c *HttpClient) SetRetryWaitMax(val int) {
	c.RetryWaitMax = val
}

func (c *HttpClient) DisableLogLevel() {
	c.client.Logger = nil
}

func (c *HttpClient) EnableLogLevel() {
	c.client.Logger = log.New(os.Stderr, "", log.LstdFlags)
}

func (c *HttpClient) GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func (c *HttpClient) GetWithResponse(addrs string) (*http.Response, error) {
	req, err := retryablehttp.NewRequest(http.MethodGet, addrs, nil)
	if err != nil {
		return &http.Response{}, err
	}
	if len(c.headers[http.MethodGet]) > 0 {
		for k, v := range c.headers[http.MethodGet] {
			req.Header.Set(k, v)
		}
	}

	rr, err := c.client.Do(req)
	if err != nil {
		return &http.Response{}, err
	}
	return rr, nil
}

func (c *HttpClient) PostWithResponse(addrs string, payload []byte) (*http.Response, error) {

	req, err := retryablehttp.NewRequest(http.MethodPost, addrs, payload)
	if err != nil {
		return &http.Response{}, err
	}
	if len(c.headers[http.MethodPost]) > 0 {
		for k, v := range c.headers[http.MethodPost] {
			req.Header.Set(k, v)
		}
	}

	rr, err := c.client.Do(req)
	if err != nil {
		return &http.Response{}, err
	}

	return rr, nil
}

func (c *HttpClient) DeleteWithResponse(addrs string, payload []byte) (*http.Response, error) {

	req, err := retryablehttp.NewRequest(http.MethodDelete, addrs, payload)
	if err != nil {
		return &http.Response{}, err
	}
	if len(c.headers[http.MethodDelete]) > 0 {
		for k, v := range c.headers[http.MethodDelete] {
			req.Header.Set(k, v)
		}
	}

	rr, err := c.client.Do(req)
	if err != nil {
		return &http.Response{}, err
	}

	return rr, nil
}

func (c *HttpClient) PostFormWithResponse(addrs string) (*http.Response, error) {
	forms := url.Values{}
	if len(c.forms[http.MethodPost]) > 0 {
		for k, v := range c.forms[http.MethodPost] {
			forms.Add(k, v)
		}
	}

	req, err := retryablehttp.NewRequest(http.MethodPost, addrs, strings.NewReader(forms.Encode()))
	if err != nil {
		return &http.Response{}, err
	}
	if len(c.headers[http.MethodPost]) > 0 {
		for k, v := range c.headers[http.MethodPost] {
			req.Header.Set(k, v)
		}
	}

	rr, err := c.client.Do(req)
	if err != nil {
		return &http.Response{}, err
	}

	return rr, nil

}

func (c *HttpClient) GetHeaders(method string) map[string]string {
	c.Lock()
	defer c.Unlock()
	return c.headers[strings.ToUpper(method)]
}

func (c *HttpClient) GetFormValue(method string) map[string]string {
	c.Lock()
	defer c.Unlock()
	return c.forms[strings.ToUpper(method)]
}
