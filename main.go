package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg, err := ReadServerConfig("config.yml")
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	mc, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoDB.URI))
	if err != nil {
		log.Fatalf("failed to connect to mongodb: %v", err)
	}
	defer mc.Disconnect(context.Background())

	ss := NewStoreService(cfg.MongoDB, mc)
	if err := ss.CreateIndexes(context.Background()); err != nil {
		log.Fatalf("failed to create mongodb indexes: %v", err)
	}
	s := NewServer(cfg, ss)

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		log.Printf("server listening on %s", cfg.BindAddr)
		if err := s.Start(cfg.BindAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to run server: %v", err)
		}
	}()
	go func() {
		defer wg.Done()
		log.Printf("starting updater")
		if err := s.StartUpdater(ctx); err != nil {
			log.Fatalf("failed to run updater: %w", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	log.Printf("gracefully shutting down")
	cancel()
	if err := s.ShutdownWithTimeout(10 * time.Second); err != nil {
		log.Fatalf("failed to shutdown server: %v", err)
	}
	wg.Wait()
}
