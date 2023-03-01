package pack

import (
	"bufio"
	"fmt"
	"github.com/Kingson4Wu/fast_proxy/common/server"
	"github.com/Kingson4Wu/fast_proxy/common/servicediscovery"
	"github.com/Kingson4Wu/fast_proxy/outproxy/outconfig"
	"github.com/Kingson4Wu/fast_proxy/test"
	"github.com/agiledragon/gomonkey"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"reflect"
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

	//-gcflags "all=-N -l"
	//mock a function
	mockSc := &servicediscovery.ServiceCenter{}
	patchSc := gomonkey.ApplyFunc(server.Center, func() *servicediscovery.ServiceCenter {
		return mockSc
	})
	defer patchSc.Reset()

	//mock a method
	patchClientName := gomonkey.ApplyMethod(reflect.TypeOf(mockSc), "ClientName", func(_ *servicediscovery.ServiceCenter, req *http.Request) string {
		return req.Header.Get("C_service_name")
	})
	defer patchClientName.Reset()

	mockConfig := test.GetOutConfig()
	patchConfig := gomonkey.ApplyFunc(outconfig.Get, func() outconfig.Config {
		return mockConfig
	})
	defer patchConfig.Reset()

	httpPostReq :=
		"POST /service_name/thrift/service HTTP/1.0\r\n" +
			"protocol:json\r\n" +
			"Content-Type:application/x-thrift\r\n" +
			"C_service_name:labali\r\n\r\n" +
			"[1,\"getUserId\",1,1,{\"1\":{\"i32\":123433123}}]\r\n"

	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(httpPostReq)))
	if err != nil {
		t.Fatal(err)
	}

	result, err := EncodeReq(req)
	Convey("数据不为空", t, func() {
		So(err, ShouldBeNil)
		So(result, ShouldNotBeNil)
	})
}
