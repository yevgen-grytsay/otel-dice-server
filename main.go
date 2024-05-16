package main

import (
	"context"
	"log"
	"net/http"
	"os"
	monitoring "yevgen-grytsay/dice/otel"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var appVersion = "unknown"
var appEnv = os.Getenv("OTEL_DICE_ENV")

func main() {
	log.Printf("Version: %s", appVersion)

	ctx := context.Background()
	shutdown, err := monitoring.SetupOTelSDK(ctx, monitoring.ParseEnv(appEnv))
	if err != nil {
		panic(err)
	}
	defer shutdown(ctx)

	otelHandler := otelhttp.NewHandler(http.HandlerFunc(rolldice), "Hello")

	http.Handle("/rolldice", otelHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
