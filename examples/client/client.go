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

	send("song_service", "token_service")
	send("chat_service", "search_service")
}

func send(clientName, destServiceName string) {
	req, err := http.NewRequest("POST", getAddress(destServiceName), strings.NewReader(fmt.Sprintf("param=%s", "hello")))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("C_ServiceName", clientName)

	client := &http.Client{}
	resp, err := client.Do(req)

	/*resp, err := http.Post(getAddress(),
	"application/x-www-form-urlencoded",
	strings.NewReader(fmt.Sprintf("param=%s", "hello")))*/
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Printf("resp code :%v\n", resp.StatusCode)

	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	bodyBytes, _ := io.ReadAll(resp.Body)

	log.Println(string(bodyBytes))
}

func getAddress(destServiceName string) string {

	proxyEnable := true

	if proxyEnable {
		addr := center.GetAddress("out_proxy")
		if addr == "" {
			panic("service address not exist")
		}
		return fmt.Sprintf("http://%s/%s/api/service", addr, destServiceName)
	}

	addr := center.GetAddress(destServiceName)

	if addr == "" {
		panic("service address not exist")
	}
	return fmt.Sprintf("http://%s/api/service", addr)
}
