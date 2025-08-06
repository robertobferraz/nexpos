package repository

import (
	"github.com/robertobff/food-service/adapter/outbound/repository/src"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"repository",
	src.UserModule,
	src.ItemModule,
	src.CategoryModule,
	src.UserOrdersModule,
	src.DiscountModule,
	src.CityModule,
	src.CountryModule,
	src.StateModule,
	src.StreetModule,
	src.DistrictModule,
	src.UserAddressModule,
)
