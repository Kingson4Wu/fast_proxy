package test

import (
	"fmt"
	"github.com/Kingson4Wu/fast_proxy/common/servicediscovery"
	"log"
	"net/http"
)

func GetSC() *servicediscovery.ServiceCenter {
	sc := servicediscovery.Create().
		AddressFunc(func(serviceName string) *servicediscovery.Address {
			return &servicediscovery.Address{
				Ip:   "127.0.0.1",
				Port: 9988,
			}
		}).ClientNameFunc(func(req *http.Request) string {
		return req.Header.Get("C_ServiceName")
	}).RegisterFunc(func(name string, ip string, port int) chan bool {
		return make(chan bool)
	}).Build()

	return sc
}

func Service() {

	mux := http.NewServeMux()
	mux.HandleFunc("/api/service", searchHandler)

	server := &http.Server{
		Addr:    ":9988",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "{\"code\": 1, \"msg\": \"search success\", \"data\": \"{}\"}")
}
