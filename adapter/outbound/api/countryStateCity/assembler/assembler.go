package assembler

import (
	"context"
	"fmt"

	"github.com/robertobff/nexpos/adapter/outbound/api/countryStateCity/dto"
	dtoDomain "github.com/robertobff/nexpos/domain/dto"
	"github.com/robertobff/nexpos/domain/entity"
	"github.com/robertobff/nexpos/domain/repository"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"country_state_city_assembler",
	fx.Provide(NewAssembler),
)

type CountryStateCity struct {
	countryRepo repository.CountryRepository
	stateRepo   repository.StateRepository
	cityRepo    repository.CityRepository
	logger      *zap.SugaredLogger
}

func NewAssembler(
	countryRepo repository.CountryRepository,
	stateRepo repository.StateRepository,
	cityRepo repository.CityRepository,
	logger *zap.SugaredLogger,
) *CountryStateCity {
	return &CountryStateCity{
		countryRepo: countryRepo,
		stateRepo:   stateRepo,
		cityRepo:    cityRepo,
		logger:      logger,
	}
}

func (a *CountryStateCity) getCountry(ctx context.Context, dto *dto.GetCountryResponse) error {
	a.logger.Infow("assembling country request", "request", dto)
	var country []entity.Country
	for _, value := range *dto {
		a.logger.Debugw("get country info", "country", value)
		newCountry, err := entity.NewCountry(
			value.Name,
			value.Iso2,
			value.Iso3,
			value.PhoneCode,
			value.Capital,
			value.Currency,
			value.Emoji,
			value.ID,
		)

		if err != nil {
			a.logger.Errorw("error creating new country", "error", err)
			return err
		}

		country = append(country, *newCountry)
	}

	err := a.countryRepo.CreateBulk(ctx, &country)
	if err != nil {
		a.logger.Errorw("error creating country array", "error", err)
		return err
	}

	a.logger.Debugw("country array created", "count", len(country))
	return nil
}

func (a *CountryStateCity) getState(ctx context.Context, dto *dto.GetStateResponse) error {
	a.logger.Debugw("assembling state request", "request", dto)
	var state []entity.State

	countries, err := a.countryRepo.Get(ctx, &dtoDomain.GormQuery{})
	if err != nil {
		a.logger.Errorw("error getting country state", "error", err)
		return err
	}
	countryMap := make(map[string]entity.Country)
	for _, c := range *countries {
		countryMap[*c.Iso2] = c
	}

	for _, value := range *dto {
		a.logger.Debugw("get state info", "country", value)
		c, ok := countryMap[*value.CountryCode]
		if !ok {
			a.logger.Errorw("country not found", "country", value)
			continue
		}

		newState, err := entity.NewState(
			value.Name,
			value.Iso2,
			value.Id,
			&c,
		)
		if err != nil {
			a.logger.Errorw("error creating new state", "error", err)
			return err
		}

		state = append(state, *newState)
	}

	err = a.stateRepo.CreateBulk(ctx, &state)
	if err != nil {
		a.logger.Errorw("error creating state array", "error", err)
		return err
	}

	return nil
}

func (a *CountryStateCity) getCity(ctx context.Context, dto *dto.GetCityResponse) error {
	a.logger.Debugw("assembling city request", "request", dto)
	var city []entity.City

	staties, err := a.stateRepo.Get(ctx, &dtoDomain.GormQuery{
		Preload: &[]dtoDomain.GormPreload{{Field: "Country"}},
	})
	if err != nil {
		a.logger.Errorw("error getting state", "error", err)
		return err
	}

	statiesMap := make(map[string]entity.State)
	for _, s := range *staties {
		key := fmt.Sprintf("%s-%s", *s.Country.Iso2, *s.Iso2)
		statiesMap[key] = s
	}

	for _, value := range *dto {
		a.logger.Debugw("get state info", "city", value)
		key := fmt.Sprintf("%s-%s", *value.CountryCode, *value.StateCode)
		s, ok := statiesMap[key]
		if !ok {
			a.logger.Errorw("state not found", "state", value)
			continue
		}
		newCity, err := entity.NewCity(
			value.Name,
			value.Id,
			&s,
		)
		if err != nil {
			return err
		}

		city = append(city, *newCity)
	}

	err = a.cityRepo.CreateBulk(ctx, &city)
	if err != nil {
		a.logger.Errorw("error creating city array", "error", err)
		return err
	}

	return nil
}

func (a *CountryStateCity) Sync(ctx context.Context, country *dto.GetCountryResponse, state *dto.GetStateResponse, city *dto.GetCityResponse) error {
	a.logger.Debugw("assembling country request", "request", country)
	a.logger.Info("starting sync")
	err := a.getCountry(ctx, country)
	if err != nil {
		a.logger.Errorw("error getting country info", "error", err)
		return err
	}

	err = a.getState(ctx, state)
	if err != nil {
		a.logger.Errorw("error getting state info", "error", err)
		return err
	}

	err = a.getCity(ctx, city)
	if err != nil {
		a.logger.Errorw("error getting city info", "error", err)
		return err
	}

	a.logger.Info("completed")
	return nil
}
