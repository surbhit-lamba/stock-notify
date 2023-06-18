package repo

import (
	"context"
	"os"
	"stock-notify/internal/env"
	"stock-notify/pkg/httpclient"
	"stock-notify/pkg/log"
)

func GetTimeSeriesMonthlyStockDataForStocks(ctx context.Context, stockSymbols []string) map[string]map[string]interface{} {
	e := env.FromContext(ctx)
	avClient := e.AlphaVantageHttpConn()
	response := make(map[string]map[string]interface{})
	for _, symbol := range stockSymbols {
		params := make(map[string]string)
		params["function"] = "TIME_SERIES_MONTHLY"
		params["symbol"] = symbol
		params["apikey"] = os.Getenv("ALPHAVANTAGE_API_KEY")
		request := &httpclient.RequestConfig{
			Method:      "GET",
			Path:        "/query",
			QueryParams: params,
		}
		responseDatum := make(map[string]interface{})
		err := avClient.MakeRequest(ctx, request, &responseDatum)
		if err != nil {
			log.ErrorfWithContext(ctx, "[GetTimeSeriesMonthlyStockDataForStocks] error - "+err.Error())
			continue
		}
		response[symbol] = responseDatum
	}
	return response
}
