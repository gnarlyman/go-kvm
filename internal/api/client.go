package api

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

// Client receives input events from server
// and emulates them on current system
type Client struct {
	Host string
	Port string
}

// NewClient returns a Client ready for use
func NewClient(host, port string) *Client {
	return &Client{
		Host: host,
		Port: port,
	}
}

// Connect will establish a connection to the server
// and emulate incoming input events
func (c *Client) Connect() {
	caCertPEM, err := ioutil.ReadFile("/home/mschon/certs/ca.crt")
	if err != nil {
		log.Fatalf("failed to read ca cert: %v", err)
	}
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(caCertPEM)
	if !ok {
		log.Fatalf("failed to add cert to pool")
	}

	cert, err := tls.LoadX509KeyPair("/home/mschon/certs/client.crt", "/home/mschon/certs/client.key")
	if err != nil {
		log.Fatalf("failed to load keypair: %v", err)
	}

	conf := tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
		RootCAs:            roots,
	}
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", c.Host, c.Port), &conf)
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Errorf("failed to close connection: %v", err)
		}
	}()

clientLoop:
	for {
		b := make([]byte, 64)
		n, err := conn.Read(b)
		if err != nil {
			log.Errorf("failed to read conn %v", err)
			break clientLoop
		}
		log.Printf("%d %s", n, b)
	}
}
