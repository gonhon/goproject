/*
 * @Author: gaoh
 * @Date: 2024-07-10 22:38:08
 * @LastEditTime: 2024-07-10 22:41:14
 */
package main

import (
	_ "bookstore/internal/store"
	"bookstore/server"
	"bookstore/store/factory"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	s, err := factory.New("mem")
	if err != nil {
		panic(err)
	}

	bss := server.NewBookStoreServer(":8090", s)
	errChan, err := bss.ListenAndServe()
	if err != nil {
		log.Println("web serve start fail:", err)
		return
	}

	log.Println("web start ...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	select {
	case err = <-errChan:
		log.Println("web serve run fail:", err)
		return
	case <-c:
		log.Println("book store exit")
		ctx, cf := context.WithTimeout(context.Background(), time.Second)
		defer cf()
		bss.Shutdown(ctx)
	}
	if err != nil {
		log.Println("exit err:", err)
		return
	}
	log.Println("book program exit ok")
}
