package repositories

import (
	"github.com/GerryLon/newc/test/config"
	"github.com/GerryLon/newc/test/pkg/errors"
)

// ProRepository a repository for example
//go:generate go run ../../../newc
type ProRepository struct {
	conf      config.Config
	db        *database
	TableName string
	version   int
}

// FindByID find something by id
func (r *ProRepository) FindByID() errors.NoFound {
	return errors.NoFound{}
}
