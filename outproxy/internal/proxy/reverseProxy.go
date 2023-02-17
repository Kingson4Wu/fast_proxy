package proxy

import (
	"bytes"
	"github.com/Kingson4Wu/fast_proxy/common/logger"
	"github.com/Kingson4Wu/fast_proxy/outproxy/internal/pack"
	"github.com/Kingson4Wu/fast_proxy/outproxy/internal/servicediscovery"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

var proxy *httputil.ReverseProxy

func init() {

	proxy = &httputil.ReverseProxy{

		Director: func(req *http.Request) {

			callUrl := servicediscovery.GetForwardAddress() + req.RequestURI

			u, _ := url.Parse(callUrl)
			req.URL = u
			req.Host = u.Host // 必须显示修改Host，否则转发可能失败
		},

		ModifyResponse: func(resp *http.Response) error {

			logger.GetLogger().Info("resp :", zap.String("status", resp.Status))
			logger.GetLogger().Info("resp headers:")
			for hk, hv := range resp.Header {
				logger.GetLogger().Info("", zap.String(hk, strings.Join(hv, ",")))
			}

			body, err := pack.DecodeResp(resp)
			if err != nil {
				resp.Header.Add("proxy_error_message", err.Msg)
				resp.StatusCode = err.Code
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			//resp.ContentLength = int64(len(bodyBytesDecrypt))
			resp.Header.Set("Content-Length", strconv.Itoa(len(body)))

			return nil
		},

		//ErrorLog: log.New(os.Stdout, "ReverseProxy:", log.LstdFlags|log.Lshortfile),

		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			if err != nil {
				logger.GetLogger().Error("ErrorHandler catch err:", zap.Error(err))

				w.WriteHeader(http.StatusBadGateway)
				//_, _ = fmt.Println(w, err.Error())
				logger.GetLogger().Error("", zap.Any("w", w), zap.Any("Encode err", err))
			}
		},
	}
}

func Delegate(res http.ResponseWriter, req *http.Request) {
	bodyBytes, err := pack.EncodeReq(req)
	if err != nil {
		writeErrorMessage(res, err.Code, err.Msg)
		return
	}

	req.ContentLength = int64(len(bodyBytes))
	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	proxy.ServeHTTP(res, req)
}

func writeErrorMessage(res http.ResponseWriter, statusCode int, errorHeader string) {
	res.Header().Add("proxy_error_message", errorHeader)
	res.WriteHeader(statusCode)
}
