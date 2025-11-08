package tracing

import (
    "github.com/ncraft-io/ncraft/go/pkg/ncraft/logs"
    "github.com/opentracing/opentracing-go"
    "github.com/uber/jaeger-lib/metrics"
    "io"
)

func New(name string) (opentracing.Tracer, io.Closer) {
    return NewWith(name, NewConfig("tracing"))
}

func NewWith(name string, cfg *Config) (opentracing.Tracer, io.Closer) {
    if cfg == nil {
        logs.Fatal("failed to create the tracer coz of given nil config")
        return nil, nil
    }

    if cfg.Enable {
        tracer, c, err := newTracer(name, metrics.NullFactory, cfg.Param, cfg.Url)
        if err != nil {
            logs.Fatalw("failed to create the tracing")
            return nil, nil
        }

        if tracer == nil {
            logs.Fatal()
        }

        opentracing.SetGlobalTracer(tracer)
        return tracer, c
    }

    return opentracing.NoopTracer{}, nil
    // return nil, nil
}
