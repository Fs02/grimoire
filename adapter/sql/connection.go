package sql

import (
	"github.com/tsenart/nap"
)

// Connection sql database
type Connection struct {
	*nap.DB
}

// NewConnection for sql database
func NewConnection(driver, dsns string) (*Connection, error) {
	db, err := nap.Open(driver, dsns)
	return &Connection{db}, err
}
