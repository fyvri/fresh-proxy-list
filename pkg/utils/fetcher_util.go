package utils

import (
	"fmt"
	"io"
	"net/http"
)

type FetcherUtil struct {
	Client         *http.Client
	NewRequestFunc func(method string, url string, body io.Reader) (*http.Request, error)
}

type FetcherUtilInterface interface {
	NewRequest(method, url string, body io.Reader) (*http.Request, error)
	Do(client *http.Client, req *http.Request) (*http.Response, error)
	FetchData(url string) ([]byte, error)
}

func NewFetcher(client *http.Client, newRequestFunc func(method, url string, body io.Reader) (*http.Request, error)) FetcherUtilInterface {
	return &FetcherUtil{
		Client:         client,
		NewRequestFunc: newRequestFunc,
	}
}

func (u *FetcherUtil) Do(client *http.Client, req *http.Request) (*http.Response, error) {
	return client.Do(req)
}

func (u *FetcherUtil) NewRequest(method string, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, url, body)
}

func (u *FetcherUtil) FetchData(url string) ([]byte, error) {
	req, err := u.NewRequestFunc("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := u.Do(u.Client, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return body, fmt.Errorf("failed to fetch data: %s", http.StatusText(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
