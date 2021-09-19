package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	done := make(chan int)
	handle := func(w http.ResponseWriter, r *http.Request) {
		done <- 1
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("handled"))
	}
	
	creat := func(p context.Context, port int, handle http.HandlerFunc) error {
		// 创建上下文
		ctx, cancel := context.WithCancel(p)
		// 创建 error group
		eg, ec := errgroup.WithContext(ctx)
		// 创建服务
		s := &http.Server{Addr: ":" + strconv.Itoa(port)}
		// 启动服务
		eg.Go(func() error {
			http.HandleFunc("/close", handle)
			return s.ListenAndServe()
		})

		// 错误处理
		eg.Go(func() error {
			select {
			case <- done:
				fmt.Println("server out...")
			case <-ec.Done():
				return ec.Err()
			}
			defer cancel()
			fmt.Println("server shutdown")
			return s.Shutdown(ec)
		})

		// 信号处理
		eg.Go(func() error {
			sChan := make(chan os.Signal, 1)
			signal.Notify(sChan, syscall.SIGINT, syscall.SIGTERM)
			for  {
				select {
				case sig := <-sChan:
					return fmt.Errorf("os signal: %v", sig)
				case <-ec.Done():
					return ec.Err()
				}
			}
		})

		return eg.Wait()
	}
	err := creat(context.Background(), 8081, handle)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("exec over!")
	}
}
