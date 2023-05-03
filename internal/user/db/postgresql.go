package db

import (
	"go.mod/pkg/client/postgresql"
	"go.mod/pkg/logging"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}
