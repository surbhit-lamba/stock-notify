package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
)

func main() {
	go runCronJobs()
	fmt.Println("ye to ho gya")
	time.Sleep(time.Minute)
}

func hello(name string) {
	message := fmt.Sprintf("Hi, %v", name)
	fmt.Println(message)
}

func runCronJobs() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(1).Minute().Do(func() {
		hello("John Doe")
	})

	s.StartBlocking()
}

func setRoutes(r *gin.RouterGroup, ctx context.Context) {
	r.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "web--healthy")
	})
}
