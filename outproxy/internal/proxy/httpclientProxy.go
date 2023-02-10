package proxy

import (
	"bytes"
	"context"
	"errors"
	"github.com/Kingson4Wu/fast_proxy/outproxy/config"
	"github.com/Kingson4Wu/fast_proxy/outproxy/internal/logger"
	"github.com/Kingson4Wu/fast_proxy/outproxy/internal/pack"
	"github.com/Kingson4Wu/fast_proxy/outproxy/internal/servicediscovery"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

var client *http.Client

func init() {
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

	bodyBytes, errn := pack.EncodeReq(r)
	if errn != nil {
		writeErrorMessage(w, errn.Code, errn.Msg)
		return
	}

	// 转发的URL
	reqURL := servicediscovery.GetForwardAddress() + r.RequestURI

	// 创建转发用的请求
	reqProxy, err := http.NewRequest(r.Method, reqURL, bytes.NewReader(bodyBytes))
	if err != nil {
		log.Println("创建转发请求发生错误")
		// 响应状态码
		w.WriteHeader(http.StatusServiceUnavailable)
		return
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
		reqServiceName := servicediscovery.GetServiceName(r)
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
			io.Copy(ioutil.Discard, responseProxy.Body)
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
		/*if strings.EqualFold(k, "Content-Length") {
			continue
		}*/
		w.Header().Set(k, v[0])
	}

	if responseProxy.StatusCode == http.StatusOK {
		body, errn := pack.DecodeResp(responseProxy)
		if errn != nil {
			writeErrorMessage(w, errn.Code, errn.Msg)
			return
		}

		resProxyBody := ioutil.NopCloser(bytes.NewBuffer(body))
		defer resProxyBody.Close() // 延时关闭
		// 复制转发的响应Body到响应Body
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		io.Copy(w, resProxyBody)
	} else {
		io.Copy(w, responseProxy.Body)
	}

	// 响应状态码
	//w.WriteHeader(responseProxy.StatusCode)

}
