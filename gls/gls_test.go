package gls

import (
	"fmt"
	"testing"
)

func test1() bool {
	ctx := GetContext()
	if ctx.Id() != 1 {
		return false
	}
	v, ok := ctx.Get("hello")
	if !ok || v.(string) != "world" {
		return false
	}
	return true
}

func TestGoroutineId(t *testing.T) {
	ctx := GetContext()
	if ctx.Id() != 1 {
		t.Errorf("first ctx id != 1")
	}
	ctx.Put("hello", "world")
	ok := test1()
	if !ok {
		t.Errorf("context failed")
	}

	//wg := sync.WaitGroup{}
	//wg.Add(1)
	//Go(func() {
	//	defer wg.Done()
	//	ctx := GetContext()
	//	if ctx == nil {
	//		t.Fatal("ctx is nil")
	//	}
	//	fmt.Printf("gorutine1: %+v \n", ctx)
	//})

	WithContext(ctx, func() {
		ctx := GetContext()
		fmt.Printf("gorutine with context :%+v \n", ctx)
	})

	WithNewContext(func() {
		ctx := GetContext()
		fmt.Printf("gorutine with new context :%+v \n", ctx)
	})
	//wg.Wait()
}

func TestStart(t *testing.T) {
	start(256, func() {
		id := getContextId()
		if id != 256 {
			t.Fatal("ctx id wrong!")
		}
	})
	start(65536, func() {
		id := getContextId()
		if id != 65536 {
			t.Fatal("ctx id wrong!")
		}
	})
	start(0xabcdef, func() {
		id := getContextId()
		if id != 0xabcdef {
			t.Fatal("ctx id wrong!")
		}
	})
}
