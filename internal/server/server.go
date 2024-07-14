package server

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"github.com/NicholasRodrigues/mini-db/internal/config"
	"github.com/NicholasRodrigues/mini-db/internal/storage"
	"net"
	"strings"
	"sync"

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

	// Load persisted data
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
	var ln net.Listener
	var err error

	if config.Cfg.Server.TLS {
		cer, err := tls.LoadX509KeyPair(config.Cfg.Server.TLSCertFile, config.Cfg.Server.TLSKeyFile)
		if err != nil {
			log.Fatalf("Failed to load TLS certificates: %v", err)
		}
		tlsConfig := &tls.Config{Certificates: []tls.Certificate{cer}}
		ln, err = tls.Listen("tcp", s.address, tlsConfig)
	} else {
		ln, err = net.Listen("tcp", s.address)
	}

	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	s.listener = ln
	log.Infof("Server listening on %s", s.address)

	for {
		conn, err := ln.Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && !opErr.Temporary() {
				log.Info("Server stopped accepting new connections")
				break
			}
			log.Errorf("Failed to accept connection: %v", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
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

	if config.Cfg.Security.AuthEnabled {
		if len(parts) < 2 || parts[0] != config.Cfg.Security.AuthToken {
			_, err := conn.Write([]byte("ERROR: Unauthorized\n"))
			if err != nil {
				return err
			}
			log.Warnf("Unauthorized access attempt from %s", conn.RemoteAddr().String())
			return fmt.Errorf("unauthorized access")
		}
		// Removing token from command for the next steps
		parts = parts[1:]
	}

	cmd := strings.ToUpper(parts[0])
	switch cmd {
	case "SET":
		if len(parts) != 3 {
			return fmt.Errorf("invalid SET command")
		}
		key := parts[1]
		value := parts[2]
		log.Printf("Received SET command for key: %s, value: %s", key, value)
		s.storage.Set(key, value)

		// Persist data
		err := s.persistence.Save(s.storage.Store())
		log.Printf("Persisted data: %v", s.storage.Store())
		if err != nil {
			log.Errorf("Failed to persist data: %v", err)
			return fmt.Errorf("failed to persist data: %v", err)
		}

		conn.Write([]byte("OK\n"))
		log.Infof("SET command successful for key: %s", key)
	case "LOOKUP":
		if len(parts) != 2 {
			return fmt.Errorf("invalid LOOKUP command")
		}
		key := parts[1]
		value, found := s.storage.Get(key)
		if !found {
			conn.Write([]byte("NOT FOUND\n"))
			log.Infof("LOOKUP command: key %s not found", key)
		} else {
			conn.Write([]byte(fmt.Sprintf("%s\n", value)))
			log.Infof("LOOKUP command successful for key: %s", key)
		}
	default:
		conn.Write([]byte("ERROR: Unknown command\n"))
		log.Warnf("Unknown command received: %s", line)
		return fmt.Errorf("unknown command")
	}
	return nil
}

func (s *Server) Stop() {
	if s.listener != nil {
		err := s.listener.Close()
		if err != nil {
			log.Errorf("Failed to close listener: %v", err)
		}
	}
}
