package main

//还没有一些基础的东西没有搞明白，借鉴了别人的代码，但是还没有调好

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func StartHttpServer(srv *http.Server) error {
	http.HandleFunc("/demo", func(writer http.ResponseWriter, request *http.Request) {
		io.WriteString(writer, "hello, golang!\n")
	})
	fmt.Println("start")
	err := srv.ListenAndServe()
	return err
}

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	srv := &http.Server{Addr: ":8080"}
	// http server
	g.Go(func() error {
		fmt.Println("http")
		go func() {
			<-ctx.Done()
			fmt.Println("http ctx done")
			srv.Shutdown(context.TODO())
		}()
		return StartHttpServer(srv)
	})

	// signal
	g.Go(func() error {
		exitSignals := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT} // SIGTERM is POSIX specific
		sig := make(chan os.Signal, len(exitSignals))
		signal.Notify(sig, exitSignals...)
		for {
			fmt.Println("signal")
			select {
			case <-ctx.Done():
				fmt.Println("signal ctx done")
				return ctx.Err()
			case <-sig:
				// do something
				return nil
			}
		}
	})

	// inject error
	g.Go(func() error {
		fmt.Println("inject")
		time.Sleep(time.Second)
		fmt.Println("inject finish")
		return errors.New("inject error")
	})

	err := g.Wait() // first error return
	fmt.Println(err)
}
