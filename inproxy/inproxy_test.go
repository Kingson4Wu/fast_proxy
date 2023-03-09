package inproxy

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/Kingson4Wu/fast_proxy/common/server"
	"github.com/Kingson4Wu/fast_proxy/inproxy/inconfig"
	"github.com/Kingson4Wu/fast_proxy/test"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func init() {

	configBytes, err := ConfigFs.ReadFile("testdata/config.yaml")
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	c := inconfig.LoadYamlConfig(configBytes)

	go NewServer(c, server.WithServiceCenter(test.GetSC()))
	go test.Service()

	time.Sleep(3 * time.Second)

}

//go test -v -run=^$ -bench=.

func TestRequestProxy(t *testing.T) {

	httpPostReq :=
		"POST /inProxy/search_service/api/service HTTP/1.0\r\n" +
			"Content-Type:application/x-www-form-urlencoded\r\n" +
			"C_ServiceName:chat_service\r\n\r\n"

	r, err := http.ReadRequest(bufio.NewReader(strings.NewReader(httpPostReq)))

	data, err := os.ReadFile("../outproxy/testdata/encodeReq.golden")
	if err != nil {
		t.Fatal(err)
	}
	r.Body = io.NopCloser(bytes.NewReader(data))

	rw := httptest.NewRecorder()
	requestProxy(rw, r)

	if rw.Code != http.StatusOK {

		fmt.Println(rw.Header().Get("proxy_error_message"))

		t.Fatal("resp code error")
	}

}

func BenchmarkRequestProxy(b *testing.B) {
	b.ReportAllocs()
	httpPostReq :=
		"POST /inProxy/search_service/api/service HTTP/1.0\r\n" +
			"Content-Type:application/x-www-form-urlencoded\r\n" +
			"C_ServiceName:chat_service\r\n\r\n"

	r, err := http.ReadRequest(bufio.NewReader(strings.NewReader(httpPostReq)))

	data, err := os.ReadFile("../outproxy/testdata/encodeReq.golden")
	if err != nil {
		b.Fatal(err)
	}
	r.Body = io.NopCloser(bytes.NewReader(data))

	rw := httptest.NewRecorder()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		requestProxy(rw, r)
	}
}

//go test -v -run=^$ -bench=^BenchmarkRequestProxy$ -benchtime=2s -cpuprofile=cpu.prof
//go tool pprof inproxy.test cpu.prof
//go test -v -run=^$ -bench=^BenchmarkRequestProxy$ -benchtime=2s -memprofile=mem.prof
//go tool pprof -sample_index=alloc_space  inproxy.test mem.prof

func BenchmarkRequestProxyParallel(b *testing.B) {

	httpPostReq :=
		"POST /inProxy/search_service/api/service HTTP/1.0\r\n" +
			"Content-Type:application/x-www-form-urlencoded\r\n" +
			"C_ServiceName:chat_service\r\n\r\n"

	r, err := http.ReadRequest(bufio.NewReader(strings.NewReader(httpPostReq)))

	data, err := os.ReadFile("../outproxy/testdata/encodeReq.golden")
	if err != nil {
		b.Fatal(err)
	}
	r.Body = io.NopCloser(bytes.NewReader(data))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		rw := httptest.NewRecorder()
		for pb.Next() {
			requestProxy(rw, r)
		}
	})
}

//go test -bench=Parallel -blockprofile=block.prof
//go tool pprof inproxy.test block.prof

func TestEncodeRespToGolden(t *testing.T) {

	httpPostReq :=
		"POST /inProxy/search_service/api/service HTTP/1.0\r\n" +
			"Content-Type:application/x-www-form-urlencoded\r\n" +
			"C_ServiceName:chat_service\r\n\r\n"

	r, err := http.ReadRequest(bufio.NewReader(strings.NewReader(httpPostReq)))

	data, err := os.ReadFile("../outproxy/testdata/encodeReq.golden")
	if err != nil {
		t.Fatal(err)
	}
	r.Body = io.NopCloser(bytes.NewReader(data))

	rw := httptest.NewRecorder()
	requestProxy(rw, r)

	if rw.Code != http.StatusOK {

		fmt.Println(rw.Header().Get("proxy_error_message"))

		t.Fatal("resp code error")
	}

	dir := "testdata"
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile("testdata/encodeResp.golden", rw.Body.Bytes(), 0644)
	if err != nil {
		t.Fatal(err)
	}
}
