package server

import (
	"bufio"
	"github.com/NicholasRodrigues/mini-db/internal/config"
	"net"
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func init() {
	viper.SetConfigFile("internal/config/test_data/test_config.yaml")
	config.LoadConfig()
}

func TestServerStartAndProcessCommand(t *testing.T) {
	err := os.Remove(config.Cfg.Storage.FilePath)
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("Failed to remove test storage file: %v", err)
	}

	server := NewServer()
	go server.Start()
	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", "localhost:9090")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("SET test_key test_value\n"))
	if err != nil {
		t.Fatal(err)
	}

	reader := bufio.NewReader(conn)
	setResponse, err := reader.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "OK\n", setResponse, "Expected OK response for SET command")

	time.Sleep(1 * time.Second)

	_, err = conn.Write([]byte("LOOKUP test_key\n"))
	if err != nil {
		t.Fatal(err)
	}

	lookupResponse, err := reader.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "test_value\n", lookupResponse, "Expected to retrieve the value set earlier")

	server.Stop()
	time.Sleep(10 * time.Second)
}

func TestServerPersistence(t *testing.T) {
	err := os.Remove(config.Cfg.Storage.FilePath)
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("Failed to remove test storage file: %v", err)
	}

	server := NewServer()
	go server.Start()
	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", "localhost:9090")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("SET persistent_key persistent_value\n"))
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)
	server.Stop()

	time.Sleep(1 * time.Second)

	server = NewServer()
	go server.Start()
	time.Sleep(1 * time.Second)

	conn, err = net.Dial("tcp", "localhost:9090")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("LOOKUP persistent_key\n"))
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "persistent_value\n", string(buf[:n]))
	server.Stop()
}
