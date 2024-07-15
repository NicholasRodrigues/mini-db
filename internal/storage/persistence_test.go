package storage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/writer"
)

func TestPersistenceSaveAndLoad(t *testing.T) {
	filePath := "test_db.parquet"
	p := NewPersistence(filePath)
	defer os.Remove(filePath)

	dataToSave := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	err := p.Save(dataToSave)
	assert.NoError(t, err)

	loadedData, err := p.Load()
	assert.NoError(t, err)
	assert.Equal(t, dataToSave, loadedData)
}

func TestPersistenceLoadEmptyFile(t *testing.T) {
	filePath := "empty_test_db.parquet"
	defer os.Remove(filePath)

	// Create a valid empty Parquet file
	fw, err := local.NewLocalFileWriter(filePath)
	assert.NoError(t, err)
	pw, err := writer.NewParquetWriter(fw, new(KeyValue), 4)
	assert.NoError(t, err)
	assert.NoError(t, pw.WriteStop())
	assert.NoError(t, fw.Close())

	p := NewPersistence(filePath)

	// Load data from empty Parquet file
	loadedData, err := p.Load()
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{}, loadedData)
}

func TestPersistenceLoadNonExistentFile(t *testing.T) {
	filePath := "nonexistent_test_db.parquet"
	p := NewPersistence(filePath)

	// Load data from non-existent Parquet file
	loadedData, err := p.Load()
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{}, loadedData)
}
