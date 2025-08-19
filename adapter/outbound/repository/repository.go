package repository

import (
	"github.com/robertobff/nexpos/adapter/outbound/repository/src"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"repository",
	src.UserModule,
	src.ItemModule,
	src.CategoryModule,
	src.OrderModule,
	src.DiscountModule,
	src.CityModule,
	src.CountryModule,
	src.StateModule,
	src.StreetModule,
	src.DistrictModule,
	src.AddressModule,
	src.ImageModule,
	src.CartModule,
	src.CartItemModule,
	src.PaymentModule,
	src.ReviewModule,
	src.OrderItemModule,
)
