package jobs

import (
	"context"
	"stock-notify/internal/service"
	"stock-notify/pkg/log"
	"stock-notify/pkg/newrelic"

	"github.com/go-co-op/gocron"
)

func NotifyNifty50Stock(ctx context.Context, s *gocron.Scheduler) {
	_, err := s.Every(80).Minute().Do(func() {
		defer newrelic.StartSegmentWithContext(ctx, "NotifyNifty50Stock").End()
		service.NotifyNifty50Stock(ctx)
	})
	if err != nil {
		log.ErrorfWithContext(ctx, "Could not run NotifyNifty50Stock cron!!! Exiting")
		newrelic.NoticeError(ctx, err)
		return
	}

	s.StartBlocking()
}
