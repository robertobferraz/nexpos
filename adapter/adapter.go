package adapter

import (
	"github.com/robertobff/food-service/adapter/connector"
	"github.com/robertobff/food-service/adapter/inbound"
	"github.com/robertobff/food-service/adapter/outbound"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"adapter",
	connector.Module,
	outbound.Module,
	inbound.Module,
)
