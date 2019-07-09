package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/tslamic/go-mux-metrics/metrics"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	m := metrics.NewMetrics()

	r := mux.NewRouter()
	r.Use(perf(m))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sleep := randSeconds(5)
		log.Printf("sleeping for %dms", sleep.Nanoseconds()/1e6)
		time.Sleep(sleep)
		w.WriteHeader(http.StatusOK)
	})

	srv := &http.Server{Handler: r, Addr: ":8088"}
	conns := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(conns)
	}()
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("HTTP server ListenAndServe: %v", err)
	}
	<-conns
}

func perf(m metrics.Metrics) mux.MiddlewareFunc {
	return mux.MiddlewareFunc(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			defer func() {
				url := r.URL.Path
				elapsed := time.Since(start)
				avg := m.Put(url, elapsed)
				log.Printf("running avg for '%s': %dms", url, avg.Nanoseconds()/1e6)
			}()
			h.ServeHTTP(w, r)
		})
	})
}

func randSeconds(max int) time.Duration {
	min := 1
	rnd := rand.Intn(max-min) + min
	return time.Duration(rnd) * time.Second
}
