local token = "4a363afc-def8-4451-8905-2aa8525f912d"
local bot_id = "2c23fd65-c6e6-40f4-acca-c59ed86a4620"

local counter = 0

wrk.method = "POST"
wrk.headers["Authorization"] = "Bearer " .. token
wrk.headers["Content-Type"] = "application/json"

request = function()
   counter = counter + 1
   local body = string.format([[{
      "bot_id": "%s",
      "status": "success",
      "msg": "Load test message #%d"
   }]], bot_id, counter)
   
   return wrk.format("POST", nil, wrk.headers, body)
end

response = function(status, headers, body)
   if status ~= 201 then
      print("Error: " .. status .. " - " .. body)
   end
end

