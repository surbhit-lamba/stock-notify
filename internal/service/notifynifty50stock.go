package service

import (
	"context"
	"fmt"
	"os"
	"stock-notify/internal/constants"
	"stock-notify/internal/utils"
	"stock-notify/pkg/log"
	"stock-notify/pkg/newrelic"
	"time"

	avClient "github.com/curtismckee/go-alpha-vantage"
)

func NotifyNifty50Stock(ctx context.Context) {
	defer newrelic.StartSegmentWithContext(ctx, "NotifyNifty50Stock").End()
	timeNow := time.Now()
	l, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		newrelic.NoticeError(ctx, err)
		log.ErrorfWithContext(ctx, "unable to load location", err.Error())
		return
	}
	timeNow = timeNow.In(l)
	weekDay := timeNow.Weekday()
	if utils.SliceContains(weekDay, constants.WeekEnds) {
		newrelic.NoticeExpectedError(ctx, fmt.Errorf("weekend skip"))
		log.ErrorfWithContext(ctx, "[NotifyNifty50Stock] Not a weekday, skipping checking", timeNow)
		return
	}

	if (timeNow.Hour() < 9 || (timeNow.Hour() == 9 && timeNow.Minute() < 15)) || (timeNow.Hour() > 15 || (timeNow.Hour() == 15 && timeNow.Minute() > 30)) {
		newrelic.NoticeExpectedError(ctx, fmt.Errorf("non trading hours skip"))
		log.ErrorfWithContext(ctx, "[NotifyNifty50Stock] Not in trading window, skipping checking", timeNow)
		return
	}
	log.InfofWithContext(ctx, "[NotifyNifty50Stock] Proceeding for checks ", timeNow)
	av := avClient.NewClient(os.Getenv("ALPHAVANTAGE_API_KEY"))
	for _, symbol := range constants.Nifty50AlphaVantageSymbols {
		resp, err := av.StockTimeSeries(avClient.TimeSeriesMonthly, symbol)
		if err != nil {
			newrelic.NoticeError(ctx, err)
			log.ErrorfWithContext(ctx, "[StockTimeSeries] err - ", err.Error())
		}
		// to make the array in descending order of time
		resp = utils.SliceReverse(resp)
		currentPrice := resp[0].Close
		var highPrice float64
		var lowPrice float64 = 9999999999999999
		monthCounter := 1
		for _, monthData := range resp {
			if monthCounter > 12 {
				break
			}
			if monthData.High > highPrice {
				highPrice = monthData.High
			}
			if monthData.Low < lowPrice {
				lowPrice = monthData.Low
			}
			monthCounter++
		}
		highLowGap := highPrice - lowPrice
		currLowGap := currentPrice - lowPrice
		diffPercentage := (currLowGap / highLowGap) * 100
		if diffPercentage <= 10 {
			utils.SendEmailWithHTMLTemplate(
				ctx,
				os.Getenv("FROM_EMAIL"),
				[]string{os.Getenv("TO_EMAIL"), os.Getenv("TO_EMAIL_2")},
				"[StockNotify] "+symbol+" near 52 week low",
				"emailtemplates/notifyNifty50Stock.html",
				struct {
					StockName string
					CurrPrice float64
					HighPrice float64
					LowPrice  float64
				}{
					StockName: symbol,
					CurrPrice: currentPrice,
					HighPrice: highPrice,
					LowPrice:  lowPrice,
				},
			)
			log.InfofWithContext(ctx, "%v is %v percent away from low, sent email", symbol, diffPercentage)
		} else {
			log.InfofWithContext(ctx, "%v is %v percent away from low", symbol, diffPercentage)
		}
		time.Sleep(20 * time.Second)
	}
	return
}
