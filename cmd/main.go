package main

import (
	"log"
	"os"

	"github.com/gocardless-ingest/internal/clients"
	"github.com/gocardless-ingest/internal/scheduler"
)

func main() {
	log.Println("Inicio del programa")

	secretID := os.Getenv("GC_CLIENT_ID")
	secretKey := os.Getenv("GC_SECRET_KEY")
	gocardless_client := clients.NewGoCardlessClient(secretID, secretKey)

	scheduler.StartCron(gocardless_client)

}
