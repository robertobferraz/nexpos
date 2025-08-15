package api

import (
	"github.com/robertobff/nexpos/adapter/outbound/api/countryStateCity"
	"github.com/robertobff/nexpos/adapter/outbound/api/openCep"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"api",
	countryStateCity.Module,
	openCep.Module,
)
