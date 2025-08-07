package application

import (
	"github.com/robertobff/nexpos/application/usecase"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"application",
	usecase.Module,
)
