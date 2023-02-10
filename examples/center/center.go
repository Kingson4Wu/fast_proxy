package center

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

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
