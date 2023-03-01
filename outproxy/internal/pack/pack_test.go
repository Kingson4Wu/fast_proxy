package pack

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func BenchmarkEncodeReq(b *testing.B) {

	httpPostReq :=
		"POST /service_name/thrift/service HTTP/1.0\r\n" +
			"protocol:json\r\n" +
			"Content-Type:application/x-thrift\r\n" +
			"C_service_name:labali\r\n\r\n" +
			"[1,\"getUserId\",1,1,{\"1\":{\"i32\":123433123}}]\r\n"

	//req, err := http.ReadRequest(bufio.NewReader(strings.NewReader("GET /hi HTTP/1.0\r\n\r\n")))
	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(httpPostReq)))
	if err != nil {
		b.Fatal(err)
	}

	//rw := httptest.NewRecorder()
	b.ResetTimer()
	/*for i := 0; i < b.N; i++ {
		EncodeReq(req)
	}*/
	fmt.Println(req)

	//TODO Benchmark test generates pprof file analysis
}

func TestEncodeReq(t *testing.T) {

	httpPostReq :=
		"POST /service_name/thrift/service HTTP/1.0\r\n" +
			"protocol:json\r\n" +
			"Content-Type:application/x-thrift\r\n" +
			"C_service_name:labali\r\n\r\n" +
			"[1,\"getUserId\",1,1,{\"1\":{\"i32\":123433123}}]\r\n"

	//req, err := http.ReadRequest(bufio.NewReader(strings.NewReader("GET /hi HTTP/1.0\r\n\r\n")))
	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(httpPostReq)))
	if err != nil {
		t.Fatal(err)
	}

	//rw := httptest.NewRecorder()
	//EncodeReq(req)
	fmt.Println(req)

	//fmt.Println(string(rw.Body.Bytes()))
}
