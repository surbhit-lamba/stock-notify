package service

import (
	"context"
	"fmt"
	"os"
	"stock-notify/internal/repo"
	"stock-notify/internal/utils"
)

func NotifyNifty50Stock(ctx context.Context) {
	var fromEmail = os.Getenv("FROM_EMAIL")
	var toEmails = []string{os.Getenv("TO_EMAIL")}
	//weekDay := time.Now().Weekday()
	//if utils.SliceContains(weekDay, constants.WeekEnds) {
	//	log.ErrorfWithContext(ctx, "[NotifyNifty50Stock] Not a weekday, skipping sending mail")
	//	return
	//}
	response := repo.GetTimeSeriesMonthlyStockDataForStocks(ctx, []string{"TITAN.BSE"})
	fmt.Println(response)
	return
	utils.SendEmailWithHTMLTemplate(
		ctx,
		fromEmail,
		toEmails,
		"emailtemplates/notifyNifty50Stock.html",
		"[StockNotify] Stocks near 52 week low",
	)
}
