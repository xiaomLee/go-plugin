package ratelimit

import (
	"io/ioutil"
	redis_pool "github.com/xiaomLee/go-plugin/redis"
	"strconv"
)

type Rule struct {
	Duration int `json:"duration"` // 一段时间，单位秒
	Limit    int `json:"limit"`    // 时间段内同一个会话允许访问的次数
}

var scriptSha string

func InitRateLimit(name, addr, port, pwd string, dbNum int) error {
	// init redis
	if err := redis_pool.AddRedisInstance(name, addr, port, pwd, dbNum); err != nil {
		return err
	}
	println("redis init success")

	redis := redis_pool.MustGetRedisInstance()

	script, err := ioutil.ReadFile("redis_ratelimit.lua")
	if err != nil {
		return err
	}

	scriptSha, err = redis.ScriptLoad(string(script)).Result()
	println(scriptSha)

	return err
}

func TokenAccess(key string, rules []*Rule) bool {
	if len(rules) == 0 {
		return true
	}

	tmp := make([]interface{}, 0)
	for _, rule := range rules {
		tmp = append(tmp, rule.Duration)
		tmp = append(tmp, rule.Limit)
	}
	return tokenAccess(key, tmp...)
}

// tokenAccess - 令牌验证具体实现
//
// PARAMS:
// - key: 需要验证令牌的key值
// - rules: duration和limie交替组成，每一组duration和limit代表一组规则
//
// RETURNS:
// bool 是否通过了令牌限制
func tokenAccess(key string, rules ...interface{}) bool {
	redis := redis_pool.MustGetRedisInstance()

	keys := []string{key, strconv.Itoa(len(rules))}

	val, err := redis.EvalSha(scriptSha, keys, rules...).Int()
	if err != nil {
		return false
	}
	if val == 1 {
		return true
	}
	return false
}
