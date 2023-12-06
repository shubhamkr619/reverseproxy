package main

import (
	"errors"
	"net/http"

	"github.com/revproxy/src/model/hostnames"
)

func getHostNameMap(request *http.Request) (string, error) {

	hostName := request.Host
	bucket, err := getBucketPerHostname(hostName)
	if err != nil {
		return "", err
	}
	return bucket, nil
}

func createHostName(hostname hostnames.HostName) {
	hostnames.CreateHostName(Db, &hostname)
}

func getBucketPerHostname(domain string) (string, error) {
	bucket := hostnames.GetBucketByHostName(Db, domain)
	//hostnames.GetAllHostNames(Db)
	//if err != nil {
	//	// log.Errorf("Error fetching the Hostnames from SQL DB")
	//	return "", err
	//}
	//for _, hostname := range hostnames_list {
	//	if hostname.Domain == domain {
	//		return hostname.Bucket, nil
	//	}
	//}
	if bucket != "" {
		return bucket, nil
	}
	return "", errors.New("failed to find a matching domain")
}

func setCustomerRequestHeaders(r *http.Request, headers map[string]string) {
	for k, v := range headers {
		r.Header.Set(k, v)
	}
}

func setDefaultRequestHeaders(r *http.Request, headers map[string]string) {
	setCustomerRequestHeaders(r, headers)
}

func setCustomerResponseHeaders(r *http.Response, headers map[string]string) {
	for k, v := range headers {
		r.Header.Set(k, v)
	}
}

func setDefaultResponseOverrideHeaders(r *http.Response, headersMap map[string]string) {
	setCustomerResponseHeaders(r, headersMap)
}
