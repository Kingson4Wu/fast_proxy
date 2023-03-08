
<pre>
8 GB 1867 MHz DDR3
2.7 GHz Dual-Core Intel Core i5
</pre>

+ mac : `sysctl -n hw.logicalcpu` : 4

----

## v1
+ ea1cfb81b0cf24996ebfc73d89d236a0dd212652

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

+ 初步看来： 4 threads and 128 connections 测试效果最佳
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

</pre>

