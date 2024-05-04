package google

import (
	"net/http"
)

type GoogleAdapter struct {
	httpClient *http.Client
}

func NewGoogle() *GoogleAdapter {
	return &GoogleAdapter{
		httpClient: &http.Client{},
	}
}
