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
    "github.com/valyala/fasthttp"
    "io"
    "log"
    "math/rand"
    "net"
    "net/http"
    "os"
    "strings"
    "strconv"
    "time"
)

var client *http.Client
var fastHttpClient *fasthttp.Client

// hop-by-hop headers that should not be forwarded between client and upstream
func isHopByHop(h string) bool {
    switch strings.ToLower(h) {
    case "connection", "proxy-connection", "keep-alive", "te", "trailer", "transfer-encoding", "upgrade":
        return true
    default:
        return false
    }
}

var (
    enableRetries = true
    maxRetries    = 2
    baseBackoff   = 100 * time.Millisecond
    maxBackoff    = 1 * time.Second
)

func initRetryConfigFromEnv() {
    if v := os.Getenv("FAST_PROXY_RETRY_ENABLED"); v != "" {
        enableRetries = strings.EqualFold(v, "1") || strings.EqualFold(v, "true") || strings.EqualFold(v, "yes")
    }
    if v := os.Getenv("FAST_PROXY_RETRY_MAX"); v != "" {
        if n, err := strconv.Atoi(v); err == nil && n >= 0 {
            maxRetries = n
        }
    }
    if v := os.Getenv("FAST_PROXY_RETRY_BASE_MS"); v != "" {
        if n, err := strconv.Atoi(v); err == nil && n > 0 {
            baseBackoff = time.Duration(n) * time.Millisecond
        }
    }
    if v := os.Getenv("FAST_PROXY_RETRY_MAX_MS"); v != "" {
        if n, err := strconv.Atoi(v); err == nil && n > 0 {
            maxBackoff = time.Duration(n) * time.Millisecond
        }
    }
}

func shouldRetryErr(err error) bool {
    if err == nil {
        return false
    }
    // Don't retry on context deadline/cancellation
    if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
        return false
    }
    // Retry on temporary/timeouts at the transport level
    var ne net.Error
    if errors.As(err, &ne) {
        return ne.Timeout() || ne.Temporary()
    }
    // Default: no
    return false
}

func shouldRetryStatus(code int) bool {
    switch code {
    case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
        return true
    default:
        return false
    }
}

func backoffDuration(attempt int) time.Duration {
    d := baseBackoff << attempt
    if d > maxBackoff {
        d = maxBackoff
    }
    // add jitter +/- 20%
    jitter := 0.2 * float64(d)
    delta := (rand.Float64()*2 - 1) * jitter
    return d + time.Duration(delta)
}

func BuildClient(c config.Config) {
    initRetryConfigFromEnv()
    tr := &http.Transport{

		MaxIdleConns: c.HttpClientMaxIdleConns(),

		MaxIdleConnsPerHost: c.HttpClientMaxIdleConnsPerHost(),
	}

	client = &http.Client{
		Transport: tr,
		//Timeout: 50 * time.Millisecond,
	}

	fastHttpClient = &fasthttp.Client{
		MaxConnsPerHost: c.HttpClientMaxIdleConnsPerHost(),
		//MaxIdleConnDuration: c.HttpClientMaxIdleConns(),
		//WriteTimeout: 50 * time.Millisecond,
		//ReadTimeout:  50 * time.Millisecond,
	}
}

// Close releases resources held by underlying clients (idle connections, etc.).
func Close() {
    if client != nil {
        if tr, ok := client.Transport.(*http.Transport); ok {
            tr.CloseIdleConnections()
        }
    }
    // fasthttp client keeps connections in pools; nothing explicit required here
}

