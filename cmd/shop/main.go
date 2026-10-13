package main 


import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"


	"github.com/gin-gonic/gin"

	"github.com/nikhil00shinde/atlas/internal/platform/config"
	"github.com/nikhil00shinde/atlas/internal/platform/db"
	"github.com/nikhil00shinde/atlas/internal/platform/httpx"
	"github.com/nikhil00shinde/atlas/internal/shop"

)


func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := run(log); err != nil {
		log.Error("shop exited", "err", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	//Cancelled on CTRL+C or SIGTERM 
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()


	cfg, err := config.Load()

	if err != nil {
		return err 
	}


	gdb, err := db.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		return err 
	}


	sqlDB, err := gdb.DB()
	if err != nil {
		return err 
	}


	defer sqlDB.Close()


	h := &shop.Handler{Store: &shop.Store{DB: gdb}, Log: log}


	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery(), httpx.Logging(log))


	router.GET("/products", h.ListProducts)
	router.GET("/healthz", h.Live)
	router.GET("/readyz", h.Ready)


	srv := &http.Server{
		Addr:								cfg.HTTPAddr,
		Handler:						router,
		ReadHeaderTimeout:	5 * time.Second,
		ReadTimeout:				10 * time.Second,
		WriteTimeout:				10 * time.Second,
		IdleTimeout:				60 * time.Second,
	}


	errCh := make(chan error, 1)
	go func() {
		log.Info("listening", "addr", cfg.HTTPAddr)
		errCh <- srv.ListenAndServe()
	}()


	select {
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			return err 
		}

	case <-ctx.Done():
		log.Info("shutting down")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)

}
