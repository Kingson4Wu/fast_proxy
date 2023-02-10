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

	addr := center.GetAddress("token_service")

	if addr == "" {
		panic("service address not exist")
	}

	resp, err := http.Post(fmt.Sprintf("http://%s/api/service", addr),
		"application/x-www-form-urlencoded",
		strings.NewReader(fmt.Sprintf("param=%s", "hello")))
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
