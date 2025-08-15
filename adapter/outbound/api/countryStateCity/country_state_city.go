package countryStateCity

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	nexpos "github.com/robertobff/nexpos/adapter/connector/countryStateCity"
	"github.com/robertobff/nexpos/adapter/outbound/api/countryStateCity/assembler"
	"github.com/robertobff/nexpos/adapter/outbound/api/countryStateCity/dto"
	"github.com/robertobff/nexpos/adapter/outbound/database/redis"
	"github.com/robertobff/nexpos/utils"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"country_state_city_api",
	ConfigModule,
	assembler.Module,
	fx.Provide(NewCountryStateCityClient),
)

type API struct {
	nexPos                  *nexpos.CountryStateCity
	assembler               *assembler.CountryStateCity
	redis                   *redis.Redis
	logger                  *zap.SugaredLogger
	CountryStateCityBaseURL *string
}

func NewCountryStateCityClient(
	c *Config,
	logger *zap.SugaredLogger,
	client *nexpos.CountryStateCity,
	redis *redis.Redis,
	assembler *assembler.CountryStateCity,
) *API {
	return &API{
		nexPos:                  client,
		assembler:               assembler,
		logger:                  logger,
		redis:                   redis,
		CountryStateCityBaseURL: c.CountryStateCityBaseURL,
	}
}

func (c *API) getCountry() (*dto.GetCountryResponse, error) {
	req, _ := http.NewRequest("GET", *c.CountryStateCityBaseURL+"/countries", nil)

	req.Header.Add("X-CSCAPI-KEY", *c.nexPos.ApiKey)

	res, err := c.nexPos.Client().Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	var response dto.GetCountryResponse

	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *API) getState() (*dto.GetStateResponse, error) {
	req, _ := http.NewRequest("GET", *c.CountryStateCityBaseURL+"/states", nil)

	req.Header.Add("X-CSCAPI-KEY", *c.nexPos.ApiKey)

	res, err := c.nexPos.Client().Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	var response dto.GetStateResponse

	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *API) getCity(countryCode *string, stateCode *string) (*dto.GetCityResponse, error) {
	req, _ := http.NewRequest(
		"GET",
		*c.CountryStateCityBaseURL+"/countries/"+*countryCode+"/states/"+*stateCode+"/cities",
		nil,
	)

	req.Header.Add("X-CSCAPI-KEY", *c.nexPos.ApiKey)

	res, err := c.nexPos.Client().Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	var response dto.GetCityResponse

	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	for i := range response {
		response[i].StateCode = stateCode
		response[i].CountryCode = countryCode
	}

	return &response, nil
}

func (c *API) StartSync(ctx context.Context) error {
	lockKey := "sync:country_state_city:lock"
	locked, err := c.redis.SetNX(ctx, utils.PString(lockKey), utils.PString("1"), utils.PDuration(2*time.Hour))
	if err != nil {
		c.logger.Errorw("failed to lock country state city", "error", err)
		return fmt.Errorf("failed to acquire redis lock: %w", err)
	}
	if !*locked {
		c.logger.Warn("job already running, skipping")
		return nil
	}
	defer func(redis *redis.Redis, matchingKeys *string) {
		err := redis.Del(matchingKeys)
		if err != nil {
			c.logger.Errorw("failed to unlock country state city", "error", err)
		}
	}(c.redis, utils.PString(lockKey))

	countries, err := c.getCountry()
	if err != nil {
		c.logger.Errorw("failed to get country", "error", err)
		return err
	}

	states, err := c.getState()
	if err != nil {
		c.logger.Errorw("failed to get state", "error", err)
		return err
	}

	var mu sync.Mutex
	var allCities dto.GetCityResponse
	sem := make(chan struct{}, 10)
	var wg sync.WaitGroup
	for _, country := range *countries {
		for _, state := range *states {
			wg.Add(1)
			sem <- struct{}{}
			go func(countryCode, stateCode string) {
				defer wg.Done()
				defer func() { <-sem }()
				if *country.Iso2 == countryCode && stateCode == *state.Iso2 {
					cities, err := c.getCity(utils.PString(countryCode), utils.PString(stateCode))
					if err != nil {
						fmt.Printf("failed to fetch cities for %s-%s: %v\n", countryCode, stateCode, err)
						return
					}
					mu.Lock()
					allCities = append(allCities, *cities...)
					mu.Unlock()
				}
			}(*state.CountryCode, *state.Iso2)
		}
	}

	wg.Wait()

	err = c.assembler.Sync(ctx, countries, states, &allCities)
	if err != nil {
		return err
	}

	err = c.redis.Set(ctx, utils.PString("sync:country_state_city:last_run"), time.Now().UTC().Format(time.RFC3339), utils.PDuration(0))
	if err != nil {
		return err
	}

	return nil
}
