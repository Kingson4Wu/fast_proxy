package servicediscovery

import "net/http"

type ServiceQuery interface {
	Get(name string) *Address

	//ClientName(req *http.Request) string
}

type Address struct {
	Ip   string
	Port int
}

/** todo struct ，在各自proxy组合 */
var Sq ServiceQuery

func GetAddress(name string) *Address {
	if Sq != nil {
		return Sq.Get(name)
	}
	return nil
}

func GetClientName(r *http.Request) string {
	if Sq != nil {
		//return Sq.ClientName(r)
		return "song_service"
	}
	return ""
}

func Register(s ServiceQuery) {
	Sq = s
}
