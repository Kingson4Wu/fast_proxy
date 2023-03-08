package proxy

import (
	"bytes"
	"context"
	"errors"
	"github.com/Kingson4Wu/fast_proxy/common/config"
	"github.com/Kingson4Wu/fast_proxy/common/server"
	"github.com/Kingson4Wu/fast_proxy/inproxy/inconfig"
	"github.com/Kingson4Wu/fast_proxy/inproxy/internal/limiter"
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

func BuildClient(c config.Config) {
	tr := &http.Transport{

		MaxIdleConns: c.HttpClientMaxIdleConns(),

		MaxIdleConnsPerHost: c.HttpClientMaxIdleConnsPerHost(),
	}

	client = &http.Client{
		Transport: tr,
		//Timeout: 50 * time.Millisecond,
	}
}

// DoProxy /** forward request */
func DoProxy(w http.ResponseWriter, r *http.Request) {

	//Parameter check
	if !strings.HasPrefix(r.RequestURI, inconfig.Get().ServerContextPath()) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	/** Authentication to determine whether the uri has authority, etc. */
	//How to prevent forgery of service name (signature verification is considered to be)
	clientServiceName := server.Center().ClientName(r)
	requestPath := servicediscovery.RealRequestUri(r.RequestURI)
	if !inconfig.Get().ContainsCallPrivilege(clientServiceName, requestPath) {
		writeErrorMessage(w, http.StatusBadRequest, "client has no privilege")
		return
	}

	if limiter.IsLimit(clientServiceName, requestPath) {
		writeErrorMessage(w, http.StatusBadRequest, "client is limit")
		return
	}

	bodyBytes, pb, error := pack.DecodeReq(r)
	defer pack.PbPool.Put(pb)
	if error != nil {
		writeErrorMessage(w, error.Code, error.Msg)
		return
	}

	callUrl, rHandler := servicediscovery.Forward(r)

	if callUrl == "" {
		writeErrorMessage(w, http.StatusServiceUnavailable, "call url is blank")
		return
	}

	// Create a request for forwarding
	reqProxy, err := http.NewRequest(r.Method, callUrl, bytes.NewReader(bodyBytes))
	if err != nil {
		writeErrorMessage(w, http.StatusServiceUnavailable, "request wrap error")
		return
	}
	if rHandler != nil {
		rHandler(reqProxy)
	}

	deadTime := int64(servicediscovery.GetRequestDeadTime(r))
	if deadTime > 0 {
		if deadTime <= time.Now().Unix() {
			writeErrorMessage(w, http.StatusGatewayTimeout, "already reach dead time")
			return
		} else {
			ctx, cancel := context.WithDeadline(context.Background(), time.Unix(deadTime, 0))
			defer cancel()
			reqProxy = reqProxy.WithContext(ctx)
		}
	} else {
		reqServiceName := server.Center().ClientName(r)
		timeout := inconfig.Get().GetTimeoutConfigByName(reqServiceName, r.RequestURI)
		if timeout > 0 {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
			defer cancel()
			reqProxy = reqProxy.WithContext(ctx)
		}
	}

	// Header of forwarding request
	for k, v := range r.Header {
		reqProxy.Header.Set(k, v[0])
	}

	// make a request
	responseProxy, err := client.Do(reqProxy)
	if responseProxy != nil {
		defer func() {
			io.Copy(io.Discard, responseProxy.Body)
			responseProxy.Body.Close()
		}()
	}
	if err != nil {
		server.GetLogger().Error("Error forwarding request", zap.Any("req forward err", err))

		if errors.Is(err, context.DeadlineExceeded) {
			writeErrorMessage(w, http.StatusGatewayTimeout, "call timeout")
			return
		}
		writeErrorMessage(w, http.StatusServiceUnavailable, "call error")
		return
	}

	// Header of the forwarded response
	for k, v := range responseProxy.Header {
		if strings.EqualFold(k, "Content-Length") {
			continue
		}
		w.Header().Set(k, v[0])
	}

	body, error := pack.EncodeResp(responseProxy, pb)
	if error != nil {
		writeErrorMessage(w, error.Code, error.Msg)
		return
	}

	resProxyBody := io.NopCloser(bytes.NewBuffer(body))
	defer resProxyBody.Close()

	// response status code
	w.WriteHeader(responseProxy.StatusCode)
	// Copy the forwarded response Body to the response Body
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	io.Copy(w, resProxyBody)

}

func writeErrorMessage(res http.ResponseWriter, statusCode int, errorHeader string) {
	res.Header().Add("proxy_error_message", errorHeader)
	res.Header().Add("proxy_name", "in_proxy")
	res.WriteHeader(statusCode)
}
