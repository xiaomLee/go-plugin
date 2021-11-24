package ratelimit

import (
	"fmt"
	"golang.org/x/time/rate"
	"strconv"
	"testing"
	"time"
)

func TestXXX(t *testing.T) {
	fmt.Println(strconv.FormatInt(time.Now().UnixNano(), 10))
}

func TestInit(t *testing.T) {
	err := InitRateLimit("", "127.0.0.1", "6379", "", 0)
	if err != nil {
		t.Error(err)
	}
}

func TestTokenAccess(t *testing.T) {
	err := InitRateLimit("", "127.0.0.1", "6379", "", 0)
	if err != nil {
		t.Error(err)
	}

	rules := make([]*Rule, 0)

	rule1 := &Rule{
		Duration: 5,
		Limit:    1,
	}
	rules = append(rules, rule1)

	// 添加规则2之后 单例测试将失败
	//rule2 := &Rule{
	//	Duration: 10,
	//	Limit: 1,
	//}
	//rules = append(rules, rule2)

	ticker := time.NewTicker(1 * time.Second)
	exit := time.After(20 * time.Second)

	i := 0
	for {
		select {
		case <-ticker.C:
			result := TokenAccess("urlpath/user/info", rules)
			fmt.Printf("index:%d result:%v \n", i, result)

			if i%5 == 0 && !result {
				t.Errorf("index:%d result:%v", i, result)
			}
			if i%5 != 0 && result {
				t.Errorf("index:%d result:%v", i, result)
			}
			i++
		case <-exit:
			return
		}
	}

}

// 测试golang标准库的限流器
func TestTimeRate(t *testing.T) {
	limiter := rate.NewLimiter(1, 10)
	limiter.Allow()
}
