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

	"github.com/joho/godotenv"
	"github.com/nskondratev/postgres-sqlpad-setup/webhooks/query"
	"github.com/nskondratev/postgres-sqlpad-setup/webhooks/sqlpad"
	"golang.org/x/sync/errgroup"
)

func main() {
	_ = godotenv.Load()

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConnsPerHost:   100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	sqlpadClient := sqlpad.NewClient("http://sqlpad:3000", "admin@sqlpad.com", "admin", httpClient)

	queryHandler := query.NewHandler(sqlpadClient, "admin@sqlpad.com")

	http.Handle("/query_created", queryHandler)

	srv := &http.Server{
		Addr:         os.Getenv("ADDR"),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Println("Start application")

	ctx, done := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)

	// Wait for interruption
	g.Go(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

		select {
		case sig := <-c:
			log.Printf("Receive interruption signal: %s\n", sig)
			done()
		case <-ctx.Done():
			log.Println("Close signal goroutine")
			return ctx.Err()
		}
		return nil
	})

	g.Go(func() error {
		c := make(chan error, 1)
		go func() {
			c <- srv.ListenAndServe()
		}()
		select {
		case err := <-c:
			if err != nil && err != http.ErrServerClosed {
				log.Fatalf("Failed to start HTTP server: %v", err)
				return err
			}
		case <-ctx.Done():
			log.Println("Close HTTP goroutine")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := srv.Shutdown(ctx); err != nil {
				log.Fatalf("Failed to shutdown HTTP server: %v", err)
			}
			return ctx.Err()
		}
		return nil
	})

	err := g.Wait()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Gracefully shutdown")
}
