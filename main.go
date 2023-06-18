package main

import (
	"context"
	"fmt"
	"net/http"
	"stock-notify/internal/constants"
	"stock-notify/internal/utils"
	"stock-notify/pkg/log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func main() {
	ctx := context.Background()

	nrApp, err := newrelic.NewApplication(
		newrelic.ConfigAppName("Stock notify"),
		newrelic.ConfigLicense("8b902c7cb9de77e972b811d71939d9f5aec3NRAL"),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)

	log.Initialize(ctx)

	go runCronJobs()

	r := gin.New()
	r.Use(nrgin.Middleware(nrApp))
	setRoutes(r.Group("stocks"), ctx)
	srv := &http.Server{
		Addr:         "0.0.0.0:81",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err = srv.ListenAndServe()
	if err != nil {
		fmt.Printf("unable to start http server")
	}
}

func runCronJobs() {
	s := gocron.NewScheduler(time.UTC)

	_, err := s.Every(10).Seconds().Do(func() {
		weekDay := time.Now().Weekday()
		if utils.SliceContains(weekDay, constants.WeekDays) {
			utils.SendEmail()
		} else {
			fmt.Println("not a weekday")
		}
	})
	if err != nil {
		return
	}

	s.StartBlocking()
}

func setRoutes(r *gin.RouterGroup, ctx context.Context) {
	r.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "yoo healthy!!")
	})
}
