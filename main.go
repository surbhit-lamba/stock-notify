package main

import (
	"context"
	"net/http"
	"os"
	"stock-notify/internal/env"
	"stock-notify/internal/router"
	"stock-notify/jobs"
	"stock-notify/pkg/httpclient"
	"stock-notify/pkg/log"
	"stock-notify/pkg/newrelic"
	"time"
)

func main() {
	ctx := context.Background()

	nrApp := newrelic.Initialize(
		&newrelic.Options{
			Name:                   os.Getenv("NEWRELIC_NAME"),
			License:                os.Getenv("NEWRELIC_KEY"),
			Enabled:                true,
			CrossApplicationTracer: true,
			DistributedTracer:      true,
		},
	)

	log.Initialize(ctx)

	alphavantageClient := &httpclient.RequestClient{
		Identifier: httpclient.AlphaVantage,
		Host:       "www.alphavantage.co",
		Scheme:     "https",
		Authority:  "www.alphavantage.co",
	}

	ev := env.NewEnv(
		env.WithAlphaVantageHttpConn(alphavantageClient),
	)

	jobs.SetupCronJobs(ev.WithContext(ctx))

	r := router.SetupRouter(ctx, ev, nrApp)
	srv := &http.Server{
		Addr:         "0.0.0.0:81",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.ErrorfWithContext(ctx, "unable to start http server")
	}
}
