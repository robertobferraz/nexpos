package database

import (
	"github.com/robertobff/nexpos/adapter/outbound/database/postgres"
	"github.com/robertobff/nexpos/adapter/outbound/database/redis"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"database",
	postgres.Module,
	redis.Module,
)
