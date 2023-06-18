package jobs

import (
	"context"
	"stock-notify/internal/service"
	"stock-notify/pkg/log"

	"github.com/go-co-op/gocron"
)

func NotifyNifty50Stock(ctx context.Context, s *gocron.Scheduler) {
	_, err := s.Every(10).Seconds().Do(func() {
		service.NotifyNifty50Stock(ctx)
	})
	if err != nil {
		log.ErrorfWithContext(ctx, "Could not run NotifyNifty50Stock cron!!! Exiting")
		return
	}

	s.StartBlocking()
}
