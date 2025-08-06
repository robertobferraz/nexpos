package application

import (
	"github.com/robertobff/food-service/application/usecase"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"application",
	usecase.Module,
)
