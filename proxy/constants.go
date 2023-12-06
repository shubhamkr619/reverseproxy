package main

import "time"

const (
	upstream     = "https://storage.googleapis.com/"
	upstreamHost = "storage.googleapis.com"
	originHost   = "origin_host"
)

var (
	customerRequestHeaderPerCustomer = map[string]map[string]string{
		"content.hsts.skuttle.io": {
			"origin_host": "content.hsts.skuttle.io",
		},
	}

	customerResponseHeaderPerCustomer = map[string]map[string]string{
		"content.hsts.skuttle.io": {
			"custom": "customer1",
		},
	}

	defaultResponseHeadersMap = map[string]string{
		"defaultExpire": "300",
	}

	defaultRequestHeadersMap = map[string]string{
		"startingAt": time.Now().String(),
	}
)