// DoProxy /** forward request */
func DoProxy(w http.ResponseWriter, r *http.Request) {

	if server.Config().FastHttpEnable() {
		fastDoProxy(w, r)
		return
	}

    bodyBytes, cerr := pack.EncodeReq(r)
    if cerr != nil {
        writeErrorMessage(w, cerr.Code, cerr.Msg)
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

    // Header of forwarding request (skip hop-by-hop)
    for k, v := range r.Header {
        if isHopByHop(k) {
            continue
        }
        if len(v) > 0 {
            reqProxy.Header.Set(k, v[0])
        }
    }

    var responseProxy *http.Response
    attempts := 0
    for {
        // make a request
        responseProxy, err = client.Do(reqProxy)
        if err != nil {
            if enableRetries && attempts < maxRetries && shouldRetryErr(err) && r.Context().Err() == nil {
                attempts++
                time.Sleep(backoffDuration(attempts - 1))
                // rebuild request body reader for next attempt
                reqProxy.Body = io.NopCloser(bytes.NewReader(bodyBytes))
                continue
            }
            server.GetLogger().Error("Error forwarding request", "req forward err", err)
            if errors.Is(err, context.DeadlineExceeded) {
                w.WriteHeader(http.StatusGatewayTimeout)
                return
            }
            w.WriteHeader(http.StatusServiceUnavailable)
            return
        }

        if responseProxy != nil && enableRetries && attempts < maxRetries && shouldRetryStatus(responseProxy.StatusCode) && r.Context().Err() == nil {
            // drain and close before retry
            _, _ = io.Copy(io.Discard, responseProxy.Body)
            responseProxy.Body.Close()
            attempts++
            time.Sleep(backoffDuration(attempts - 1))
            reqProxy.Body = io.NopCloser(bytes.NewReader(bodyBytes))
            continue
        }
        break
    }

    if responseProxy != nil {
        defer func() {
            _, _ = io.Copy(io.Discard, responseProxy.Body)
            responseProxy.Body.Close()
        }()
    }

    // Header of the forwarded response (skip hop-by-hop)
    for k, v := range responseProxy.Header {
        if isHopByHop(k) {
            continue
        }
        if len(v) > 0 {
            w.Header().Set(k, v[0])
        }
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

		if _, err := io.Copy(w, resProxyBody); err != nil {
			server.GetLogger().Error("Error forwarding request", "io.Copy err", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
	} else {
		if _, err := io.Copy(w, responseProxy.Body); err != nil {
			server.GetLogger().Error("Error forwarding request", "io.Copy err", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(responseProxy.StatusCode)
	}

	// response status code
	//w.WriteHeader(responseProxy.StatusCode)

}

func fastDoProxy(w http.ResponseWriter, r *http.Request) {

	bodyBytes, erro := pack.EncodeReq(r)
	if erro != nil {
		writeErrorMessage(w, erro.Code, erro.Msg)
		return
	}

	reqURL := outconfig.Get().ForwardAddress() + r.RequestURI

	// Create a new request with fasthttp
	reqProxy := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(reqProxy)
	reqProxy.SetRequestURI(reqURL)
	reqProxy.Header.SetMethod(r.Method)

    // Copy request headers to fasthttp request (skip hop-by-hop)
    headers := &reqProxy.Header
    for k, v := range r.Header {
        if isHopByHop(k) {
            continue
        }
        if len(v) > 0 {
            headers.Set(k, v[0])
        }
    }

	reqProxy.SetBody(bodyBytes)

	resProxy := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resProxy)

    var err error
    
    deadTime := int64(servicediscovery.GetRequestDeadTime(r))
    reqServiceName := server.Center().ClientName(r)
    timeout := outconfig.Get().GetTimeoutConfigByName(reqServiceName, r.RequestURI)
    attempts := 0
    for {
        if deadTime > 0 {
            if deadTime <= time.Now().Unix() {
                w.WriteHeader(http.StatusGatewayTimeout)
                return
            } else {
                err = fastHttpClient.DoDeadline(reqProxy, resProxy, time.Unix(deadTime, 0))
            }
        } else if timeout > 0 {
            err = fastHttpClient.DoTimeout(reqProxy, resProxy, time.Duration(timeout)*time.Millisecond)
        } else {
            err = fastHttpClient.Do(reqProxy, resProxy)
        }

        if err != nil {
            if enableRetries && attempts < maxRetries && shouldRetryErr(err) && r.Context().Err() == nil {
                attempts++
                time.Sleep(backoffDuration(attempts - 1))
                // resProxy may contain body; reset for next attempt
                resProxy.Reset()
                continue
            }
            server.GetLogger().Error("Error forwarding request", "req forward err", err)
            if errors.Is(err, context.DeadlineExceeded) {
                w.WriteHeader(http.StatusGatewayTimeout)
                return
            }
            w.WriteHeader(http.StatusServiceUnavailable)
            return
        }

        if enableRetries && attempts < maxRetries && shouldRetryStatus(resProxy.StatusCode()) && r.Context().Err() == nil {
            attempts++
            time.Sleep(backoffDuration(attempts - 1))
            resProxy.Reset()
            continue
        }
        break
    }

    // Set response headers (skip hop-by-hop)
    resHeader := w.Header()
    resProxy.Header.VisitAll(func(k, v []byte) {
        if isHopByHop(string(k)) {
            return
        }
        resHeader.Set(string(k), string(v))
    })

	if resProxy.StatusCode() == http.StatusOK {
		body, errn := pack.DecodeFastResp(resProxy.Body())
		if errn != nil {
			writeErrorMessage(w, errn.Code, errn.Msg)
			return
		}

		resProxyBody := io.NopCloser(bytes.NewBuffer(body))
		defer resProxyBody.Close() // Delay off
		// Copy the forwarded response Body to the response Body
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		if _, err := io.Copy(w, resProxyBody); err != nil {
			server.GetLogger().Error("Error forwarding request", "io.Copy err", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
	} else {
		if _, err := w.Write(resProxy.Body()); err != nil {
			server.GetLogger().Error("Error forwarding request", "w.Write err", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(resProxy.StatusCode())
	}
}
