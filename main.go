package main

import (
	"context"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()
	shutdown, err := setupOTelSDK(ctx)
	if err != nil {
		panic(err)
	}
	defer shutdown(ctx)

	http.HandleFunc("/rolldice", rolldice)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
