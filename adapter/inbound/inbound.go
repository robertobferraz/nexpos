package inbound

import (
	"github.com/robertobff/nexpos/adapter/inbound/http"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"inbound",
	http.Module,
)
