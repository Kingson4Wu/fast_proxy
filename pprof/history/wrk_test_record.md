
<pre>
8 GB 1867 MHz DDR3
2.7 GHz Dual-Core Intel Core i5
</pre>

+ mac : `sysctl -n hw.logicalcpu` : 4

----

# v1
+ `git rev-parse --short=10 HEAD`
+ `97b1353554`

## wrk

<pre>
 /usr/local/Cellar/wrk/4.1.0/bin/wrk -d 60 -c 100  -t 32 -s post.lua http://127.0.0.1:8034/search_service/api/service
Running 1m test @ http://127.0.0.1:8034/search_service/api/service
  32 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    35.19ms   31.61ms 396.86ms   85.74%
    Req/Sec    98.68     33.87   650.00     72.90%
  187581 requests in 1.00m, 33.27MB read
Requests/sec:   3122.39
Transfer/sec:    567.15KB

Running 1m test @ http://127.0.0.1:8034/search_service/api/service
  32 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    32.59ms   25.60ms 397.52ms   78.12%
    Req/Sec   101.79     32.62   450.00     75.28%
  194562 requests in 1.00m, 34.51MB read
Requests/sec:   3237.24
Transfer/sec:    588.02KB

 /usr/local/Cellar/wrk/4.1.0/bin/wrk -d 60 -c 256  -t 64 -s post.lua http://127.0.0.1:8034/search_service/api/service
Running 1m test @ http://127.0.0.1:8034/search_service/api/service
  64 threads and 256 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    76.49ms   46.12ms 465.86ms   74.95%
    Req/Sec    54.81     19.19   303.00     74.83%
  209396 requests in 1.00m, 37.14MB read
Requests/sec:   3485.38
Transfer/sec:    633.09KB

 /usr/local/Cellar/wrk/4.1.0/bin/wrk -d 60 -c 128  -t 8 -s post.lua http://127.0.0.1:8034/search_service/api/service
Running 1m test @ http://127.0.0.1:8034/search_service/api/service
  8 threads and 128 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    39.80ms   28.94ms 343.62ms   77.02%
    Req/Sec   437.55     91.18     1.05k    71.82%
  209224 requests in 1.00m, 37.11MB read
Requests/sec:   3481.19
Transfer/sec:    632.33KB

Running 1m test @ http://127.0.0.1:8034/search_service/api/service
  4 threads and 128 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    37.72ms   25.78ms 308.25ms   75.69%
    Req/Sec     0.91k   148.38     1.43k    71.45%
  216784 requests in 1.00m, 38.45MB read
Requests/sec:   3611.41
Transfer/sec:    655.98KB

 /usr/local/Cellar/wrk/4.1.0/bin/wrk -d 60 -c 256  -t 4 -s post.lua http://127.0.0.1:8034/search_service/api/service
Running 1m test @ http://127.0.0.1:8034/search_service/api/service
  4 threads and 256 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    75.53ms   45.00ms 425.92ms   74.59%
    Req/Sec     0.88k   192.90     1.52k    68.68%
  211149 requests in 1.00m, 37.45MB read
Requests/sec:   3515.03
Transfer/sec:    638.47KB

/usr/local/Cellar/wrk/4.1.0/bin/wrk -d 60 -c 64  -t 4 -s post.lua http://127.0.0.1:8034/search_service/api/service
Running 1m test @ http://127.0.0.1:8034/search_service/api/service
  4 threads and 64 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    23.90ms   21.50ms 294.76ms   83.90%
    Req/Sec   784.45    161.79     1.32k    70.20%
  187577 requests in 1.00m, 33.27MB read
Requests/sec:   3121.13
Transfer/sec:    566.92KB

</pre>

+ At first glance, the best test result is achieved with 4 threads and 128 connections.
+ qps: 3500


<pre>

top -pid $(ps -ef | grep outproxy | grep -v 'grep' | awk '{print $2}')

Processes: 454 total, 151 running, 303 sleeping, 2351 threads          18:39:48
Load Avg: 131.83, 50.24, 30.50  CPU usage: 56.62% user, 41.85% sys, 1.52% idle
SharedLibs: 271M resident, 47M data, 78M linkedit.
MemRegions: 407272 total, 3061M resident, 66M private, 860M shared.
PhysMem: 8115M used (2131M wired), 76M unused.
VM: 45T vsize, 3100M framework vsize, 4544797007(318) swapins, 4576659064(0) swa
Networks: packets: 68524825/36G in, 60182771/18G out.
Disks: 193858870/18T read, 114611261/18T written.

