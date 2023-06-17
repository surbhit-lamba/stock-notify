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

	ctx := context.Background()
	r := gin.New()
	setRoutes(r.Group("stocks"), ctx)
	srv := &http.Server{
		Addr:         "0.0.0.0:81",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err := srv.ListenAndServe()
	if err != nil {
		fmt.Printf("unable to start http server")
	}
}

func hello(name string) {
	message := fmt.Sprintf("Hi, %v", name)
	fmt.Println(message)
}

func runCronJobs() {
	s := gocron.NewScheduler(time.UTC)
	

	s.Every(10).Seconds().Do(func() {
		hello("John Doe")
	})

	s.StartBlocking()
}

func setRoutes(r *gin.RouterGroup, ctx context.Context) {
	r.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "yoo healthy!!")
	})
}
