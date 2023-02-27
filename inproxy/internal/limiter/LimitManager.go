package limiter

import (
	"golang.org/x/time/rate"
	"sync"
)

var limitMap map[string]*rate.Limiter
var keyLocks sync.Map

func init() {
	// https://cloud.tencent.com/developer/article/1847918
	// limiter := rate.NewLimiter(10, 100)

	limitMap = make(map[string]*rate.Limiter)
}

func IsLimit(serviceName string, uri string) bool {
	key := serviceName + "_" + uri

	limiter, ok := limitMap[key]
	if ok {
		return limiter.Allow()
	}

	keyLock, _ := keyLocks.LoadOrStore(key, &sync.Mutex{})
	keyLock.(*sync.Mutex).Lock()
	defer keyLock.(*sync.Mutex).Unlock()

	limiter, ok = limitMap[key]
	if ok {
		return limiter.Allow()
	}
	qps := 10
	limiter = rate.NewLimiter(rate.Limit(qps), qps)
	limitMap[key] = limiter

	return limiter.Allow()
}
