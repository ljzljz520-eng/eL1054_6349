package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	app "schoolbusauth"
	"schoolbusauth/internal/domain"
	"syscall"
	"time"
)

func main() {
	address := flag.String("address", "127.0.0.1:8080", "HTTP listen address")
	data := flag.String("data", "schoolbus-auth.db", "bbolt database path")
	date := flag.String("date", "2026-09-01", "deterministic business date")
	flag.Parse()
	businessDate, err := time.Parse("2006-01-02", *date)
	if err != nil {
		log.Fatal(err)
	}
	application, err := app.Open(app.Config{DatabasePath: *data, BusinessDate: businessDate, OnComplete: func(bundle domain.AuthorizationBundle) error { return nil }})
	if err != nil {
		log.Fatal(err)
	}
	defer application.Close()
	server := &http.Server{Addr: *address, Handler: application.Handler, ReadHeaderTimeout: 5 * time.Second}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-done
		_ = server.Close()
	}()
	fmt.Printf("school bus authorization server listening on http://%s\n", *address)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
