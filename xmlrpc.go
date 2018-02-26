package xmlrpc

import (
	"bytes"
	"net/http"
)

// Client ...
type Client interface {
	Call(method string, args ...*Value) (*Value, error)
}

// Call XMLRPC
func Call(url, method string, args ...*Value) (*Value, error) {
	return call(httpClient, url, method, args)
}

var httpClient = &http.Client{}

type client struct {
	url string
}

// NewClient ...
func NewClient(url string) Client {
	return &client{
		url: url,
	}
}

func (c *client) Call(method string, args ...*Value) (*Value, error) {
	return call(httpClient, c.url, method, args)
}

func call(client *http.Client, url, method string, args []*Value) (*Value, error) {
	req, err := NewRequest(method, args...)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = req.Write(buf)
	if err != nil {
		return nil, err
	}

	resp, err := client.Post(url, "text/xml", buf)
	if err != nil {
		return nil, err
	}

	return ParseResponse(resp.Body)
}
