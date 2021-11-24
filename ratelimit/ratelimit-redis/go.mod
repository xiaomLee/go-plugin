module github.com/xiaomLee/go-plugin/ratelimit-redis/ratelimit

go 1.16

require (
	github.com/xiaomLee/go-plugin/redis v0.0.0
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
)

replace github.com/xiaomLee/go-plugin/redis => ../../redis
