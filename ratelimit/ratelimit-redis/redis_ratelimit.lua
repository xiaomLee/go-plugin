-- 脚本说明
-- EVAL script numkeys key [key ...] arg [arg ...]
-- 示例 对limitKey 实行10s一次 的限流
-- EVAL script 2 limitKey 2 10 1

-- 由于在脚本中调用了redis.call("TIME"), 在持久化及主从复制时会导致写入结果存在不确定性。
-- redis防止随机写入 https://yq.aliyun.com/articles/195914
-- TODO 该方法在大流量写入时不建议用 时间参数改为由外部传入
redis.replicate_commands()

-- 获取当前时间戳 微秒
local timestamp = redis.call("TIME")
local time_now = timestamp[1] * 1e6 + timestamp[2]
local N = tonumber(KEYS[2])
local keys = {}
local remains = {}
local lastTimes = {}
local durations = {}
local limits = {}

-- 将所有规则整理出来，包括key，duration和limit
local j = 1
for i=1, N-1, 2 do
	durations[j] = tonumber(ARGV[i])
	limits[j] = tonumber(ARGV[i+1])
	keys[j] = KEYS[1]..durations[j]..limits[j]
	j = j+1
end

j = j-1

-- 遍历每一条规则，判断是否还有token剩余，是否满足重新补充的条件
for i=1, j, 1 do
	local ratelimit_info=redis.pcall("HMGET",keys[i],"remain_token","last_fill_time")
	remains[i] = tonumber(ratelimit_info[1])
	lastTimes[i] = tonumber(ratelimit_info[2])

	-- 之前不存在，创建，并设置过期为一小时
	if (lastTimes[i]==nil) then
    	redis.call("HMSET",keys[i],"remain_token",limits[i],"last_fill_time",time_now)
    	redis.call("EXPIRE", keys[i], 3600)
    	lastTimes[i] = time_now
    	remains[i] = limits[i]
	end

	-- 剩余token不足，判断是否需要补充
	if (remains[i] == 0) then
		if (time_now>lastTimes[i]) then
			local a,b
			a = math.floor((durations[i] * 1e6)/limits[i])
			remains[i],b = math.modf((time_now - lastTimes[i])/a)
			if (remains[i]>limits[i]) then
				remains[i] = limits[i]
			end
		    lastTimes[i] = time_now
		end
	end

	-- 任意一条规则不满足，则返回失败
	if (remains[i] == 0) then
		return 0
	end
end

-- 对每一条规则减去一个token，并返回成功
for i=1, j, 1 do
	redis.pcall("HMSET", keys[i],"remain_token",remains[i]-1,"last_fill_time",lastTimes[i])
end

return 1