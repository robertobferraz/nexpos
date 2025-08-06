package connector

import (
	"github.com/robertobff/food-service/adapter/connector/firebase"
	"github.com/robertobff/food-service/adapter/connector/mercadopago"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"connector",
	mercadopago.Module,
	firebase.Module,
	//stripe.Module,
)
