package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
)

var version = "1.0.1"
var buildDate = "2025-09-17"

const keyServerAddr = "serverAddr"

func getRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	second := r.URL.Query().Get("second") // Getting Query Params
	body, err := io.ReadAll(r.Body)       // Getting BODY
	if err != nil {
		fmt.Printf("could not read body: %s\n", err)
	}
	fmt.Printf("%s: got / request. second=%s body:\n%s\n", ctx.Value(keyServerAddr), second, body)
	io.WriteString(w, "This is my website!\n")
}

func getHello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fmt.Printf("%s: got /hello request\n", ctx.Value(keyServerAddr))
	myName := r.PostFormValue("myName")
	if myName == "" {
		w.Header().Set("x-missing-field", "myName")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	io.WriteString(w, fmt.Sprintf("Hello, %s!\n", myName))
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fmt.Printf("%s: healthy request\n", ctx.Value(keyServerAddr))
	io.WriteString(w, "healthy")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", getHello)
	mux.HandleFunc("/health", getHealth)
	mux.HandleFunc("/", getRoot)
	ctx, cancelCtx := context.WithCancel(context.Background())

	serverOne := &http.Server{
		Addr:    ":3333",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx := context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	serverTwo := &http.Server{
		Addr:    ":4444",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx := context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	go func() {
		err := serverOne.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server one closed\n")
		} else if err != nil {
			fmt.Printf("error starting server: %s\n", err)
		}
		cancelCtx()
	}()

	go func() {
		err := serverTwo.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server two closed\n")
		} else if err != nil {
			fmt.Printf("error starting server: %s\n", err)
		}
		cancelCtx()
	}()

	<-ctx.Done()
}
