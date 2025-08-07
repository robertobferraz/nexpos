package outbound

import (
	"github.com/robertobff/nexpos/adapter/outbound/auth"
	"github.com/robertobff/nexpos/adapter/outbound/database"
	"github.com/robertobff/nexpos/adapter/outbound/logger"
	"github.com/robertobff/nexpos/adapter/outbound/repository"
	"github.com/robertobff/nexpos/adapter/outbound/scheduler"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"outbound",
	logger.Module,
	database.Module,
	repository.Module,
	scheduler.Module,
	auth.Module,
)
