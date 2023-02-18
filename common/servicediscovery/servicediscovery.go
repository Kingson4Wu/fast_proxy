package servicediscovery

import (
	"net/http"
	"strconv"
)

type ServiceCenter struct {
	addressFunc    func(serviceName string) *Address
	clientNameFunc func(req *http.Request) string
	registerFunc   func(name string, ip string, port int) chan bool
}

func (sc *ServiceCenter) Address(serviceName string) *Address {
	return sc.addressFunc(serviceName)
}

func (sc *ServiceCenter) ClientName(req *http.Request) string {
	return sc.clientNameFunc(req)
}

func (sc *ServiceCenter) Register(name string, ip string, port int) chan bool {
	return sc.registerFunc(name, ip, port)
}

type ServiceCenterBuilder struct {
	addressFunc    func(serviceName string) *Address
	clientNameFunc func(req *http.Request) string
	registerFunc   func(name string, ip string, port int) chan bool
}

func (b *ServiceCenterBuilder) AddressFunc(f func(serviceName string) *Address) *ServiceCenterBuilder {
	b.addressFunc = f
	return b
}
func (b *ServiceCenterBuilder) ClientNameFunc(f func(req *http.Request) string) *ServiceCenterBuilder {
	b.clientNameFunc = f
	return b
}
func (b *ServiceCenterBuilder) RegisterFunc(f func(name string, ip string, port int) chan bool) *ServiceCenterBuilder {
	b.registerFunc = f
	return b
}

func (b *ServiceCenterBuilder) Build() *ServiceCenter {
	return &ServiceCenter{
		addressFunc:    b.addressFunc,
		clientNameFunc: b.clientNameFunc,
		registerFunc:   b.registerFunc,
	}
}

func Create() *ServiceCenterBuilder {
	return &ServiceCenterBuilder{}
}

type ServiceQuery interface {
	Get(name string) *Address

	ClientName(req *http.Request) string
}

type Address struct {
	Ip   string
	Port int
}

func GetRequestDeadTime(req *http.Request) int {
	timestamp := req.Header.Get("request_dead_time")
	if timestamp == "" {
		return 0
	}
	op, err := strconv.Atoi(timestamp)
	if err != nil {
		return 0
	}
	return op
}
