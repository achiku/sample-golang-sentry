package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/getsentry/raven-go"
)

func loggingMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}
	return http.HandlerFunc(fn)
}

func normalHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "normal")
}

func panicHandler(w http.ResponseWriter, r *http.Request) {
	err := errors.New("test error for sentry")
	panic(err)
	fmt.Fprintf(w, "panic")
}

func main() {

	http.Handle("/normal", loggingMiddleware(
		http.HandlerFunc(raven.RecoveryHandler(normalHandler))))
	http.Handle("/panic", loggingMiddleware(
		http.HandlerFunc(raven.RecoveryHandler(panicHandler))))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
