package redis

import (
	"context"
	"testing"
	"time"
)

func TestAddRedisInstance(t *testing.T) {
	if err := AddRedisInstance(
		"default",
		"127.0.0.1",
		"6379",
		"",
		0,
		); err!=nil {
		t.Fatal(err)
	}
}

func TestGetRedisInstance(t *testing.T) {
	TestAddRedisInstance(t)
	redis, ok:= GetRedisInstance("default")
	if !ok {
		t.Fatal("get redis instance err")
	}
	_, err:= redis.Set(context.Background(), "test:redis:key1", "this is key1 value", 3600*time.Second).Result()
	if err!=nil {
		t.Fatal(err)
	}
	v, err:= redis.Get(context.Background(), "test:redis:key1").Result()
	if err!=nil {
		t.Fatal(err)
	}
	if v != "this is key1 value" {
		t.Error("result not expected")
	}
	t.Log(v)
}