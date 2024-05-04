package facebook

import (
	"net/http"
)

type FacebookAdapter struct {
	httpClient *http.Client
}

func NewFacebook() *FacebookAdapter {
	return &FacebookAdapter{
		httpClient: &http.Client{},
	}
}
