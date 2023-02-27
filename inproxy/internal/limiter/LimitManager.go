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

	// 并发安全地从 limitMap 中读取 limiter
	limiterI, ok := limitMap.Load(key)
	if ok {
		limiter := limiterI.(*rate.Limiter)
		return !limiter.Allow()
	}

	// 获取 key 对应的锁，避免不同 key 互相影响
	keyLock, _ := keyLocks.LoadOrStore(key, &sync.Mutex{})
	keyLock.(*sync.Mutex).Lock()
	defer keyLock.(*sync.Mutex).Unlock()

	// 再次从 limitMap 中读取 limiter
	limiterI, ok = limitMap.Load(key)
	if ok {
		limiter := limiterI.(*rate.Limiter)
		return !limiter.Allow()
	}

	qps := inconfig.Get().ServiceQps(serviceName, uri)
	limiter := rate.NewLimiter(rate.Limit(qps), qps)

	// 并发安全地向 limitMap 中写入 limiter
	limitMap.Store(key, limiter)

	return !limiter.Allow()
}