PID   COMMAND      %CPU  TIME     #TH  #WQ  #POR MEM   PURG CMPRS  PGRP PPID
5917  ___go_build_ 70.4  09:25.41 16/3 0    26   16M-  0B   1004K  5917 2364

stop test

PID   COMMAND      %CPU  TIME     #TH  #WQ  #POR MEM    PURG CMPRS  PGRP PPID
5917  ___go_build_ 0.0   09:39.83 16   0    26   9148K  0B   1016K  5917 2364

-----

top -pid $(ps -ef | grep inproxy | grep -v 'grep' | awk '{print $2}')

Processes: 446 total, 98 running, 348 sleeping, 2287 threads           18:38:34
Load Avg: 59.10, 21.95, 19.70  CPU usage: 54.2% user, 44.0% sys, 1.97% idle
SharedLibs: 270M resident, 47M data, 78M linkedit.
MemRegions: 406809 total, 3089M resident, 63M private, 835M shared.
PhysMem: 8146M used (2122M wired), 45M unused.
VM: 45T vsize, 3100M framework vsize, 4544722384(767) swapins, 4576587767(0) swa
Networks: packets: 66749686/36G in, 58408675/18G out.
Disks: 193856436/18T read, 114609069/18T written.

PID   COMMAND      %CPU  TIME     #TH  #WQ  #POR MEM  PURG CMPRS PGRP PPID
5729  ___go_build_ 79.7  10:49.36 17/4 0    27   31M+ 0B   976K  5729 2364

stop test

PID   COMMAND      %CPU  TIME     #TH  #WQ  #POR MEM  PURG CMPRS PGRP PPID
5729  ___go_build_ 0.0   11:57.72 17   0    27   15M  0B   992K  5729 2364

</pre>

+ wrk -> outproxy(fasthttp client, fasthttp server) ->  outproxy(golang http client, fasthttp server) -> server

+ Why in proxy use more memory than out proxy ?
+ Why is there a lot of memory not released after stopping the in proxy test ?
+ Use fasthttp client , and fix it.


## pprof

### CPU
+ in : `go tool pprof -http :8081 http://localhost:8033/debug/pprof/profile`
+ out : `go tool pprof -http :8082 http://localhost:8034/debug/pprof/profile`

### Memory
+ in : `go tool pprof -http :8081 http://localhost:8033/debug/pprof/heap`
+ out : `go tool pprof -http :8082 http://localhost:8034/debug/pprof/heap`


## Benchmark
+ cd outproxy

### cpu

+ go test -v -run=^$ -bench=^BenchmarkRequestProxy$ -benchtime=2s -cpuprofile=cpu.prof

<pre>
goos: darwin
goarch: amd64
pkg: github.com/Kingson4Wu/fast_proxy/outproxy
cpu: Intel(R) Core(TM) i5-5257U CPU @ 2.70GHz
BenchmarkRequestProxy
BenchmarkRequestProxy-4             3576            688034 ns/op           10726 B/op        242 allocs/op
PASS
ok      github.com/Kingson4Wu/fast_proxy/outproxy       8.672s


</pre>

