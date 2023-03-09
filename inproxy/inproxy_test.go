package inproxy

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/Kingson4Wu/fast_proxy/common/server"
	"github.com/Kingson4Wu/fast_proxy/examples/center"
	"github.com/Kingson4Wu/fast_proxy/inproxy/inconfig"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func init() {

	configBytes, err := ConfigFs.ReadFile("internal/config.yaml")
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	c := inconfig.LoadYamlConfig(configBytes)

	sc := center.GetSC(func() string { return inconfig.Get().ServiceRpcHeaderName() })

	go NewServer(c, server.WithServiceCenter(sc))

	//TODO ，mock，todo

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
