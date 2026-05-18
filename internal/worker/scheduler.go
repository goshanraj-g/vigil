package worker

import (
	"context"
	"log"
	"time"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Scheduler struct {
	pool     *pgxpool.Pool
	client   *asynq.Client
	interval time.Duration
}

func NewScheduler(pool *pgxpool.Pool, client *asynq.Client, interval time.Duration) *Scheduler {
	return &Scheduler{pool: pool, client: client, interval: interval}
}

func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop() // stops when this function returns

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.tick(ctx)
		}
	}
}

func (s *Scheduler) tick(ctx context.Context) {
	rows, err := s.pool.Query(ctx,
		`SELECT id FROM monitors WHERE status = 'active'`,
	)
	if err != nil {
		log.Printf("Scheduler query: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			log.Printf("scheduler scan: %v", err)
			continue
		}

		payload, err := NewSearchPayload(id)
		if err != nil {
			log.Printf("scheduler payload: %v", err)
			continue
		}

		task := asynq.NewTask(TypeSearchMonitor, payload)
		if _, err := s.client.EnqueueContext(ctx, task); err != nil {
			log.Printf("scheduler enqueue %s: %v", id, err)
		}
	}
}
