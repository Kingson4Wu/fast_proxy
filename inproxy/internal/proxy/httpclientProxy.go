package proxy

import (
	"bytes"
	"context"
	"errors"
	sc "github.com/Kingson4Wu/fast_proxy/common/servicediscovery"
	"github.com/Kingson4Wu/fast_proxy/inproxy/config"
	"github.com/Kingson4Wu/fast_proxy/inproxy/internal/limiter"
	"github.com/Kingson4Wu/fast_proxy/inproxy/internal/logger"
	"github.com/Kingson4Wu/fast_proxy/inproxy/internal/pack"
	"github.com/Kingson4Wu/fast_proxy/inproxy/internal/servicediscovery"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

var client *http.Client

func init() {

	//TODO 可配置
	tr := &http.Transport{

		MaxIdleConns: 5000,

		MaxIdleConnsPerHost: 3000,
	}

	client = &http.Client{
		Transport: tr,
		//Timeout: 50 * time.Millisecond,
	}
}

// DoProxy /** 转发请求*/
func DoProxy(w http.ResponseWriter, r *http.Request) {

	//参数校验
	if !strings.HasPrefix(r.RequestURI, config.Get().ServerContextPath()) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	/** 鉴权  判断uri是否有权限等 */
	//如何防止伪造服务名(签名验证通过即认为是)
	clientServiceName := servicediscovery.GetServiceName(r)
	requestPath := servicediscovery.RealRequestUri(r.RequestURI)
	if !config.Get().ContainsCallPrivilege(clientServiceName, requestPath) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if limiter.IsLimit(clientServiceName, requestPath) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bodyBytes, pb, errn := pack.DecodeReq(r)
	defer pack.PbPool.Put(pb)
	if errn != nil {
		writeErrorMessage(w, errn.Code, errn.Msg)
		return
	}

	callUrl, rHandler := servicediscovery.Forward(r)

	if callUrl == "" {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// 创建转发用的请求
	reqProxy, err := http.NewRequest(r.Method, callUrl, bytes.NewReader(bodyBytes))
	if err != nil {
		// 响应状态码
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	if rHandler != nil {
		rHandler(reqProxy)
	}

	deadTime := int64(servicediscovery.GetRequestDeadTime(r))
	if deadTime > 0 {
		if deadTime <= time.Now().Unix() {
			w.WriteHeader(http.StatusGatewayTimeout)
			return
		} else {
			ctx, cancel := context.WithDeadline(context.Background(), time.Unix(deadTime, 0))
			defer cancel()
			reqProxy = reqProxy.WithContext(ctx)
		}
	} else {
		reqServiceName := sc.GetClientName(r)
		timeout := config.Get().GetTimeoutConfigByName(reqServiceName, r.RequestURI)
		if timeout > 0 {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
			defer cancel()
			reqProxy = reqProxy.WithContext(ctx)
		}
	}

	// 转发请求的 Header
	for k, v := range r.Header {
		reqProxy.Header.Set(k, v[0])
	}

	// 发起请求
	responseProxy, err := client.Do(reqProxy)
	if responseProxy != nil {
		defer func() {
			io.Copy(io.Discard, responseProxy.Body)
			responseProxy.Body.Close()
		}()
	}
	if err != nil {
		logger.GetLogger().Error("转发请求发生错误", zap.Any("req forward err", err))

		if errors.Is(err, context.DeadlineExceeded) {
			w.WriteHeader(http.StatusGatewayTimeout)
			return
		}

		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// 转发响应的 Header
	for k, v := range responseProxy.Header {
		if strings.EqualFold(k, "Content-Length") {
			continue
		}
		w.Header().Set(k, v[0])
	}

	body, errn := pack.EncodeResp(responseProxy, pb)
	if errn != nil {
		writeErrorMessage(w, errn.Code, errn.Msg)
		return
	}

	resProxyBody := io.NopCloser(bytes.NewBuffer(body))
	defer resProxyBody.Close() // 延时关闭

	// 响应状态码
	w.WriteHeader(responseProxy.StatusCode)
	// 复制转发的响应Body到响应Body
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	io.Copy(w, resProxyBody)

}

func writeErrorMessage(res http.ResponseWriter, statusCode int, errorHeader string) {
	res.Header().Add("proxy_error_message", errorHeader)
	res.WriteHeader(statusCode)
}
