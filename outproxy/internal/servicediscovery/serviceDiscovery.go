package servicediscovery

import (
	"github.com/Kingson4Wu/fast_proxy/outproxy/outconfig"
	"net/http"
	"strconv"
)

func GetServiceName(req *http.Request) string {
	return req.Header.Get(outconfig.Get().ServiceRpcHeaderName())
}

func GetForwardAddress() string {
	return outconfig.Get().ForwardAddress()
}

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
