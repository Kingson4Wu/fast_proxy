package pack

import (
	"bufio"
	"bytes"
	"github.com/Kingson4Wu/fast_proxy/common/compress"
	"github.com/Kingson4Wu/fast_proxy/common/encrypt"
	"github.com/Kingson4Wu/fast_proxy/common/proto/protobuf"
	"github.com/Kingson4Wu/fast_proxy/common/server"
	"github.com/Kingson4Wu/fast_proxy/common/servicediscovery"
	"github.com/Kingson4Wu/fast_proxy/outproxy/outconfig"
	"github.com/Kingson4Wu/fast_proxy/test"
	"github.com/agiledragon/gomonkey"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func BenchmarkEncodeReq(b *testing.B) {

	c := runBeforeMock()
	defer c()

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
	for i := 0; i < b.N; i++ {
		EncodeReq(req)
	}

	//Benchmark test generates pprof file analysis
	//go test ./outproxy/internal/pack -bench="^\QBenchmarkEncodeReq\E$" -gcflags "all=-N -l" -memprofile=mem.pprof
}

func TestEncodeReq(t *testing.T) {

	c := runBeforeMock()
	defer c()

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
	Convey("data not null", t, func() {
		So(err, ShouldBeNil)
		So(result, ShouldNotBeNil)
	})
}

func runBeforeMock() func() {
	//-gcflags "all=-N -l"
	//mock a function
	mockSc := &servicediscovery.ServiceCenter{}
	patchSc := gomonkey.ApplyFunc(server.Center, func() *servicediscovery.ServiceCenter {
		return mockSc
	})

	//mock a method
	patchClientName := gomonkey.ApplyMethod(reflect.TypeOf(mockSc), "ClientName", func(_ *servicediscovery.ServiceCenter, req *http.Request) string {
		return req.Header.Get("C_service_name")
	})

	mockConfig := test.GetOutConfig()
	patchConfig := gomonkey.ApplyFunc(outconfig.Get, func() outconfig.Config {
		return mockConfig
	})

	return func() {
		patchSc.Reset()
		patchClientName.Reset()
		patchConfig.Reset()
	}
}

func TestDecodeResp(t *testing.T) {

	mockConfig := test.GetOutConfig()
	patchConfig := gomonkey.ApplyFunc(outconfig.Get, func() outconfig.Config {
		return mockConfig
	})
	defer patchConfig.Reset()

	text := "Hello, world!"
	encryptKeyName := "encrypt.key.room.v2"
	encryptKey := mockConfig.GetEncryptKeyByName(encryptKeyName)
	b, err := encrypt.Encode([]byte(text), encryptKey)

	if err != nil {
		t.Fatal(err)
	}

	b, err = compress.Encode(b, 0)
	if err != nil {
		t.Fatal(err)
	}

	stSend := &protobuf.ProxyRespData{}
	stSend.Compress = true
	stSend.Payload = b
	stSend.EncryptEnable = true
	stSend.EncryptKeyName = encryptKeyName
	stSend.CompressAlgorithm = 0

	pData, err := proto.Marshal(stSend)
	if err != nil {
		t.Fatal(err)
	}

	fakeResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(pData)),
	}

	result, err := DecodeResp(fakeResponse)
	Convey("data not null", t, func() {
		So(err, ShouldBeNil)
		So(result, ShouldNotBeNil)
		So(string(result), ShouldEqual, text)
	})
}
