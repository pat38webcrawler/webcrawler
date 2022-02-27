package main

import (
	"io"
	"net/http"
)

// struct  representing the http client
type wcClient struct {
	*http.Client
	baseURL string
}

// use to send a GET request to the webcrawler server in order to initialize a sitemap computing for
// the url pass in parameter

func (c *wcClient) NewRequest(method, path string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, c.baseURL+path, body)
}

// function that return a http client to ocnnect to webrawler server (port 8900)
func (c *wcClient) getClient(host string, port string) (baseURL string, client *http.Client, err error) {
	client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	baseURL = "http://" + host + ":" + port
	return baseURL, client, nil
}

// function to get a client
func client(host, port string) (*wcClient, error) {
	var err error
	wcClient := wcClient{}

	wcClient.baseURL, wcClient.Client, err = wcClient.getClient(host, port)
	if err != nil {
		return nil, err
	}
	return &wcClient, nil
}
