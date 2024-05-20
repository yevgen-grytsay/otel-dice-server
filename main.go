package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	monitoring "yevgen-grytsay/dice/otel"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var appVersion = "unknown"
var appEnv = os.Getenv("OTEL_DICE_ENV")

var livenessHandler = func(w http.ResponseWriter, r *http.Request) {
	if _, err := io.WriteString(w, "OK"); err != nil {
		log.Printf("Write failed: %v\n", err)
	}
}

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

	http.HandleFunc("/liveness", livenessHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
