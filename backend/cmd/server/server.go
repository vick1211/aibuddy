package server

import (
	httptransport "aibuddy/internal/transport/http"
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"aibuddy/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Server struct{
	cfg config.Config
	log *slog.Logger
}

func New() (*Server, error){
	cfg, err := config.Load(os.Getenv)
	if err != nil{
		return nil, err
	}
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))
	return &Server{
		cfg: cfg,
		log: log,
	}, nil
}

func (s *Server) Start() error{
	db, err := sql.Open("pgx", s.cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			s.log.Error("failed to close database", "error", err)
		}
	}()

	pingCtx, cancelPing := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelPing()
	if err := db.PingContext(pingCtx); err != nil{
		s.log.Error("No response from database")
		return err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	srv := &http.Server{
		Addr:              	s.cfg.HTTPAddr,
		Handler:          	httptransport.NewRouter(httptransport.Deps{Log: s.log, DB: db}),
		ReadHeaderTimeout: 	5 * time.Second,
		ReadTimeout:       	10 * time.Second,
        WriteTimeout:     	30 * time.Second,
		IdleTimeout:       	60 * time.Second,
		ErrorLog: 			slog.NewLogLogger(s.log.Handler(), slog.LevelError),
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		s.log.Info("api started", "addr", s.cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed){
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
			return err
	case <- ctx.Done():
		s.log.Info("shutting down")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	 if err := srv.Shutdown(shutdownCtx); err != nil {
		s.log.Error("Shutdown")
        return err
    }
    return  <-errCh
}