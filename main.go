package main

import (
	"log"
	"net/http"

	"github.com/scottd018/rosa-windows-overcommit-webhook/webhook"
)

func main() {
	// create the webhook
	w, err := webhook.NewWebhook()
	if err != nil {
		log.Fatalf("failed to create webhook: %v", err)
	}

	// set the handlers and start
	http.HandleFunc("/validate", w.Validate)
	http.HandleFunc("/healthz", w.HealthZ)
	log.Println("Starting webhook server on :8443")
	log.Fatal(http.ListenAndServeTLS(":8443", "/ssl_certs/tls.crt", "/ssl_certs/tls.key", nil))
	//log.Fatal(http.ListenAndServe(":8080", nil))
}
