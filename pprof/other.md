+ top -pid 5729
+ `top -pid $(ps -ef | grep inproxy | grep -v 'grep' | awk '{print $2}')`
+ `top -pid $(ps -ef | grep outproxy | grep -v 'grep' | awk '{print $2}')`

### TimeWait
+ `netstat -n | awk '/^tcp/ {++S[$NF]} END {for(a in S) print a, S[a]}'`
+ `netstat -n |grep TIME_WAIT|grep '127.0.0.1'|awk '{print $5}'|sort | uniq -c | sort -k1,1nr | head -20`


### Client
+ curl -i -H 'C_ServiceName:chat_service'  -d "" "http://127.0.0.1:8034/search_service/api/service"

