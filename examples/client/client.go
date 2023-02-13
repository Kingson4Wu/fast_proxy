package main

import (
	"fmt"
	"github.com/Kingson4Wu/fast_proxy/examples/center"
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {

	req, err := http.NewRequest("POST", getAddress(), strings.NewReader(fmt.Sprintf("param=%s", "hello")))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("C_ServiceName", "song_service")

	client := &http.Client{}
	resp, err := client.Do(req)

	/*resp, err := http.Post(getAddress(),
	"application/x-www-form-urlencoded",
	strings.NewReader(fmt.Sprintf("param=%s", "hello")))*/
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

func getAddress() string {

	proxyEnable := true

	if proxyEnable {
		addr := center.GetAddress("out_proxy")
		if addr == "" {
			panic("service address not exist")
		}
		return fmt.Sprintf("http://%s/token_service/api/service", addr)
	}

	addr := center.GetAddress("token_service")

	if addr == "" {
		panic("service address not exist")
	}
	return fmt.Sprintf("http://%s/api/service", addr)
}
