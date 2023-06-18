package jobs

import (
	"context"
	"fmt"
	"stock-notify/internal/service"
	"stock-notify/pkg/log"
	"stock-notify/pkg/newrelic"

	"github.com/go-co-op/gocron"
)

func NotifyNifty50Stock(ctx context.Context, s *gocron.Scheduler) {
	_, err := s.Every(80).Minute().Do(func() {
		txn := newrelic.StartTransaction(fmt.Sprintf("stock-notify/NotifyNifty50Stock"), nil, nil)
		defer func() {
			txn.End()
		}()
		nrctx := newrelic.NewContext(context.Background(), txn)
		service.NotifyNifty50Stock(nrctx)
	})
	if err != nil {
		log.ErrorfWithContext(ctx, "Could not run NotifyNifty50Stock cron!!! Exiting")
		newrelic.NoticeError(ctx, err)
		return
	}

	s.StartBlocking()
}
