package servicediscovery

import (
	"github.com/Kingson4Wu/fast_proxy/common/server"
	"github.com/Kingson4Wu/fast_proxy/common/servicediscovery"
	"github.com/Kingson4Wu/fast_proxy/inproxy/inconfig"
	"net/http"
	"strconv"
	"strings"
)

func GetRequestDeadTime(req *http.Request) int {
	return servicediscovery.GetRequestDeadTime(req)
}

func RealRequestUri(uri string) string {
	return strings.Replace(uri, inconfig.Get().ServerContextPath(), "", 1)
}

func Forward(r *http.Request) (string, func(*http.Request)) {
	requestPath := RealRequestUri(r.RequestURI)

	arr := strings.Split(requestPath, "/")
	serviceName := arr[1]
	requestPath = strings.Replace(requestPath, "/"+serviceName, "", 1)

	var callUrl string
	var rHandler func(*http.Request)

	addr := server.Center().Address(serviceName)
	if addr == nil {
		return "", nil
	}

	builder := strings.Builder{}
	builder.WriteString("http://")
	builder.WriteString(addr.Ip)
	builder.WriteString(":")
	builder.WriteString(strconv.Itoa(addr.Port))
	builder.WriteString(requestPath)
	callUrl = builder.String()

	if len(callUrl) > 0 {
		return callUrl, rHandler
	}

	return "", nil
}
