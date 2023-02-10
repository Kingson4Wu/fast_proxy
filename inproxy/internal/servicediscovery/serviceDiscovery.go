package servicediscovery

import (
	"github.com/Kingson4Wu/fast_proxy/common/servicediscovery"
	"github.com/Kingson4Wu/fast_proxy/inproxy/config"
	"net/http"
	"strconv"
	"strings"
)

func GetRequestDeadTime(req *http.Request) int {
	timestamp := req.Header.Get("request_request_dead_time")
	if timestamp == "" {
		return 0
	}
	op, err := strconv.Atoi(timestamp)
	if err != nil {
		return 0
	}
	return op
}

func RealRequestUri(uri string) string {
	return strings.Replace(uri, config.Get().ServerContextPath(), "", 1)
}

func Forward(r *http.Request) (string, func(*http.Request)) {
	requestPath := RealRequestUri(r.RequestURI)

	arr := strings.Split(requestPath, "/")
	serviceName := arr[1]
	requestPath = strings.Replace(requestPath, "/"+serviceName, "", 1)

	var callUrl string
	var rHandler func(*http.Request)

	addr := servicediscovery.GetAddress(serviceName)
	if addr == nil {
		return "", nil
	}

	callUrl = "http://" + addr.Ip + ":" + strconv.Itoa(addr.Port) + requestPath

	if callUrl != "" {
		return callUrl, rHandler
	}

	return "", nil

}
