package logger

import (
	"context"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func Init() {
	ConfigLogger(Config{
		AppName:       "test.logger",
		Level:         5,
		Formatter:     "text",
		Dir:           "./",
		Filename:      "default.log",
		MaxAge:        -1,
		RotationSize:  0,
		RotationTime:  0,
		RotationCount: 5,
	},
	)
}

func TestLogger(t *testing.T) {
	Init()
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			Error("test_err", logrus.Fields{"some_info": i})
		} else if i%3 == 0 {
			Debug("test_debug", logrus.Fields{"some_info": i})
		} else if i%4 == 0 {
			Warning("test_warning", logrus.Fields{"some_info": i})
		} else {
			Info("test_info", logrus.Fields{"some_info": i})
		}
	}


	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			Error("test_err", logrus.Fields{"some_info": i})
		} else if i%3 == 0 {
			Debug("test_debug", logrus.Fields{"some_info": i})
		} else if i%4 == 0 {
			Warning("test_warning", logrus.Fields{"some_info": i})
		} else {
			Info("test_info", logrus.Fields{"some_info": i})
		}
		time.Sleep(time.Second*1)
	}

	time.Sleep(time.Second*10)
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			Error("test_err", logrus.Fields{"some_info": i})
		} else if i%3 == 0 {
			Debug("test_debug", logrus.Fields{"some_info": i})
		} else if i%4 == 0 {
			Warning("test_warning", logrus.Fields{"some_info": i})
		} else {
			Info("test_info", logrus.Fields{"some_info": i})
		}
	}
}

func TestErrorWithContext(t *testing.T) {
	Init()
	ctx := context.WithValue(context.Background(), "requestid", "122222")
	ErrorWithContext(ctx, "test err withcontext", nil)
}
