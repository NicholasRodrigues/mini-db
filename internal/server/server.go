package server

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/NicholasRodrigues/mini-db/internal/config"
	"github.com/NicholasRodrigues/mini-db/internal/storage"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type Server struct {
	address     string
	storage     *storage.Storage
	persistence *storage.Persistence
	mu          sync.Mutex
	listener    net.Listener
}

func init() {
	config.LoadConfig()
	level, err := logrus.ParseLevel(config.Cfg.Logging.Level)
	if err != nil {
		log.Fatalf("Invalid log level: %v", err)
	}
	log.SetLevel(level)
}

func NewServer() *Server {
	address := fmt.Sprintf(":%s", config.Cfg.Server.Port)
	newStorage := storage.NewStorage()
	persistence := storage.NewPersistence(config.Cfg.Storage.FilePath)

	data, err := persistence.Load()
	if err != nil {
		log.Fatalf("Failed to load persisted data: %v", err)
	}
	for key, value := range data {
		newStorage.Set(key, value)
	}

	return &Server{
		address:     address,
		storage:     newStorage,
		persistence: persistence,
	}
}

func (s *Server) Start() {
	var err error
	if config.Cfg.Server.TLS {
		cer, err := tls.LoadX509KeyPair(config.Cfg.Server.TLSCertFile, config.Cfg.Server.TLSKeyFile)
		if err != nil {
			log.Fatalf("Failed to load TLS certificates: %v", err)
		}
		tlsConfig := &tls.Config{Certificates: []tls.Certificate{cer}, MinVersion: tls.VersionTLS12}
		s.listener, err = tls.Listen("tcp", s.address, tlsConfig)
		if err != nil {
			log.Fatalf("Failed to start TLS listener: %v", err)
		}
	} else {
		s.listener, err = net.Listen("tcp", s.address)
		if err != nil {
			log.Fatalf("Failed to start listener: %v", err)
		}
	}

	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	log.Infof("Server listening on %s", s.address)
	defer s.listener.Close()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Errorf("Failed to accept connection: %v", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Errorf("Failed to close connection: %v", err)
		}
	}(conn)
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		if err := s.processCommand(line, conn); err != nil {
			log.Errorf("Error processing command: %v", err)
		}
	}
}

func (s *Server) processCommand(line string, conn net.Conn) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	parts := strings.Fields(line)
	if len(parts) < 1 {
		return fmt.Errorf("invalid command")
	}

	cmd := strings.ToUpper(parts[0])
	startIndex := 1

	if config.Cfg.Security.AuthEnabled {
		if len(parts) < 2 {
			return fmt.Errorf("authentication required")
		}
		authToken := parts[1]
		if authToken != config.Cfg.Security.AuthToken {
			if _, err := conn.Write([]byte("ERROR: Authentication failed\n")); err != nil {
				return err
			}
			return fmt.Errorf("authentication failed")
		}
		startIndex = 2
	}

	switch cmd {
	case "SET":
		if len(parts) < startIndex+2 {
			return fmt.Errorf("invalid SET command")
		}
		key := parts[startIndex]
		value := parts[startIndex+1]
		s.storage.Set(key, value)

		err := s.persistence.Save(s.storage.Store())
		if err != nil {
			log.Errorf("Failed to persist data: %v", err)
			return err
		}

		if _, err := conn.Write([]byte("OK\n")); err != nil {
			return err
		}
		log.Infof("SET command successful for key: %s", key)
	case "LOOKUP":
		if len(parts) < startIndex+1 {
			return fmt.Errorf("invalid LOOKUP command")
		}
		key := parts[startIndex]
		value, found := s.storage.Get(key)
		if !found {
			if _, err := conn.Write([]byte("NOT FOUND\n")); err != nil {
				return err
			}
			log.Infof("LOOKUP command: key %s not found", key)
		} else {
			if _, err := conn.Write([]byte(fmt.Sprintf("%s\n", value))); err != nil {
				return err
			}
			log.Infof("LOOKUP command successful for key: %s", key)
		}
	default:
		if _, err := conn.Write([]byte("ERROR: Unknown command\n")); err != nil {
			return err
		}
		log.Warnf("Unknown command received: %s", line)
		return fmt.Errorf("unknown command")
	}
	return nil
}

func (s *Server) Stop() {
	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			log.Errorf("Failed to close listener: %v", err)
		}
	}
}
