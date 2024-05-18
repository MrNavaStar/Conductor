package main

import (
	"github.com/docker/docker/client"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	socket, err := net.Listen("unix", "/tmp/conductor.sock")
	if err != nil {
		panic(err)
	}

	server := http.Server{
		Handler: getMux(cli),
		Addr:    "127.0.0.1:9999",
	}

	go func() {
		if err := server.Serve(socket); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	<-done

	server.Close()
	cli.Close()
	os.Remove("/tmp/conductor.sock")
	os.Exit(1)
}
