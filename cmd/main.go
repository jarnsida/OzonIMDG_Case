package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jarsida/OzonIMDG_Case/server"
)

// Mat Ryer main() hack
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	// ctx := context.Background()

	// config set
	//	cfg := config.Get()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTSTP)

	srv := server.NewServer()

	select {
	case <-stop:
		srv.Stop()
	}
	return nil
}
