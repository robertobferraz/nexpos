package adapter

import (
	"github.com/robertobff/nexpos/adapter/connector"
	"github.com/robertobff/nexpos/adapter/inbound"
	"github.com/robertobff/nexpos/adapter/outbound"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"adapter",
	connector.Module,
	outbound.Module,
	inbound.Module,
)
