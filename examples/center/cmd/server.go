package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var services map[string]address

type address struct {
	Ip   string
	Port int
}

func (r *address) Get() string {
	if r == nil {
		return ""
	}
	return fmt.Sprintf("%s:%v", r.Ip, r.Port)
}

func init() {
	services = make(map[string]address)
}

func register(name string, ip string, port int) {
	log.Printf("register %s to serivce center\n", name)
	services[name] = address{
		Ip:   ip,
		Port: port,
	}
}

func getAddress(name string) *address {
	if addr, ok := services[name]; ok {
		return &addr

	}
	return nil
}

func main() {
	serverPort := 8080
	http.HandleFunc("/api/register", registerHandler)
	http.HandleFunc("/api/address/get", getAddressHandler)

	fmt.Println("center start ...")
	if err := http.ListenAndServe(":"+strconv.Itoa(serverPort), nil); err != nil {
		log.Fatal(err)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	name := r.Form.Get("name")
	ip := r.Form.Get("ip")
	port, _ := strconv.Atoi(r.Form.Get("port"))
	register(name, ip, port)

	fmt.Fprintln(w, "success")
}

func getAddressHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	name := r.Form.Get("name")
	fmt.Fprintln(w, getAddress(name).Get())
}
