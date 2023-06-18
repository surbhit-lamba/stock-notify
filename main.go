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
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
)

func main() {
	ctx := context.Background()

	nrApp, err := newrelic.NewApplication(
		newrelic.ConfigAppName(os.Getenv("NEWRELIC_NAME")),
		newrelic.ConfigLicense(os.Getenv("NEWRELIC_KEY")),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)

	log.Initialize(ctx)

	alphavantageClient := &httpclient.RequestClient{
		Identifier: httpclient.AlphaVantage,
		Host:       "https://www.alphavantage.co",
		Scheme:     "https",
		Authority:  "https://www.alphavantage.co",
	}

	ev := env.NewEnv(
		env.WithAlphaVantageHttpConn(alphavantageClient),
	)

	jobs.SetupCronJobs(ctx)

	r := router.SetupRouter(ctx, ev, nrApp)
	srv := &http.Server{
		Addr:         "0.0.0.0:81",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.ErrorfWithContext(ctx, "unable to start http server")
	}
}
