local b64 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

function b64decode(data)
  data = data:gsub("[-]", "+"):gsub("[_]", "/")
  local pad = #data % 4
  if pad > 0 then data = data .. string.rep("=", 4 - pad) end
  local result = ""
  for i = 1, #data, 4 do
    local a, b, c, d = data:byte(i, i+3)
    a = b64:find(string.char(a)) - 1
    b = b64:find(string.char(b)) - 1
    c = c and b64:find(string.char(c)) or 64
    d = d and b64:find(string.char(d)) or 64
    if c then c = c - 1 else c = 0 end
    if d then d = d - 1 else d = 0 end
    local n = a * 262144 + b * 4096 + c * 64 + d
    result = result .. string.char(math.floor(n / 65536) % 256)
    if c ~= 63 then result = result .. string.char(math.floor(n / 256) % 256) end
    if d ~= 63 then result = result .. string.char(n % 256) end
  end
  return result
end

function extract_tenant(token)
  local parts = {}
  for part in token:gmatch("[^.]+") do
    table.insert(parts, part)
  end
  if #parts ~= 3 then return nil end
  local payload = b64decode(parts[2])
  local tenant = payload:match('"tenant"%s*:%s*"([^"]+)"')
  return tenant
end

function get_subdomain(host)
  local h = host:match("^([^:]+)")
  if not h then h = host end
  local sub = h:match("^([^.]+)%.")
  return sub
end

function envoy_on_request(request_handle)
  local path = request_handle:headers():get(":path")

  if path and path:find("^/auth/") then
    return
  end

  local auth = request_handle:headers():get("authorization")
  if not auth then
    request_handle:respond(
      {[":status"] = "401", ["content-type"] = "application/json"},
      '{"error":"missing authorization header"}'
    )
    return
  end

  local token = auth:match("^Bearer%s+(.+)$")
  if not token then
    request_handle:respond(
      {[":status"] = "401", ["content-type"] = "application/json"},
      '{"error":"invalid authorization format"}'
    )
    return
  end

  local tenant = extract_tenant(token)
  if not tenant then
    request_handle:respond(
      {[":status"] = "401", ["content-type"] = "application/json"},
      '{"error":"cannot extract tenant from token"}'
    )
    return
  end

  local host = request_handle:headers():get(":authority")
  local subdomain = get_subdomain(host)

  if not subdomain or subdomain ~= tenant then
    request_handle:respond(
      {[":status"] = "403", ["content-type"] = "application/json"},
      '{"error":"token not valid for this subdomain"}'
    )
    return
  end
end
