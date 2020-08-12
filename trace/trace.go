package trace

import (
	"io"

	"github.com/opentracing/opentracing-go"
)

func NewJaegerTracer(serviceName string) (opentracing.Tracer, io.Closer) {

}
