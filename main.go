// Copyright 2021, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2021/1/17

package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"proxy-go/config"
	"proxy-go/logger"
)

func main() {
	config.InitConfig()

	var handler http.Handler = ChanCreate()
	srv := &http.Server{
		Handler:   handler,
		TLSConfig: nil,
	}
	// http
	httpListener, err := net.Listen("tcp", ":8080")
	if err != nil {
		logger.Fatal(err.Error())
	}
	go func() {
		logger.Info("http server running...")
		if err := srv.Serve(httpListener); err != nil && err != http.ErrServerClosed {
			logger.Fatal(err.Error())
		}
	}()
	// https
	httpsListener, err := net.Listen("tcp", ":8081")
	if err != nil {
		logger.Fatal(err.Error())
	}
	go func() {
		logger.Info("https server running...")
		if err := srv.ServeTLS(httpsListener, "", ""); err != nil && err != http.ErrServerClosed {
			logger.Fatal(err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("server forced to shutdown:%v", err)
	}
	logger.Info("server exit!")
}

// a channel (just for the fun of it)
type Chan chan int

func ChanCreate() Chan {
	c := make(Chan)
	go func(c Chan) {
		for x := 0; ; x++ {
			c <- x
		}
	}(c)
	return c
}

func (ch Chan) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, fmt.Sprintf("channel send #%d\n", <-ch))
}
