package inbound

import (
	"github.com/robertobff/food-service/adapter/inbound/http"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"inbound",
	http.Module,
)
