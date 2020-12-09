package main

import (
	"context"
	"fmt"
	"errors"
	"net/http"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var done = make(chan int)
	group, ctx := errgroup.WithContext(context.Background())
	group.Go(func() error {
		mux := http.NewServeMux()
		mux.HandleFunc("/close", func(writer http.ResponseWriter, request *http.Request) { done <- 1 })
		s := NewServer(":8080", mux)
		go func() {
			err := s.Start()
			if err != nil {
				fmt.Println(err)
			}
		}()

		select {
		case <-done:
			return s.Stop()
		case <- ctx.Done():
			fmt.Println("[notify] signal close")
			return nil
		}
	})

	group.Go(func() error {
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-quit:
			return errors.New("signal close")
		case <-ctx.Done():
			return errors.New("http server close")
		}
	})

	fmt.Println("start to catch signal")
	err := group.Wait()
	fmt.Println("===========", err)
}