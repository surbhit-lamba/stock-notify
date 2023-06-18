package service

import (
	"context"
	"fmt"
	"os"
	"stock-notify/internal/constants"
	"stock-notify/internal/utils"
	"stock-notify/pkg/log"
	"time"

	avClient "github.com/curtismckee/go-alpha-vantage"
)

func NotifyNifty50Stock(ctx context.Context) {
	timeNow := time.Now()
	//weekDay := timeNow.Weekday()
	//if utils.SliceContains(weekDay, constants.WeekEnds) {
	//	log.ErrorfWithContext(ctx, "[NotifyNifty50Stock] Not a weekday, skipping sending mail")
	//	return
	//}
	fmt.Println(timeNow)
	l, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.ErrorfWithContext(ctx, "unable to load location", err.Error())
		return
	}
	fmt.Println(timeNow.In(l))
	l, err = time.LoadLocation("Asia/Dubai")
	if err != nil {
		log.ErrorfWithContext(ctx, "unable to load location", err.Error())
		return
	}
	fmt.Println(timeNow.In(l))
	return
	av := avClient.NewClient(os.Getenv("ALPHAVANTAGE_API_KEY"))
	for _, symbol := range constants.Nifty50AlphaVantageSymbols {
		resp, err := av.StockTimeSeries(avClient.TimeSeriesMonthly, "TITAN.BSE")
		if err != nil {
			log.ErrorfWithContext(ctx, "[StockTimeSeries] err - ", err.Error())
		}
		// to make the array in descending order of time
		resp = utils.SliceReverse(resp)
		currentPrice := resp[0].Close
		var highPrice float64
		var lowPrice float64
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
				[]string{os.Getenv("TO_EMAIL")},
				"emailtemplates/notifyNifty50Stock.html",
				"[StockNotify] "+symbol+" near 52 week low",
			)
			log.InfofWithContext(ctx, "%v is %v percent away from low, sent email", symbol, diffPercentage)
		} else {
			log.InfofWithContext(ctx, "%v is %v percent away from low", symbol, diffPercentage)
		}
	}
	return
}