<pre>
go tool pprof outproxy.test cpu.prof 
File: outproxy.test
Type: cpu
Time: Mar 9, 2023 at 9:52am (CST)
Duration: 2.90s, Total samples = 840ms (28.92%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top -cum
Showing nodes accounting for 20ms, 2.38% of 840ms total
Showing top 10 nodes out of 152
      flat  flat%   sum%        cum   cum%
         0     0%     0%      660ms 78.57%  github.com/Kingson4Wu/fast_proxy/outproxy.BenchmarkRequestProxy
         0     0%     0%      660ms 78.57%  github.com/Kingson4Wu/fast_proxy/outproxy.requestProxy
         0     0%     0%      660ms 78.57%  github.com/Kingson4Wu/fast_proxy/outproxy/internal/proxy.DoProxy
         0     0%     0%      660ms 78.57%  github.com/Kingson4Wu/fast_proxy/outproxy/internal/proxy.ForwardProxy
         0     0%     0%      660ms 78.57%  github.com/Kingson4Wu/fast_proxy/outproxy/internal/proxy.Func.Handle (inline)
      20ms  2.38%  2.38%      660ms 78.57%  github.com/Kingson4Wu/fast_proxy/outproxy/internal/proxy.fastDoProxy
         0     0%  2.38%      660ms 78.57%  testing.(*B).launch
         0     0%  2.38%      660ms 78.57%  testing.(*B).runN
         0     0%  2.38%      480ms 57.14%  github.com/valyala/fasthttp.(*Client).Do
         0     0%  2.38%      480ms 57.14%  github.com/valyala/fasthttp.(*HostClient).Do

</pre>

<pre>
(pprof) list fastDoProxy
Total: 840ms
ROUTINE ======================== github.com/Kingson4Wu/fast_proxy/outproxy/internal/proxy.fastDoProxy in /Users/kingsonwu/Personal/go-src/f_proxy/outproxy/internal/proxy/httpclientProxy.go
      20ms      660ms (flat, cum) 78.57% of Total
         .          .    140:
         .          .    141:}
         .          .    142:
         .          .    143:func fastDoProxy(w http.ResponseWriter, r *http.Request) {
         .          .    144:
         .       60ms    145:   bodyBytes, erro := pack.EncodeReq(r)
         .          .    146:   if erro != nil {
         .          .    147:           writeErrorMessage(w, erro.Code, erro.Msg)
         .          .    148:           return
         .          .    149:   }
         .          .    150:
         .       10ms    151:   reqURL := outconfig.Get().ForwardAddress() + r.RequestURI
         .          .    152:
         .          .    153:   // Create a new request with fasthttp
         .          .    154:   reqProxy := fasthttp.AcquireRequest()
         .          .    155:   defer fasthttp.ReleaseRequest(reqProxy)
         .          .    156:   reqProxy.SetRequestURI(reqURL)
         .          .    157:   reqProxy.Header.SetMethod(r.Method)
         .          .    158:
         .          .    159:   // Copy request headers to fasthttp request
         .          .    160:   headers := &reqProxy.Header
         .          .    161:
         .       10ms    162:   for k, v := range r.Header {
         .       10ms    163:           headers.Set(k, v[0])
         .          .    164:   }
         .          .    165:
         .          .    166:   reqProxy.SetBody(bodyBytes)
         .          .    167:
         .          .    168:   resProxy := fasthttp.AcquireResponse()
         .          .    169:   defer fasthttp.ReleaseResponse(resProxy)
         .          .    170:
         .          .    171:   var err error
         .          .    172:
         .          .    173:   deadTime := int64(servicediscovery.GetRequestDeadTime(r))
         .          .    174:   reqServiceName := server.Center().ClientName(r)
         .          .    175:   timeout := outconfig.Get().GetTimeoutConfigByName(reqServiceName, r.RequestURI)
         .          .    176:   if deadTime > 0 {
         .          .    177:           if deadTime <= time.Now().Unix() {
         .          .    178:                   w.WriteHeader(http.StatusGatewayTimeout)
         .          .    179:                   return
         .          .    180:           } else {
         .          .    181:                   err = fastHttpClient.DoDeadline(reqProxy, resProxy, time.Unix(deadTime, 0))
         .          .    182:           }
         .          .    183:   } else if timeout > 0 {
         .          .    184:           err = fastHttpClient.DoTimeout(reqProxy, resProxy, time.Duration(timeout)*time.Millisecond)
         .          .    185:   } else {
         .      480ms    186:           err = fastHttpClient.Do(reqProxy, resProxy)
         .          .    187:   }
         .          .    188:
         .          .    189:   if err != nil {
         .          .    190:           server.GetLogger().Error("Error forwarding request", "req forward err", err)
         .          .    191:
         .          .    192:           if errors.Is(err, context.DeadlineExceeded) {
         .          .    193:                   w.WriteHeader(http.StatusGatewayTimeout)
         .          .    194:                   return
         .          .    195:           }
         .          .    196:
         .          .    197:           w.WriteHeader(http.StatusServiceUnavailable)
         .          .    198:           return
         .          .    199:   }
         .          .    200:
         .          .    201:   // Set response headers
         .          .    202:   resHeader := w.Header()
      10ms       10ms    203:   resProxy.Header.VisitAll(func(k, v []byte) {
         .          .    204:           resHeader.Set(string(k), string(v))
         .          .    205:   })
         .          .    206:
         .          .    207:   if resProxy.StatusCode() == http.StatusOK {
         .       70ms    208:           body, errn := pack.DecodeFastResp(resProxy.Body())
         .          .    209:           if errn != nil {
         .          .    210:                   writeErrorMessage(w, errn.Code, errn.Msg)
         .          .    211:                   return
         .          .    212:           }
         .          .    213:
         .          .    214:           resProxyBody := io.NopCloser(bytes.NewBuffer(body))
      10ms       10ms    215:           defer resProxyBody.Close() // Delay off
         .          .    216:           // Copy the forwarded response Body to the response Body
         .          .    217:           w.Header().Set("Content-Length", strconv.Itoa(len(body)))
         .          .    218:           io.Copy(w, resProxyBody)
         .          .    219:   } else {
         .          .    220:           w.Write(resProxy.Body())
(pprof) 


</pre>

<pre>

(pprof) list pack.Encode   
Total: 840ms
ROUTINE ======================== github.com/Kingson4Wu/fast_proxy/outproxy/internal/pack.Encode in /Users/kingsonwu/Personal/go-src/f_proxy/outproxy/internal/pack/pack.go
         0       40ms (flat, cum)  4.76% of Total
         .          .    215:}
         .          .    216:
         .          .    217://----------
         .          .    218:
         .          .    219:func Encode(bodyBytes []byte, serviceName string) ([]byte, error) {
         .       20ms    220:   serviceConfig := outconfig.Get().GetServiceConfig(serviceName)
         .          .    221:   if serviceConfig == nil {
         .          .    222:           return nil, errors.New("get serviceConfig failure")
         .          .    223:   }
         .          .    224:
         .          .    225:   /** Data copy, to ensure that the current data remains unchanged */
         .          .    226:   sc := newSc(serviceConfig)
         .          .    227:   defer scPool.Put(sc)
         .          .    228:
         .          .    229:   middlewares := []Middleware{Encrypt, Compress, ProtobufEncode}
         .          .    230:
         .       20ms    231:   result, err := ApplyMiddlewares(bodyBytes, sc, middlewares...)
         .          .    232:   if err != nil {
         .          .    233:           return nil, err
         .          .    234:   }
         .          .    235:   return result, nil
         .          .    236:}


</pre>

<pre>

(pprof) list pack.Encode()
Total: 840ms
ROUTINE ======================== github.com/Kingson4Wu/fast_proxy/outproxy/internal/pack.Encode in /Users/kingsonwu/Personal/go-src/f_proxy/outproxy/internal/pack/pack.go
         0       40ms (flat, cum)  4.76% of Total
         .          .    215:}
         .          .    216:
         .          .    217://----------
         .          .    218:
         .          .    219:func Encode(bodyBytes []byte, serviceName string) ([]byte, error) {
         .       20ms    220:   serviceConfig := outconfig.Get().GetServiceConfig(serviceName)
         .          .    221:   if serviceConfig == nil {
         .          .    222:           return nil, errors.New("get serviceConfig failure")
         .          .    223:   }
         .          .    224:
         .          .    225:   /** Data copy, to ensure that the current data remains unchanged */
         .          .    226:   sc := newSc(serviceConfig)
         .          .    227:   defer scPool.Put(sc)
         .          .    228:
         .          .    229:   middlewares := []Middleware{Encrypt, Compress, ProtobufEncode}
         .          .    230:
         .       20ms    231:   result, err := ApplyMiddlewares(bodyBytes, sc, middlewares...)
         .          .    232:   if err != nil {
         .          .    233:           return nil, err
         .          .    234:   }
         .          .    235:   return result, nil
         .          .    236:}

