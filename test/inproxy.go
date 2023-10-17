package test

import (
	"log"
	"net/http"
	"os"
)

func MockInProxyServe() {

	data, err := os.ReadFile("../inproxy/testdata/encodeResp.golden")
	if err != nil {
		log.Fatal(err)
	}

	handler := func(w http.ResponseWriter, r *http.Request) {

		_, _ = w.Write(data)
	}

	http.HandleFunc("/", handler)
	err = http.ListenAndServe(":9033", nil)
	if err != nil {
		log.Fatal(err)
	}
}
