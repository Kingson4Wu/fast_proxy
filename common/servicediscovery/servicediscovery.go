package servicediscovery

import "net/http"

type ServiceQuery interface {
	Get(name string) *address

	ClientName(req *http.Request) string
}

type address struct {
	Ip   string
	Port int
}

var sq ServiceQuery

func GetAddress(name string) *address {
	if sq != nil {
		return sq.Get(name)
	}
	return nil
}

func GetClientName(r *http.Request) string {
	if sq != nil {
		return sq.ClientName(r)
	}
	return ""
}

func Register(s ServiceQuery) {
	sq = s
}
