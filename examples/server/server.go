package main

import (
	"fmt"
	"github.com/Kingson4Wu/fast_proxy/common/network"
	"github.com/Kingson4Wu/fast_proxy/examples/center"
	"log"
	"net/http"
	"strconv"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/api/service", handler)
	go service(8101, "token_service", mux)

	mux2 := http.NewServeMux()
	mux2.HandleFunc("/api/service", searchHandler)
	go service(8102, "search_service", mux2)

	// prevent program from exiting
	select {}

}

func service(serverPort int, serviceName string, mux http.Handler) {

	server := &http.Server{
		Addr:    ":" + strconv.Itoa(serverPort),
		Handler: mux,
	}

	intranetIp := network.GetIntranetIp()
	/**
	curl "http://127.0.0.1:8080/api/service" |jq '.'
	*/

	stop := center.RegisterAsync(serviceName, intranetIp, serverPort)

	if err := server.ListenAndServe(); err != nil {
		close(stop)
		log.Fatal(err)
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "{\"code\": 1, \"msg\": \"search success\", \"data\": \"{}\"}")
}

func handler(w http.ResponseWriter, r *http.Request) {
	//time.Sleep(time.Duration(5) * time.Millisecond)
	fmt.Fprintln(w, "{\"code\": 0, \"msg\": \"success\", \"data\": \"{}\"}")
}
