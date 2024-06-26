package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

const redisAddr = "redis://:QAxtXFLVXIW3bXf9@172.25.1.167:6379/11"

func main() {
	redisOpt, _ := asynq.ParseRedisURI(redisAddr)
	srv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"sdxj:image:analyze": 6,
				"default":            3,
				"low":                1,
			},
			// See the godoc for other configuration options
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.Handle("sdxj:image:analyze", loggingMiddleware(&MyTaskHandler{}))

	// ...register other handlers...

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

type MyTaskHandler struct {
	// ... fields
}

func (h *MyTaskHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	fmt.Println(t.Type())

	return nil
}

func loggingMiddleware(h asynq.Handler) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
		start := time.Now()
		log.Printf("Start processing %q", t.Type())
		err := h.ProcessTask(ctx, t)
		if err != nil {
			return err
		}
		log.Printf("Finished processing %q: Elapsed Time = %v", t.Type(), time.Since(start))
		return nil
	})
}
