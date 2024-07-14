package storage

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
	"github.com/xitongsys/parquet-go/writer"
)

// KeyValue represents the key-value pair structure for Parquet
type KeyValue struct {
	Key   string `parquet:"name=key, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	Value string `parquet:"name=value, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
}

// Persistence represents the persistence layer for the storage.
type Persistence struct {
	filePath string
	mu       sync.Mutex
}

// NewPersistence creates a new Persistence instance.
func NewPersistence(filePath string) *Persistence {
	return &Persistence{filePath: filePath}
}

// Save stores the key-value pairs to the Parquet file.
func (p *Persistence) Save(data map[string]string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	err := os.MkdirAll("data", os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	fw, err := local.NewLocalFileWriter(p.filePath)
	if err != nil {
		return fmt.Errorf("failed to create file writer: %w", err)
	}
	defer fw.Close()

	pw, err := writer.NewParquetWriter(fw, new(KeyValue), 4)
	if err != nil {
		return fmt.Errorf("failed to create parquet writer: %w", err)
	}
	defer func(pw *writer.ParquetWriter) {
		err := pw.WriteStop()
		if err != nil {
			log.Printf("Failed to stop parquet writer: %v", err)
		}
	}(pw)

	for key, value := range data {
		if err := pw.Write(&KeyValue{Key: key, Value: value}); err != nil {
			return fmt.Errorf("failed to write data to parquet file: %w", err)
		}
	}
	log.Println("Data successfully saved to Parquet file.")
	return nil
}

// Load loads the key-value pairs from the Parquet file.
func (p *Persistence) Load() (map[string]string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	data := make(map[string]string)
	fr, err := local.NewLocalFileReader(p.filePath)
	log.Println("File path: ", p.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return data, nil // Empty map if file does not exist
		}
		return nil, fmt.Errorf("failed to open file reader: %w", err)
	}
	defer fr.Close()

	pr, err := reader.NewParquetReader(fr, new(KeyValue), 4)
	if err != nil {
		return nil, fmt.Errorf("failed to create parquet reader: %w", err)
	}
	defer pr.ReadStop()

	numRows := pr.GetNumRows()
	if numRows == 0 {
		return data, nil
	}

	for i := int64(0); i < numRows; i++ {
		kv := make([]KeyValue, 1)
		if err := pr.Read(&kv); err != nil {
			return nil, fmt.Errorf("error reading from Parquet file: %w", err)
		}
		if len(kv) > 0 && kv[0].Key != "" {
			data[kv[0].Key] = kv[0].Value
		}
	}
	return data, nil
}
