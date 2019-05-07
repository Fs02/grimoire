package schema

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInferScanners(t *testing.T) {
	var (
		address = "address"
		record  = struct {
			ID      int
			Name    string
			Skip    bool `db:"-"`
			Number  float64
			Address *string
			Data    []byte
		}{
			ID:      1,
			Name:    "name",
			Number:  10.5,
			Address: &address,
			Data:    []byte("data"),
		}
		fields   = []string{"name", "id", "skip", "data", "number", "address", "not_exist"}
		expected = []interface{}{
			Value(&record.Name),
			Value(&record.ID),
			&sql.RawBytes{},
			Value(&record.Data),
			Value(&record.Number),
			&record.Address,
			&sql.RawBytes{},
		}
	)

	assert.Equal(t, expected, InferScanners(&record, fields))
}

func TestInferScanners_usingInterface(t *testing.T) {
	var (
		record = CustomSchema{
			UUID:  "abc123",
			Price: 100,
		}
		fields   = []string{"_uuid", "_price"}
		expected = []interface{}{Value(&record.UUID), Value(&record.Price)}
	)

	assert.Equal(t, expected, InferScanners(&record, fields))
}

func TestInferScanners_sqlScanner(t *testing.T) {
	var (
		record   = sql.NullBool{}
		fields   = []string{}
		expected = []interface{}{&sql.NullBool{}}
	)

	assert.Equal(t, expected, InferScanners(&record, fields))
}

func TestInferScanners_notPointer(t *testing.T) {
	assert.Panics(t, func() {
		InferScanners(struct{}{}, []string{})
	})
}
