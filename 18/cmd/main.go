package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wb_l12/18/config"
	"wb_l12/18/internal/handler"
	"wb_l12/18/internal/middleware"
	"wb_l12/18/internal/service"
	"wb_l12/18/pkg/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	cnf := config.Load()

	storage := storage.NewInMemoryStorage()
	service := service.NewService(storage)
	eventHandler := handler.NewEventHandler(service)

	router := gin.New()
	router.Use(middleware.Logging())

	router.POST("/create_event", eventHandler.CreateEvent)
	router.POST("/delete_event", eventHandler.DeleteEvent)
	router.POST("/update_event", eventHandler.UpdateEvent)
	router.GET("/events_for_day", eventHandler.GetByDay)
	router.GET("/events_for_week", eventHandler.GetByWeek)
	router.GET("/events_for_month", eventHandler.GetByMonth)

	srv := &http.Server{
		Addr:    net.JoinHostPort(cnf.Host, cnf.Port),
		Handler: router,
	}

	go func() {
		log.Printf("HTTP server run on http://%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error run server: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Stop server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Error stop server: %v", err)
		srv.Close()
	}
	log.Println("Server stopped")
}
