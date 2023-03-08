wrk.method = "POST"
wrk.body = 'param=hello'
--wrk.headers["protocol"] = "json"
wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"
wrk.headers["C_ServiceName"] = "chat_service"
wrk.headers["Accept"]= "*/*"
function response(status, headers, body)
end
--[[ function response(status, headers, body)
  if status == 200 then
    print("Response: " .. body)
  else
    print("Error: " .. status)
    for key, value in pairs(headers) do
          print(key .. ": " .. value)
    end
  end
end ]]

-- curl -i -H 'C_ServiceName:chat_service'  -d "" "http://127.0.0.1:8034/search_service/api/service"