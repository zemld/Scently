package config

import (
	"net/http"
)

var (
	HTTPClient = &http.Client{
		Timeout: HTTPClientTimeout,
		Transport: &http.Transport{
			MaxIdleConns:        HTTPClientMaxIdleConns,
			MaxIdleConnsPerHost: HTTPClientMaxIdleConnsHost,
			IdleConnTimeout:     HTTPClientIdleConnTimeout,
		},
	}
)
