package main

import (
	"context"
	"log"
	"net/http"
	"yevgen-grytsay/dice/otel"
)

var appVersion = "unknown"

func main() {
	log.Printf("Version: %s", appVersion)

	ctx := context.Background()
	shutdown, err := otel.SetupOTelSDK(ctx)
	if err != nil {
		panic(err)
	}
	defer shutdown(ctx)

	http.HandleFunc("/rolldice", rolldice)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
