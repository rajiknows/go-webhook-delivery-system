package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rajiknows/webhook-mock/internal/api"
	"github.com/rajiknows/webhook-mock/internal/services"
	"github.com/rajiknows/webhook-mock/internal/workers"
	"github.com/rajiknows/webhook-mock/pkg/db"
	"github.com/rajiknows/webhook-mock/pkg/queue"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or failed to load")
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please specify 'api' or 'worker'")
	}
	command := os.Args[1]

	db.ConnectDatabase()
	services.InitRedis()
	queue.InitQueue()

	switch command {
	case "api":
		r := api.SetupRouter()
		if err := r.Run(":8080"); err != nil {
			log.Fatal("Failed to start API server:", err)
		}
	case "worker":
		workers.StartWorker()
	default:
		log.Fatal("Unknown command:", command)
	}
}
