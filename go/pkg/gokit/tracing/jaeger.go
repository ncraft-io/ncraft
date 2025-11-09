package tracing

import (
    "github.com/ncraft-io/ncraft/go/pkg/ncraft/logs"
    "github.com/opentracing/opentracing-go"
    "github.com/uber/jaeger-client-go"
    jaegercfg "github.com/uber/jaeger-client-go/config"
    "github.com/uber/jaeger-client-go/rpcmetrics"
    "github.com/uber/jaeger-lib/metrics"
    "io"
    "time"
)

// Init creates a new instance of Jaeger tracer.
func newTracer(serviceName string, metricsFactory metrics.Factory, param float64, backendHostPort string) (opentracing.Tracer, io.Closer, error) {
    var err error
    cfg := jaegercfg.Configuration{
        ServiceName: serviceName,
        Sampler: &jaegercfg.SamplerConfig{
            Type:  jaeger.SamplerTypeRateLimiting,
            Param: param,
        },
    }

    var sender jaeger.Transport
    if sender, err = jaeger.NewUDPTransport(backendHostPort, 0); err != nil {
        logs.Error("cannot initialize UDP sender", err)
        return nil, nil, err
    }

    return cfg.NewTracer(
        jaegercfg.Reporter(jaeger.NewRemoteReporter(
            sender,
            jaeger.ReporterOptions.BufferFlushInterval(1*time.Second),
        )),
        jaegercfg.Metrics(metricsFactory),
        jaegercfg.Observer(rpcmetrics.NewObserver(metricsFactory, rpcmetrics.DefaultNameNormalizer)),
    )
}
