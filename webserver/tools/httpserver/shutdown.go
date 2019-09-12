package httpserver

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func ShutDown(server *http.Server) {
	var quitSignal = make(chan os.Signal)
	signal.Notify(quitSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)

	var (
		pid       = os.Getpid()
		signalMsg = <-quitSignal
	)
	log.Printf("server pid[%d] received shutdown signal: [%v]", pid, signalMsg)

	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("server pid[%d] shutdown failed, err: %v", pid, err)
	}

	<-ctx.Done()

	log.Printf("server pid[%d] stoped", pid)
}
