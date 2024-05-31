package main

import (
	"fmt"
	"go.gllm.dev/vanity/handlers/gohdl"
	"go.gllm.dev/vanity/services/gosvc"
	"log"
	"net/http"
	"os"
)

func main() {
	domain := os.Getenv("DOMAIN")
	if domain == "" {
		log.Fatal("DOMAIN is required")
	}

	repository := os.Getenv("REPOSITORY")
	if repository == "" {
		log.Fatal("REPOSITORY is required")
	}

	goHdl := gohdl.New(gosvc.New(domain, repository))

	mux := http.NewServeMux()

	mux.HandleFunc("/", goHdl.Handle)

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