</pre>

### memory

<pre>

 go test -v -run=^$ -bench=^BenchmarkRequestProxy$ -benchtime=2s -memprofile=mem.prof
goos: darwin
goarch: amd64
pkg: github.com/Kingson4Wu/fast_proxy/outproxy
cpu: Intel(R) Core(TM) i5-5257U CPU @ 2.70GHz
BenchmarkRequestProxy

BenchmarkRequestProxy-4             2833           1063248 ns/op           10744 B/op        242 allocs/op
PASS
ok      github.com/Kingson4Wu/fast_proxy/outproxy       3.834s

</pre>

<pre>
 go tool pprof -sample_index=alloc_space  outproxy.test mem.prof
File: outproxy.test
Type: alloc_space
Time: Mar 9, 2023 at 10:12am (CST)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top -cum
Showing nodes accounting for 0, 0% of 30.64MB total
Showing top 10 nodes out of 97
      flat  flat%   sum%        cum   cum%
         0     0%     0%    28.14MB 91.84%  github.com/Kingson4Wu/fast_proxy/outproxy.BenchmarkRequestProxy
         0     0%     0%    28.14MB 91.84%  github.com/Kingson4Wu/fast_proxy/outproxy.requestProxy
         0     0%     0%    28.14MB 91.84%  github.com/Kingson4Wu/fast_proxy/outproxy/internal/proxy.DoProxy
         0     0%     0%    28.14MB 91.84%  github.com/Kingson4Wu/fast_proxy/outproxy/internal/proxy.ForwardProxy
         0     0%     0%    28.14MB 91.84%  github.com/Kingson4Wu/fast_proxy/outproxy/internal/proxy.Func.Handle (inline)
         0     0%     0%    28.14MB 91.84%  testing.(*B).launch
         0     0%     0%    28.14MB 91.84%  testing.(*B).runN
         0     0%     0%    27.64MB 90.20%  github.com/Kingson4Wu/fast_proxy/outproxy/internal/proxy.fastDoProxy
         0     0%     0%    18.50MB 60.39%  github.com/spf13/viper.(*Viper).Get
         0     0%     0%    17.50MB 57.12%  github.com/Kingson4Wu/fast_proxy/outproxy/internal/pack.Encode
