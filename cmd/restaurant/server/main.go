package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-feast/resty-backend/internal/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx, cancelFunc := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelFunc()

	addr := config.Addr()

	engine := gin.Default()

	engine.Use(gin.Recovery())
	log.Printf("Starting %s%s on addr: %s", serviceName, version, addr)

	routes(engine)

	serv := http.Server{
		Addr:    addr,
		Handler: engine.Handler(),
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		if err := serv.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Error shutting down server: %v", err)
		}
	}()

	err := serv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
