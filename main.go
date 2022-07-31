package main

import (
	"flag"
	"go-kvm/internal/api"
	"go-kvm/internal/domain"
	"os"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

func main() {
	clientArgs := flag.NewFlagSet("client", flag.ExitOnError)
	clientConnect := clientArgs.String("connect", "", "ip address or hostname of server (ex 192.168.0.5:5772)")

	serverArgs := flag.NewFlagSet("server", flag.ExitOnError)
	serverListen := serverArgs.String("listen", "", "listen address and port (ex 0.0.0.0:5772)")

	switch os.Args[1] {
	case "client":
		err := clientArgs.Parse(os.Args[2:])
		if err != nil {
			log.Fatalf("failed to parse client args %v", err)
		}
		if *clientConnect == "" {
			log.Fatalln("no connect address provided")
		}
		log.Printf("client connecting to %q", *clientConnect)

		parts := strings.Split(*clientConnect, ":")
		client := api.NewClient(parts[0], parts[1])
		client.Connect()

	case "server":
		err := serverArgs.Parse(os.Args[2:])
		if err != nil {
			log.Fatalf("failed to parse server args %v", err)
		}
		if *serverListen == "" {
			log.Fatalln("no listen address provided")
		}
		log.Printf("server listening on %q", *serverListen)

		wg := &sync.WaitGroup{}
		wg.Add(1)

		inputChan := make(domain.InputChan)

		input := api.NewInput(inputChan)
		input.Start()
		log.Println(input.Devices)

		parts := strings.Split(*serverListen, ":")
		server := api.NewServer(parts[0], parts[1],
			"/home/mschon/certs/ca.crt",
			"/home/mschon/certs/server.crt",
			"/home/mschon/certs/server.key",
			inputChan, wg)

		server.Listen()

		wg.Wait()

	default:
		log.Printf("invalid argument: %s", os.Args[1])
	}
}
