```sh
go get go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp

go get go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp
```

```sh
kubectl run dice --image=yevhenhrytsai/dice:v1.0.0
```

## Collector
Лог TracesExporter з'являється одразу піся запиту до `/rolldice`.

Лог MetricsExporter з'являється з інтервалом, вказаним в налаштуваннях `metric.MeterProvider`.
```go
	meterProvider := metric.NewMeterProvider(metric.WithReader(metric.NewPeriodicReader(exp, metric.WithInterval(10*time.Second))))
```
```
2024-05-15T16:08:45.059Z        info    TracesExporter  {"kind": "exporter", "data_type": "traces", "name": "debug", "resource spans": 1, "spans": 1}
2024-05-15T16:08:46.062Z        info    TracesExporter  {"kind": "exporter", "data_type": "traces", "name": "debug", "resource spans": 1, "spans": 1}
2024-05-15T16:08:47.923Z        info    MetricsExporter {"kind": "exporter", "data_type": "metrics", "name": "debug", "resource metrics": 1, "metrics": 1, "data points": 1}
2024-05-15T16:08:50.075Z        info    TracesExporter  {"kind": "exporter", "data_type": "traces", "name": "debug", "resource spans": 1, "spans": 3}
2024-05-15T16:08:57.952Z        info    MetricsExporter {"kind": "exporter", "data_type": "metrics", "name": "debug", "resource metrics": 1, "metrics": 1, "data points": 3}
2024-05-15T16:09:07.983Z        info    MetricsExporter {"kind": "exporter", "data_type": "metrics", "name": "debug", "resource metrics": 1, "metrics": 1, "data points": 3}
2024-05-15T16:09:17.975Z        info    TracesExporter  {"kind": "exporter", "data_type": "traces", "name": "debug", "resource spans": 1, "spans": 1}
2024-05-15T16:09:18.016Z        info    MetricsExporter {"kind": "exporter", "data_type": "metrics", "name": "debug", "resource metrics": 1, "metrics": 1, "data points": 3}
2024-05-15T16:09:18.977Z        info    TracesExporter  {"kind": "exporter", "data_type": "traces", "name": "debug", "resource spans": 1, "spans": 1}
2024-05-15T16:09:20.984Z        info    TracesExporter  {"kind": "exporter", "data_type": "traces", "name": "debug", "resource spans": 1, "spans": 1}
```


## Resources
- [otel/exporters/otlp/otlpmetric/otlpmetrichttp](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp)
