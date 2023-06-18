package jobs

import (
	"context"
	"stock-notify/internal/constants"
	"stock-notify/internal/utils"
	"stock-notify/pkg/log"
	"time"

	"github.com/go-co-op/gocron"
)

func NotifyNifty50Stock(ctx context.Context, s *gocron.Scheduler) {
	_, err := s.Every(10).Seconds().Do(func() {
		weekDay := time.Now().Weekday()
		if utils.SliceContains(weekDay, constants.WeekDays) {
			utils.SendEmailWithHTMLTemplate(ctx)
		} else {
			log.ErrorfWithContext(ctx, "[NotifyNifty50Stock] Not a weekday, skipping sending mail")
		}
	})
	if err != nil {
		log.ErrorfWithContext(ctx, "Could not run NotifyNifty50Stock cron!!! Exiting")
		return
	}

	s.StartBlocking()
}
