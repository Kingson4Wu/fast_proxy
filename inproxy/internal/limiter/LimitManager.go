package limiter

import "golang.org/x/time/rate"

var limitMap map[string]*rate.Limiter

func init() {
	//TODO 流量控制
	//https://cloud.tencent.com/developer/article/1847918
	//limiter := rate.NewLimiter(10, 100)
}

func IsLimit(serviceName string, uri string) bool {

	limiter, ok := limitMap[serviceName+"_"+uri]
	if ok {
		return limiter.Allow()
	}
	return false

}
