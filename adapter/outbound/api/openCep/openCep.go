package openCep

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	nexpos "github.com/robertobff/nexpos/adapter/connector/openCep"
	"github.com/robertobff/nexpos/adapter/outbound/api/openCep/dto"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module("open_cep_api", fx.Provide(NewOpenCepApi))

type Api struct {
	nexpos *nexpos.OpenCep
	logger *zap.SugaredLogger
}

func NewOpenCepApi(logger *zap.SugaredLogger, client *nexpos.OpenCep) *Api {
	return &Api{
		nexpos: client,
		logger: logger,
	}
}

func (o *Api) GetByCep(cep *string) (*dto.GetCepResponse, error) {
	req, _ := http.NewRequest("GET", *o.nexpos.BaseUrl+*cep, nil)
	res, err := o.nexpos.Client().Do(req)
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
		o.logger.Errorw("failed to read response body", "error", err)
		return nil, err
	}

	var response dto.GetCepResponse

	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
