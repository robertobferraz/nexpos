package connector

import (
	"github.com/robertobff/nexpos/adapter/connector/countryStateCity"
	"github.com/robertobff/nexpos/adapter/connector/firebase"
	"github.com/robertobff/nexpos/adapter/connector/mercadopago"
	"github.com/robertobff/nexpos/adapter/connector/openCep"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"connector",
	mercadopago.Module,
	firebase.Module,
	countryStateCity.Module,
	openCep.Module,
	//stripe.Module,
)