(pprof) inuse_space
(pprof) top -cum   
Showing nodes accounting for 1536.97kB, 75.01% of 2049.05kB total
Showing top 10 nodes out of 19
      flat  flat%   sum%        cum   cum%
 1024.41kB 49.99% 49.99%  1024.41kB 49.99%  runtime.malg
         0     0% 49.99%  1024.41kB 49.99%  runtime.newproc.func1
         0     0% 49.99%  1024.41kB 49.99%  runtime.newproc1
         0     0% 49.99%  1024.41kB 49.99%  runtime.systemstack
  512.56kB 25.01% 75.01%   512.56kB 25.01%  runtime.allocm
         0     0% 75.01%   512.56kB 25.01%  runtime.mstart
         0     0% 75.01%   512.56kB 25.01%  runtime.mstart0
         0     0% 75.01%   512.56kB 25.01%  runtime.mstart1
         0     0% 75.01%   512.56kB 25.01%  runtime.newm
         0     0% 75.01%   512.56kB 25.01%  runtime.resetspinning
(pprof) 


</pre>

<pre>
(pprof) alloc_space      
(pprof) list requestProxy
Total: 30.64MB
ROUTINE ======================== github.com/Kingson4Wu/fast_proxy/outproxy.requestProxy in /Users/kingsonwu/Personal/go-src/f_proxy/outproxy/outproxy.go
         0    28.14MB (flat, cum) 91.84% of Total
         .          .     19:
         .          .     20:   proxyType := proxy.FORWARD
         .          .     21:
         .          .     22:   p := proxy.GetProxy(proxyType)
         .          .     23:   if p != nil {
         .    28.14MB     24:           p.Handle(res, req)
         .          .     25:   }
         .          .     26:
         .          .     27:}
         .          .     28:
         .          .     29:func NewServer(c outconfig.Config, opts ...server.Option) {
(pprof) 

(pprof) list Handle      
Total: 30.64MB
ROUTINE ======================== github.com/Kingson4Wu/fast_proxy/outproxy/internal/proxy.Func.Handle in /Users/kingsonwu/Personal/go-src/f_proxy/outproxy/internal/proxy/proxy.go
         0    28.14MB (flat, cum) 91.84% of Total
         .          .     35:   Handle(http.ResponseWriter, *http.Request)
         .          .     36:}
         .          .     37:type Func func(http.ResponseWriter, *http.Request)
         .          .     38:
         .          .     39:func (f Func) Handle(w http.ResponseWriter, r *http.Request) {
         .    28.14MB     40:   f(w, r)
         .          .     41:}
         .          .     42:
         .          .     43:func ReverseProxy(w http.ResponseWriter, r *http.Request) {
         .          .     44:   Delegate(w, r)
         .          .     45:}
(pprof) list f
Total: 30.64MB


</pre>

### concurrent test

<pre>

go test -bench=Parallel -blockprofile=block.prof
goos: darwin
goarch: amd64
pkg: github.com/Kingson4Wu/fast_proxy/outproxy
cpu: Intel(R) Core(TM) i5-5257U CPU @ 2.70GHz
BenchmarkRequestProxyParallel-4             3241            371766 ns/op
PASS
ok      github.com/Kingson4Wu/fast_proxy/outproxy       3.104s

 go tool pprof outproxy.test block.prof
File: outproxy.test
Type: delay
Time: Mar 9, 2023 at 10:26am (CST)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 4.23s, 99.53% of 4.25s total
Dropped 95 nodes (cum <= 0.02s)
Showing top 10 nodes out of 15
      flat  flat%   sum%        cum   cum%
     2.13s 50.05% 50.05%      2.13s 50.05%  runtime.chanrecv1
     2.10s 49.48% 99.53%      2.10s 49.48%  sync.(*WaitGroup).Wait
         0     0% 99.53%      2.10s 49.48%  github.com/Kingson4Wu/fast_proxy/outproxy.BenchmarkRequestProxyParallel
         0     0% 99.53%      2.12s 49.85%  main.main
         0     0% 99.53%      2.12s 49.85%  runtime.main
         0     0% 99.53%      2.12s 49.85%  testing.(*B).Run
         0     0% 99.53%      2.10s 49.48%  testing.(*B).RunParallel
         0     0% 99.53%      2.11s 49.68%  testing.(*B).doBench
         0     0% 99.53%      2.10s 49.39%  testing.(*B).launch
         0     0% 99.53%      2.11s 49.68%  testing.(*B).run
