package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type ReverseProxy struct {
	p *httputil.ReverseProxy
}

func (s *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.p.ServeHTTP(w, r)
}

func NewProxy(upstream string) (*ReverseProxy, error) {
	upstreamUrl, err := url.Parse(upstream)
	if err != nil {
		panic(fmt.Errorf("error while working with uri %w", err))
	}
	fmt.Println(upstreamUrl)
	proxy := &ReverseProxy{httputil.NewSingleHostReverseProxy(upstreamUrl)}
	oldDirector := proxy.p.Director
	proxy.p.Director = func(request *http.Request) {
		oldDirector(request)
		modifyRequest(request)
	}
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
	}
	//Modify transport as per need
	proxy.p.Transport = transport
	proxy.p.ModifyResponse = func(response *http.Response) error {
		return modifyResponse(response)
	}
	proxy.p.ErrorHandler = ErrHandle
	return proxy, nil
}

func ErrHandle(res http.ResponseWriter, req *http.Request, err error) {
	fmt.Println(err)
}

func modifyResponse(r *http.Response) error {
	oldHost := r.Request.Header.Get(originHost)
	responseHeaders, found := customerResponseHeaderPerCustomer[oldHost]
	setDefaultResponseOverrideHeaders(r, defaultResponseHeadersMap)
	if !found {
		log.Printf("key does not exist in the map: " + oldHost)
	} else {
		setCustomerResponseHeaders(r, responseHeaders)
	}
	return nil
}

func modifyRequest(request *http.Request) {
	bucket, err := getHostNameMap(request)
	if err != nil {
		request.URL.Path = "/404"
		request.Host = "localhost:8080"
		request.URL.Scheme = "http"
		request.Header.Set("Host", "localhost")
	} else {
		request.Header.Set("Host", upstreamHost)
		setDefaultRequestHeaders(request, defaultRequestHeadersMap)
		reqHeaders, found := customerRequestHeaderPerCustomer[request.Host]
		if !found {
			log.Printf("modify requeset key does not exist in the map: " + request.Host)
		} else {
			setCustomerRequestHeaders(request, reqHeaders)
		}
		originHeader := request.Host
		nPath, _ := url.JoinPath("/", bucket, request.URL.Path)
		request.URL.Path = nPath
		request.Host = upstreamHost
		request.URL.Scheme = "https"
		request.Header.Set(originHost, originHeader)
		fmt.Printf("Upstream Host: %s\n", request.Host)
		fmt.Printf("Path : %s\n", nPath)
		fmt.Printf("Bucket : %s\n", bucket)
	}
}
