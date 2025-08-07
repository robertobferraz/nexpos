package connector

import (
	"github.com/robertobff/nexpos/adapter/connector/firebase"
	"github.com/robertobff/nexpos/adapter/connector/mercadopago"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"connector",
	mercadopago.Module,
	firebase.Module,
	//stripe.Module,
)