(pprof) list requestProxy
Total: 4.25s
ROUTINE ======================== github.com/Kingson4Wu/fast_proxy/outproxy.requestProxy in /Users/kingsonwu/Personal/go-src/f_proxy/outproxy/outproxy.go
         0   221.46us (flat, cum) 0.0052% of Total
         .          .     19:
         .          .     20:   proxyType := proxy.FORWARD
         .          .     21:
         .          .     22:   p := proxy.GetProxy(proxyType)
         .          .     23:   if p != nil {
         .   221.46us     24:           p.Handle(res, req)
         .          .     25:   }
         .          .     26:
         .          .     27:}
         .          .     28:
         .          .     29:func NewServer(c outconfig.Config, opts ...server.Option) {
(pprof) 


</pre>


+ cd inproxy

### cpu

<pre>
go test -v -run=^$ -bench=^BenchmarkRequestProxy$ -benchtime=2s -cpuprofile=cpu.prof

goos: darwin
goarch: amd64
pkg: github.com/Kingson4Wu/fast_proxy/inproxy
cpu: Intel(R) Core(TM) i5-5257U CPU @ 2.70GHz
BenchmarkRequestProxy
BenchmarkRequestProxy-4             7467            301794 ns/op            7945 B/op        148 allocs/op
PASS
ok      github.com/Kingson4Wu/fast_proxy/inproxy        6.307s

</pre>

### memory

<pre>
go test -v -run=^$ -bench=^BenchmarkRequestProxy$ -benchtime=2s -memprofile=mem.prof

goos: darwin
goarch: amd64
pkg: github.com/Kingson4Wu/fast_proxy/inproxy
cpu: Intel(R) Core(TM) i5-5257U CPU @ 2.70GHz
BenchmarkRequestProxy
BenchmarkRequestProxy-4             7167            302430 ns/op            7932 B/op        148 allocs/op
PASS
ok      github.com/Kingson4Wu/fast_proxy/inproxy        6.426s

</pre>

<pre>

go tool pprof -sample_index=alloc_space  inproxy.test mem.prof
File: inproxy.test
Type: alloc_space
Time: Mar 9, 2023 at 11:41am (CST)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top -cum
Showing nodes accounting for 3.50MB, 5.43% of 64.51MB total
Showing top 10 nodes out of 128
      flat  flat%   sum%        cum   cum%
         0     0%     0%    53.51MB 82.94%  github.com/Kingson4Wu/fast_proxy/inproxy.BenchmarkRequestProxy
         0     0%     0%    53.51MB 82.94%  github.com/Kingson4Wu/fast_proxy/inproxy.requestProxy (inline)
         0     0%     0%    53.51MB 82.94%  github.com/Kingson4Wu/fast_proxy/inproxy/internal/proxy.DoProxy
         0     0%     0%    53.51MB 82.94%  testing.(*B).launch
         0     0%     0%    53.51MB 82.94%  testing.(*B).runN
         0     0%     0%       25MB 38.75%  github.com/spf13/viper.(*Viper).Get
         0     0%     0%       24MB 37.20%  github.com/spf13/viper.(*Viper).find
         0     0%     0%       16MB 24.80%  github.com/spf13/viper.(*Viper).isPathShadowedInFlatMap
         0     0%     0%       15MB 23.25%  github.com/spf13/cast.ToStringMap (inline)
    3.50MB  5.43%  5.43%       15MB 23.25%  github.com/spf13/cast.ToStringMapE
(pprof) list requestProxy
Total: 64.51MB
ROUTINE ======================== github.com/Kingson4Wu/fast_proxy/inproxy.requestProxy in /Users/kingsonwu/Personal/go-src/f_proxy/inproxy/inproxy.go
         0    53.51MB (flat, cum) 82.94% of Total
         .          .     14:
         .          .     15://go:embed *
         .          .     16:var ConfigFs embed.FS
         .          .     17:
         .          .     18:func requestProxy(res http.ResponseWriter, req *http.Request) {
         .    53.51MB     19:   proxy.DoProxy(res, req)
         .          .     20:
         .          .     21:}
         .          .     22:
         .          .     23:func NewServer(c inconfig.Config, opts ...server.Option) {
         .          .     24:   inconfig.Read(c)
(pprof) list DoProxy
Total: 64.51MB
ROUTINE ======================== github.com/Kingson4Wu/fast_proxy/inproxy/internal/proxy.DoProxy in /Users/kingsonwu/Personal/go-src/f_proxy/inproxy/internal/proxy/httpclientProxy.go
         0    53.51MB (flat, cum) 82.94% of Total
         .          .     37:
         .          .     38:// DoProxy /** forward request */
         .          .     39:func DoProxy(w http.ResponseWriter, r *http.Request) {
         .          .     40:
         .          .     41:   //Parameter check
         .        2MB     42:   if !strings.HasPrefix(r.RequestURI, inconfig.Get().ServerContextPath()) {
         .          .     43:           w.WriteHeader(http.StatusBadRequest)
         .          .     44:           return
         .          .     45:   }
         .          .     46:
         .          .     47:   /** Authentication to determine whether the uri has authority, etc. */
         .          .     48:   //How to prevent forgery of service name (signature verification is considered to be)
         .        2MB     49:   clientServiceName := server.Center().ClientName(r)
         .     3.50MB     50:   requestPath := servicediscovery.RealRequestUri(r.RequestURI)
         .     5.50MB     51:   if !inconfig.Get().ContainsCallPrivilege(clientServiceName, requestPath) {
         .          .     52:           writeErrorMessage(w, http.StatusBadRequest, "client has no privilege")
         .          .     53:           return
         .          .     54:   }
         .          .     55:
         .          .     56:   if limiter.IsLimit(clientServiceName, requestPath) {
         .          .     57:           writeErrorMessage(w, http.StatusBadRequest, "client is limit")
         .          .     58:           return
         .          .     59:   }
         .          .     60:
         .        1MB     61:   bodyBytes, pb, error := pack.DecodeReq(r)
         .          .     62:   defer pack.PbPool.Put(pb)
         .          .     63:   if error != nil {
         .          .     64:           writeErrorMessage(w, error.Code, error.Msg)
         .          .     65:           return
         .          .     66:   }
         .          .     67:
         .     7.50MB     68:   callUrl, rHandler := servicediscovery.Forward(r)
         .          .     69:
         .          .     70:   if callUrl == "" {
         .          .     71:           writeErrorMessage(w, http.StatusServiceUnavailable, "call url is blank")
         .          .     72:           return
         .          .     73:   }
         .          .     74:
         .          .     75:   // Create a request for forwarding
         .     9.50MB     76:   reqProxy, err := http.NewRequest(r.Method, callUrl, bytes.NewReader(bodyBytes))
         .          .     77:   if err != nil {
         .          .     78:           writeErrorMessage(w, http.StatusServiceUnavailable, "request wrap error")
         .          .     79:           return
         .          .     80:   }
         .          .     81:   if rHandler != nil {
         .          .     82:           rHandler(reqProxy)
         .          .     83:   }
         .          .     84:
         .          .     85:   deadTime := int64(servicediscovery.GetRequestDeadTime(r))
         .          .     86:   if deadTime > 0 {
         .          .     87:           if deadTime <= time.Now().Unix() {
         .          .     88:                   writeErrorMessage(w, http.StatusGatewayTimeout, "already reach dead time")
         .          .     89:                   return
         .          .     90:           } else {
         .          .     91:                   ctx, cancel := context.WithDeadline(context.Background(), time.Unix(deadTime, 0))
         .          .     92:                   defer cancel()
         .          .     93:                   reqProxy = reqProxy.WithContext(ctx)
         .          .     94:           }
         .          .     95:   } else {
         .        3MB     96:           reqServiceName := server.Center().ClientName(r)
         .     5.50MB     97:           timeout := inconfig.Get().GetTimeoutConfigByName(reqServiceName, r.RequestURI)
         .          .     98:           if timeout > 0 {
         .          .     99:                   ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
         .          .    100:                   defer cancel()
         .          .    101:                   reqProxy = reqProxy.WithContext(ctx)
         .          .    102:           }
         .          .    103:   }
         .          .    104:
         .          .    105:   // Header of forwarding request
         .          .    106:   for k, v := range r.Header {
         .        4MB    107:           reqProxy.Header.Set(k, v[0])
         .          .    108:   }
         .          .    109:
         .          .    110:   // make a request
         .        9MB    111:   responseProxy, err := client.Do(reqProxy)
         .          .    112:   if responseProxy != nil {
         .          .    113:           defer func() {
         .          .    114:                   io.Copy(io.Discard, responseProxy.Body)
         .          .    115:                   responseProxy.Body.Close()
         .          .    116:           }()
         .          .    117:   }
         .          .    118:   if err != nil {
         .          .    119:           server.GetLogger().Error("Error forwarding request", zap.Any("req forward err", err))
         .          .    120:
         .          .    121:           if errors.Is(err, context.DeadlineExceeded) {
         .          .    122:                   writeErrorMessage(w, http.StatusGatewayTimeout, "call timeout")
         .          .    123:                   return
         .          .    124:           }
         .          .    125:           writeErrorMessage(w, http.StatusServiceUnavailable, "call error")
         .          .    126:           return
         .          .    127:   }
         .          .    128:
         .          .    129:   // Header of the forwarded response
         .          .    130:   for k, v := range responseProxy.Header {
         .          .    131:           if strings.EqualFold(k, "Content-Length") {
         .          .    132:                   continue
         .          .    133:           }
         .          .    134:           w.Header().Set(k, v[0])
         .          .    135:   }
         .          .    136:
         .          .    137:   body, error := pack.EncodeResp(responseProxy, pb)
         .          .    138:   if error != nil {
         .          .    139:           writeErrorMessage(w, error.Code, error.Msg)
         .          .    140:           return
         .          .    141:   }
         .          .    142:
         .   512.02kB    143:   resProxyBody := io.NopCloser(bytes.NewBuffer(body))
         .          .    144:   defer resProxyBody.Close()
         .          .    145:
         .          .    146:   // response status code
         .          .    147:   w.WriteHeader(responseProxy.StatusCode)
         .          .    148:   // Copy the forwarded response Body to the response Body
         .          .    149:   w.Header().Set("Content-Length", strconv.Itoa(len(body)))
         .          .    150:   io.Copy(w, resProxyBody)
         .          .    151:
         .   516.01kB    152:}
         .          .    153:
         .          .    154:func writeErrorMessage(res http.ResponseWriter, statusCode int, errorHeader string) {
         .          .    155:   res.Header().Add("proxy_error_message", errorHeader)
         .          .    156:   res.Header().Add("proxy_name", "in_proxy")
         .          .    157:   res.WriteHeader(statusCode)
ROUTINE ======================== github.com/Kingson4Wu/fast_proxy/inproxy/internal/proxy.DoProxy.func1 in /Users/kingsonwu/Personal/go-src/f_proxy/inproxy/internal/proxy/httpclientProxy.go
         0   516.01kB (flat, cum)  0.78% of Total
         .          .    109:
         .          .    110:   // make a request
         .          .    111:   responseProxy, err := client.Do(reqProxy)
         .          .    112:   if responseProxy != nil {
         .          .    113:           defer func() {
         .   516.01kB    114:                   io.Copy(io.Discard, responseProxy.Body)
         .          .    115:                   responseProxy.Body.Close()
         .          .    116:           }()
         .          .    117:   }
         .          .    118:   if err != nil {
         .          .    119:           server.GetLogger().Error("Error forwarding request", zap.Any("req forward err", err))
(pprof) 


</pre>


### concurrent test

<pre>

 go test -bench=Parallel -blockprofile=block.prof

goos: darwin
goarch: amd64
pkg: github.com/Kingson4Wu/fast_proxy/inproxy
cpu: Intel(R) Core(TM) i5-5257U CPU @ 2.70GHz
BenchmarkRequestProxyParallel-4             7880            127795 ns/op
PASS
ok      github.com/Kingson4Wu/fast_proxy/inproxy        5.591s


 go tool pprof inproxy.test block.prof
File: inproxy.test
Type: delay
Time: Mar 9, 2023 at 11:52am (CST)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 8681.30ms, 99.89% of 8690.80ms total
Dropped 43 nodes (cum <= 43.45ms)
Showing top 10 nodes out of 36
      flat  flat%   sum%        cum   cum%
 6480.28ms 74.56% 74.56%  6480.28ms 74.56%  runtime.selectgo
 1178.35ms 13.56% 88.12%  1178.35ms 13.56%  runtime.chanrecv1
 1022.67ms 11.77% 99.89%  1022.67ms 11.77%  sync.(*WaitGroup).Wait
         0     0% 99.89%   145.68ms  1.68%  bytes.(*Buffer).ReadFrom
         0     0% 99.89%  1022.67ms 11.77%  github.com/Kingson4Wu/fast_proxy/inproxy.BenchmarkRequestProxyParallel
         0     0% 99.89%  2748.88ms 31.63%  github.com/Kingson4Wu/fast_proxy/inproxy.BenchmarkRequestProxyParallel.func1
         0     0% 99.89%  2752.32ms 31.67%  github.com/Kingson4Wu/fast_proxy/inproxy.requestProxy (inline)
         0     0% 99.89%   145.68ms  1.68%  github.com/Kingson4Wu/fast_proxy/inproxy/internal/pack.EncodeResp
         0     0% 99.89%  2752.32ms 31.67%  github.com/Kingson4Wu/fast_proxy/inproxy/internal/proxy.DoProxy
         0     0% 99.89%   145.69ms  1.68%  io.Copy (inline)
(pprof) 


</pre>