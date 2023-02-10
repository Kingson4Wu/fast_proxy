package main

import (
	"fmt"
	"github.com/Kingson4Wu/fast_proxy/common/network"
	"github.com/Kingson4Wu/fast_proxy/examples/center"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {

	serverPort := 8101
	intranetIp := network.GetIntranetIp()
	/**
	curl "http://127.0.0.1:8080/api/service" |jq '.'
	*/

	http.HandleFunc("/api/service", handler)

	/** register service to center async */
	stop := make(chan bool)
	go func(stop chan bool) {
		for {
			select {
			case <-stop:
				log.Println("stop register service")
				return
			default:
				log.Println("try register service")
				if network.Telnet(intranetIp, serverPort) {
					center.Register("token_service", intranetIp, serverPort)
					return
				}
				time.Sleep(3 * time.Second)
			}
		}
	}(stop)

	if err := http.ListenAndServe(":"+strconv.Itoa(serverPort), nil); err != nil {
		close(stop)
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	//time.Sleep(time.Duration(5) * time.Millisecond)
	fmt.Fprintln(w, "{\"code\": 0, \"msg\": \"success\", \"data\": \"{}\"}")
}
