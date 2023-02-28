package proxy

import (
	"bytes"
	"context"
	"errors"
	"github.com/Kingson4Wu/fast_proxy/common/config"
	"github.com/Kingson4Wu/fast_proxy/common/server"
	"github.com/Kingson4Wu/fast_proxy/common/servicediscovery"
	"github.com/Kingson4Wu/fast_proxy/outproxy/internal/pack"
	"github.com/Kingson4Wu/fast_proxy/outproxy/outconfig"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
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

	bodyBytes, error := pack.EncodeReq(r)
	if error != nil {
		writeErrorMessage(w, error.Code, error.Msg)
		return
	}

	reqURL := outconfig.Get().ForwardAddress() + r.RequestURI

	reqProxy, err := http.NewRequest(r.Method, reqURL, bytes.NewReader(bodyBytes))
	if err != nil {
		log.Println("Error creating forward request")
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
		reqServiceName := server.Center().ClientName(r)
		timeout := outconfig.Get().GetTimeoutConfigByName(reqServiceName, r.RequestURI)
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
		server.GetLogger().Error("Error forwarding request", "req forward err", err)

		if errors.Is(err, context.DeadlineExceeded) {
			w.WriteHeader(http.StatusGatewayTimeout)
			return
		}

		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// Header of the forwarded response
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

		resProxyBody := io.NopCloser(bytes.NewBuffer(body))
		defer resProxyBody.Close() // Delay off
		// Copy the forwarded response Body to the response Body
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		io.Copy(w, resProxyBody)
	} else {
		io.Copy(w, responseProxy.Body)
		w.WriteHeader(responseProxy.StatusCode)
	}

	// response status code
	//w.WriteHeader(responseProxy.StatusCode)

}
