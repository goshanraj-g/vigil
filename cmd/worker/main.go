package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/goshanraj-g/vigil/internal/ai"
	"github.com/goshanraj-g/vigil/internal/db"
	"github.com/goshanraj-g/vigil/internal/dedup"
	"github.com/goshanraj-g/vigil/internal/search"
	"github.com/goshanraj-g/vigil/internal/worker"
	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	ctx := context.Background()

	pool, err := db.Connect(ctx)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer pool.Close()

	sc, err := search.New()
	if err != nil {
		log.Fatalf("search: %v", err)
	}

	ac, err := ai.New()
	if err != nil {
		log.Fatalf("ai: %v", err)
	}

	dc, err := dedup.New()
	if err != nil {
		log.Fatalf("dedup: %v", err)
	}

	h := worker.New(pool, sc, ac, dc)

	redisURL := os.Getenv("REDIS_URL")

	// creatre asynq worker server
	srv := asynq.NewServer(asynq.RedisClientOpt{Addr: redisURL}, asynq.Config{})

	// create multiplexer to route incoming jobs to right handler
	mux := asynq.NewServeMux()

	// register the route
	mux.HandleFunc(worker.TypeSearchMonitor, h.Handle)

	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisURL})
	defer client.Close()
	sched := worker.NewScheduler(pool, client, time.Minute)

	go sched.Run(ctx)
	if err := srv.Run(mux); err != nil {
		log.Fatalf("worker: %v", err)
	}
}
