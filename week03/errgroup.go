package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	group, ctx := errgroup.WithContext(context.Background())

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(" world"))
	})

	server := http.Server{
		Addr:    "8090",
		Handler: mux,
	}

	//模拟退出的server
	serverOut := make(chan struct{})
	mux.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		serverOut <- struct{}{}
	})

	group.Go(func() error {
		fmt.Println("group1 server start...")
		return server.ListenAndServe()
	})

	group.Go(func() error {
		select {
		case <-ctx.Done():
			fmt.Println("group2 errgroup exit...")
		case <-serverOut:
			fmt.Println("group2 server will out...")
		}

		timeoutCtx, _ := context.WithTimeout(context.Background(), 3*time.Second)
		fmt.Println("group2 shutdown server...")
		return server.Shutdown(timeoutCtx)
	})

	//捕获到os.Signal后的退出
	group.Go(func() error {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGEMT)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case s := <-quit:
			return errors.Errorf("group3 get os signal: %v", s)

		}
	})

	fmt.Printf("errgroup exiting:%v\n", group.Wait())
}
