local token = "aaaaaaaa-0000-1111-2222-333333333333"


wrk.method = "GET"
wrk.headers["Authorization"] = "Bearer " .. token
wrk.headers["Content-Type"] = "application/json"

request = function()
   return wrk.format(nil, nil, wrk.headers)
end

