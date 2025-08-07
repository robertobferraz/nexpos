package database

import (
	"github.com/robertobff/nexpos/adapter/outbound/database/postgres"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"database",
	postgres.Module,
)
