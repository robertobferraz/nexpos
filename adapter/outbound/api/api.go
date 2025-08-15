package api

import (
	"github.com/robertobff/nexpos/adapter/outbound/api/countryStateCity"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"api",
	countryStateCity.Module,
)
