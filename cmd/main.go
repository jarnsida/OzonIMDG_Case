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
	// for future try of
	//ctx := context.Background()

	stop := make(chan os.Signal, 1)

	// Catch Ctrl+C, Ctrl+Z commands to stop server
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT /*, syscall.SIGTSTP*/)

	//Start server
	srv := server.NewServer()

	//Graceful Shutdown with data backup
	select {
	case <-stop:
		srv.Stop()
	}
	return nil
}
