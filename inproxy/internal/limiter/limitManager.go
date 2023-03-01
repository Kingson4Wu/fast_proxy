package limiter

import (
	"github.com/Kingson4Wu/fast_proxy/inproxy/inconfig"
	"golang.org/x/time/rate"
	"sync"
)

var limitMap sync.Map
var keyLocks sync.Map

func IsLimit(serviceName string, uri string) bool {
	key := serviceName + "_" + uri

	// Concurrently read limiter from limitMap
	limiterI, ok := limitMap.Load(key)
	if ok {
		limiter := limiterI.(*rate.Limiter)
		return !limiter.Allow()
	}

	// Get the lock for the key to avoid affecting different keys
	keyLock, _ := keyLocks.LoadOrStore(key, &sync.Mutex{})
	keyLock.(*sync.Mutex).Lock()
	defer keyLock.(*sync.Mutex).Unlock()

	// Read limiter again from limitMap
	limiterI, ok = limitMap.Load(key)
	if ok {
		limiter := limiterI.(*rate.Limiter)
		return !limiter.Allow()
	}

	qps := inconfig.Get().ServiceQps(serviceName, uri)
	limiter := rate.NewLimiter(rate.Limit(qps), qps)

	// Concurrently write limiter to limitMap
	limitMap.Store(key, limiter)

	return !limiter.Allow()
}
