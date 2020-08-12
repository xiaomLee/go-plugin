package jaeger

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"testing"
	"time"
)

func TestInitJaeger(t *testing.T) {
	tracer, closer := InitJaeger("hello-world")
	defer closer.Close()


	opentracing.InitGlobalTracer(tracer)

	path := "hello"
	traceId := "22344556"

	begin := time.Now()
	span := tracer.StartSpan(path, opentracing.Tag{
		Key:   "traceId",
		Value: traceId,
	})
	span.LogFields(
		log.String("tag", "request_out"),
		log.String("params", "a=b&c=d"),
		)

	// do something
	time.Sleep(time.Duration(2) * time.Millisecond)

	// start an other span
	span2 := tracer.StartSpan(path, opentracing.ChildOf(span.Context()),
		opentracing.Tag{
			Key:   "traceId",
			Value: traceId,
		},
		)

	time.Sleep(1 *time.Millisecond)
	span2.LogFields(
		log.String("tag", "span2 do something"),
		)
	span2.Finish()

	span.LogFields(
		log.String("tag", "request_out"),
		log.String("response", "hello world"),
		log.Int64("latency", time.Since(begin).Microseconds()),
		)
	span.Finish()
}
