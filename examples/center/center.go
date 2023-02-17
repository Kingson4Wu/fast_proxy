package center

import (
	"fmt"
	"github.com/Kingson4Wu/fast_proxy/common/network"
	"github.com/Kingson4Wu/fast_proxy/common/servicediscovery"
	"io"
	"log"
	"net/http"
	"strconv"
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

func GetSC(f func() string) *servicediscovery.ServiceCenter {
	sc := servicediscovery.Create().
		AddressFunc(func(serviceName string) *servicediscovery.Address {
			addr := GetAddress(serviceName)
			arr := strings.Split(addr, ":")
			if len(arr) == 2 {
				ip := arr[0]
				port, _ := strconv.Atoi(arr[1])
				return &servicediscovery.Address{
					Ip:   ip,
					Port: port,
				}
			}
			return nil
		}).ClientNameFunc(func(req *http.Request) string {
		return req.Header.Get(f())
	}).RegisterFunc(func(name string, ip string, port int) chan bool {
		return RegisterAsync(name, ip, port)
	}).Build()

	return sc
}
