package proxy

import (
	"bytes"
	"github.com/Kingson4Wu/fast_proxy/common/server"
	"github.com/Kingson4Wu/fast_proxy/outproxy/internal/pack"
	"github.com/Kingson4Wu/fast_proxy/outproxy/outconfig"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

var proxy *httputil.ReverseProxy

func init() {

	proxy = &httputil.ReverseProxy{

		Director: func(req *http.Request) {

			callUrl := outconfig.Get().ForwardAddress() + req.RequestURI

			u, _ := url.Parse(callUrl)
			req.URL = u
			req.Host = u.Host
		},

		ModifyResponse: func(resp *http.Response) error {

			server.GetLogger().Info("resp :", "status", resp.Status)
			server.GetLogger().Info("resp headers:")
			for hk, hv := range resp.Header {
				server.GetLogger().Info("", hk, strings.Join(hv, ","))
			}

			body, err := pack.DecodeResp(resp)
			if err != nil {
				resp.Header.Add("proxy_error_message", err.Msg)
				resp.StatusCode = err.Code
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(body))
			//resp.ContentLength = int64(len(bodyBytesDecrypt))
			resp.Header.Set("Content-Length", strconv.Itoa(len(body)))

			return nil
		},

		//ErrorLog: log.New(os.Stdout, "ReverseProxy:", log.LstdFlags|log.Lshortfile),

		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			if err != nil {
				server.GetLogger().Errorf("ErrorHandler catch err:%s", err)

				w.WriteHeader(http.StatusBadGateway)
				//_, _ = fmt.Println(w, err.Error())
				server.GetLogger().Error("", "w", w, "Encode err", err)
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
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	proxy.ServeHTTP(res, req)
}

func writeErrorMessage(res http.ResponseWriter, statusCode int, errorHeader string) {
	res.Header().Add("proxy_error_message", errorHeader)
	res.WriteHeader(statusCode)
}
