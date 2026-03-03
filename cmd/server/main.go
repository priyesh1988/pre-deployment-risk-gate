package main

import (
	"log"
	"net/http"

	"github.com/yourname/guardrail-saas/internal/githubapp"
)

func main() {
	http.HandleFunc("/webhook", githubapp.HandleWebhook)
	log.Println("Guardrail SaaS running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
