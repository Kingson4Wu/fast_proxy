package center

import (
	"fmt"
	"github.com/Kingson4Wu/fast_proxy/common/network"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func RegisterAsync(name string, ip string, port int) chan bool {
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
				if network.Telnet(ip, port) {
					Register(name, ip, port)
					return
				}
				time.Sleep(3 * time.Second)
			}
		}
	}(stop)
	return stop
}

func Register(name string, ip string, port int) {
	log.Printf("register %s to serivce center\n", name)

	resp, err := http.Post("http://127.0.0.1:8080/api/register",
		"application/x-www-form-urlencoded",
		strings.NewReader(fmt.Sprintf("name=%s&ip=%s&port=%v", name, ip, port)))
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	bodyBytes, _ := io.ReadAll(resp.Body)

	log.Println(string(bodyBytes))

}

func GetAddress(name string) string {
	resp, err := http.Post("http://127.0.0.1:8080/api/address/get",
		"application/x-www-form-urlencoded",
		strings.NewReader(fmt.Sprintf("name=%s", name)))
	if err != nil {
		fmt.Println(err)
		return ""
	}

	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	bodyBytes, _ := io.ReadAll(resp.Body)

	return strings.TrimSpace(string(bodyBytes))
}
