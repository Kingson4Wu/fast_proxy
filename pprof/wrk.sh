/usr/local/Cellar/wrk/4.1.0/bin/wrk -d 60 -c 100  -t 32 -s post.lua http://127.0.0.1:8034/search_service/api/service

# cd pprof
# /usr/local/Cellar/wrk/4.1.0/bin/wrk -d 60 -t2 -c10 -d10s -s post.lua 'http://127.0.0.1:8034/search_service/api/service'