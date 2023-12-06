package main

import (
	"fmt"
	"github.com/revproxy/src/model/hostnames"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/revproxy/src/gorm-library/database"
	"gorm.io/gorm"
)

var (
	Db            *gorm.DB
	dbOnce        sync.Once
	hostname_maps = map[string]string{
		"customer1.com": "subnewdomain01-staticawsdiscoverytest",
		"customer2.com": "globalnewdomain01-staticawsdiscoverytest",
	}
)

func NewDB() *gorm.DB {
	dbOnce.Do(func() {
		Db = database.UseDB()
	})
	return Db
}

func startBackground() {
	NewDB()
}

func main() {
	startBackground()
	config, err := database.ReadConfig()
	if err != nil {
		panic(fmt.Sprintf("cant read the property file.", err))
	}

	for k, v := range hostname_maps {

		hostname := hostnames.HostName{
			Domain: k,
			Bucket: v,
		}
		hostnames.CreateHostName(Db, &hostname)
	}

	//fmt.Println(hostnames.GetAllHostNames(Db))
	proxy, err := NewProxy(upstream)
	if err != nil {
		fmt.Errorf("error received %w", err)
	}
	r := mux.NewRouter()
	//Read all the hostnames in DB, make a sub router for each hostname.
	//Gorilla mux matches routes based on the order of registration. Not the longest path match.
	//the wildcards with PathPrefix are suggested to be registered last.
	r.HandleFunc("/404", notFoundHandler)
	r.PathPrefix("/").Handler(proxy)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port), r))
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("page does not exist"))
}
