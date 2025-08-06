package outbound

import (
	"github.com/robertobff/food-service/adapter/outbound/auth"
	"github.com/robertobff/food-service/adapter/outbound/database"
	"github.com/robertobff/food-service/adapter/outbound/logger"
	"github.com/robertobff/food-service/adapter/outbound/repository"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"outbound",
	logger.Module,
	database.Module,
	repository.Module,
	auth.Module,
)
