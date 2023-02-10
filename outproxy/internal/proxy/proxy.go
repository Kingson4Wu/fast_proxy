package proxy

import (
	"net/http"
)

var (
	reverseProxy Func = Func(ReverseProxy)
	forwardProxy Func = Func(ForwardProxy)
)

type Type int

const (
	REVERSE Type = iota
	FORWARD
)

func GetProxy(t Type) Func {

	switch t {
	case REVERSE:
		/* once.Do(func() {
			reverseProxy = ProxyFunc(ReverseProxy)
		})   */
		return reverseProxy
	case FORWARD:
		return forwardProxy
	}

	return nil
}

type Handler interface {
	Handle(http.ResponseWriter, *http.Request)
}
type Func func(http.ResponseWriter, *http.Request)

func (f Func) Handle(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

func ReverseProxy(w http.ResponseWriter, r *http.Request) {
	Delegate(w, r)
}

func ForwardProxy(w http.ResponseWriter, r *http.Request) {
	DoProxy(w, r)
}
