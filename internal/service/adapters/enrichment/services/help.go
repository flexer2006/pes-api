package api

import (
	"errors"
	"net/http"
)

var (
	ErrEmptyName      = errors.New("name cannot be empty")
	ErrNon200Response = errors.New("API returned non-200 status code")
)

type APIClient struct {
	baseURL    string
	httpClient HTTPClient
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
