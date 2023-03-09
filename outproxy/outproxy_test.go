package outproxy

import (
	"bufio"
	"fmt"
	"github.com/Kingson4Wu/fast_proxy/common/server"
	"github.com/Kingson4Wu/fast_proxy/examples/center"
	"github.com/Kingson4Wu/fast_proxy/outproxy/internal/pack"
	"github.com/Kingson4Wu/fast_proxy/outproxy/outconfig"
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

	c := outconfig.LoadYamlConfig(configBytes)

	sc := center.GetSC(func() string { return outconfig.Get().ServiceRpcHeaderName() })

	go NewServer(c, server.WithServiceCenter(sc))

	//TODO ，mock，todo

	time.Sleep(3 * time.Second)

}

//go test -v -run=^$ -bench=.

func TestRequestProxy(t *testing.T) {
	httpPostReq :=
		"POST /search_service/api/service HTTP/1.0\r\n" +
			"Content-Type:application/x-www-form-urlencoded\r\n" +
			"C_ServiceName:chat_service\r\n\r\n" +
			"param=hello\r\n"

	r, err := http.ReadRequest(bufio.NewReader(strings.NewReader(httpPostReq)))
	if err != nil {
		t.Fatal(err)
	}
	rw := httptest.NewRecorder()
	requestProxy(rw, r)
	fmt.Println(string(rw.Body.Bytes()))

	if rw.Code != http.StatusOK {
		t.Fatal("resp code error")
	}

}

func BenchmarkRequestProxy(b *testing.B) {
	b.ReportAllocs()
	httpPostReq :=
		"POST /search_service/api/service HTTP/1.0\r\n" +
			"Content-Type:application/x-www-form-urlencoded\r\n" +
			"C_ServiceName:chat_service\r\n\r\n" +
			"param=hello\r\n"

	r, err := http.ReadRequest(bufio.NewReader(strings.NewReader(httpPostReq)))
	if err != nil {
		b.Fatal(err)
	}
	rw := httptest.NewRecorder()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		requestProxy(rw, r)
	}
}

//go test -v -run=^$ -bench=^BenchmarkRequestProxy$ -benchtime=2s -cpuprofile=cpu.prof
//go tool pprof outproxy.test cpu.prof
//go test -v -run=^$ -bench=^BenchmarkRequestProxy$ -benchtime=2s -memprofile=mem.prof
//go tool pprof -sample_index=alloc_space  outproxy.test mem.prof

func BenchmarkRequestProxyParallel(b *testing.B) {

	httpPostReq :=
		"POST /search_service/api/service HTTP/1.0\r\n" +
			"Content-Type:application/x-www-form-urlencoded\r\n" +
			"C_ServiceName:chat_service\r\n\r\n" +
			"param=hello\r\n"

	r, err := http.ReadRequest(bufio.NewReader(strings.NewReader(httpPostReq)))
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		rw := httptest.NewRecorder()
		for pb.Next() {
			requestProxy(rw, r)
		}
	})
}

//go test -bench=Parallel -blockprofile=block.prof
//go tool pprof outproxy.test block.prof

func TestEncodeReqToGolden(t *testing.T) {

	httpPostReq :=
		"POST /search_service/api/service HTTP/1.0\r\n" +
			"Content-Type:application/x-www-form-urlencoded\r\n" +
			"C_ServiceName:chat_service\r\n\r\n" +
			"param=hello\r\n"

	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(httpPostReq)))
	if err != nil {
		t.Fatal(err)
	}

	result, erro := pack.EncodeReq(req)
	if erro != nil {
		t.Fatal(erro)
	}
	dir := "testdata"
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile("testdata/encodeReq.golden", result, 0644)
	if err != nil {
		t.Fatal(err)
	}
}
