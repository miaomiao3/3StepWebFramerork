package step3

import (
	"github.com/miaomiao3/log"
	"mtest/http_router/step3/web"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//log.SetLogFuncCall(true)
	s := web.NewAPIServer()
	s.Start()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt, syscall.SIGKILL, syscall.SIGTERM)

	// Block until we receive our signal.
	<-c
	s.Stop()
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Debug("shutting down")
	os.Exit(0)

}
