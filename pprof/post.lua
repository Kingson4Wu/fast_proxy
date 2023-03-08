wrk.method = "POST"
wrk.body = 'param=hello'
--wrk.headers["protocol"] = "json"
wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"
wrk.headers["C_ServiceName"] = "chat_service"
wrk.headers["Accept"]= "*/*"
response = function(status, headers, body)
--print(body)
end

-- curl -i -H 'C_ServiceName:chat_service'  -d "" "http://127.0.0.1:8034/search_service/api/service"