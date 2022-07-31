package api

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"go-kvm/internal/domain"
	"io/ioutil"
	"net"
	"sync"

	log "github.com/sirupsen/logrus"
)

// Server will listen on a port and send input events to connected clients
type Server struct {
	// Host is the listen address (ex 0.0.0.0)
	Host string
	// Port is the listen port (ex 5772)
	Port string

	CACertPath     string
	ServerCertPath string
	ServerKeyPath  string

	// ServerWait ...
	ServerWait *sync.WaitGroup
	InputChan  domain.InputChan
}

// NewServer returns a new Server ready for use
func NewServer(host, port, caCert, serverCert, serverKey string, inputChan domain.InputChan, wg *sync.WaitGroup) *Server {
	server := &Server{
		Host:           host,
		Port:           port,
		ServerWait:     wg,
		CACertPath:     caCert,
		ServerCertPath: serverCert,
		ServerKeyPath:  serverKey,
		InputChan:      inputChan,
	}

	return server
}

func (s *Server) handleConnection(conn net.Conn) {
	log.Printf("server connected to client %s", conn.RemoteAddr())

	defer func() {
		if err := conn.Close(); err != nil {
			log.Errorf("failed to close connection %v", err)
		}
	}()

clientLoop:
	for {
		select {
		case inputEvents := <-s.InputChan:
			for _, e := range inputEvents {
				ev := fmt.Sprintf("%d %d %d %d", e.Type, e.Code, e.Value, e.Time.Nano())
				n, err := conn.Write([]byte(ev))
				if err != nil {
					log.Errorf("server write to client failed: %v", err)
					break clientLoop
				}
				log.Debugf("server wrote %d bytes to %s", n, conn.RemoteAddr())
			}
		}
	}
}

// Listen will start a routine which listens on the configured port
func (s *Server) Listen() {
	go func() {
		defer s.ServerWait.Done()

		caCertPem, err := ioutil.ReadFile(s.ServerCertPath)
		if err != nil {
			log.Fatalf("failed to read ca cert %v", err)
		}

		roots := x509.NewCertPool()
		if ok := roots.AppendCertsFromPEM(caCertPem); !ok {
			log.Fatalf("failed to add ca cert to pool")
		}

		cert, err := tls.LoadX509KeyPair(s.ServerCertPath, s.ServerKeyPath)
		if err != nil {
			log.Fatalf("failed to read ca certificate file %v", err)
		}

		cfg := &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.NoClientCert,
			ClientCAs:    roots,
		}

		listener, err := tls.Listen("tcp", fmt.Sprintf(":%s", s.Port), cfg)
		if err != nil {
			log.Fatalf("failed to start server %v", err)
		}

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Errorf("listener error handling connection %v", err)
				continue
			}
			go s.handleConnection(conn)
		}
	}()
}
