package workers

import (
	"log"

	"github.com/hibiken/asynq"

	"github.com/rajiknows/webhook-mock/internal/utils"
	"github.com/rajiknows/webhook-mock/pkg/queue"
)

func StartWorker() {
	mux := asynq.NewServeMux()
	mux.HandleFunc(queue.DeliveryTaskType, utils.ProcessDeliveryTask)
	if err := queue.Server.Run(mux); err != nil {
		log.Fatalf("Failed to start worker: %v", err)
	}
}
