package jobs

import (
	"context"
	"time"

	"github.com/go-co-op/gocron"
)

func SetupCronJobs(ctx context.Context) {
	s := gocron.NewScheduler(time.UTC)

	go NotifyNifty50Stock(ctx, s)
}
