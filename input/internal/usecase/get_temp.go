package usecase

import (
	"context"

	"github.com/rafaelmascaro/Weather-By-CEP-With-Tracing/input/internal/entity"
)

type TempInputDTO struct {
	Zipcode string `json:"cep"`
}

type TempOutputDTO struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type GetTempUseCase struct {
	OrchestratorClient entity.OrchestratorClientInterface
}

func NewGetTempUseCase(
	orchestratorClient entity.OrchestratorClientInterface,
) *GetTempUseCase {
	return &GetTempUseCase{
		OrchestratorClient: orchestratorClient,
	}
}

func (g *GetTempUseCase) Execute(ctx context.Context, input TempInputDTO) (TempOutputDTO, error) {
	cep, err := entity.NewCEP(input.Zipcode)
	if err != nil {
		return TempOutputDTO{}, err
	}

	temp, err := g.OrchestratorClient.GetTemp(ctx, cep)
	if err != nil {
		return TempOutputDTO{}, err
	}

	dto := TempOutputDTO{
		City:  temp.City,
		TempC: temp.TempC,
		TempF: temp.TempF,
		TempK: temp.TempK,
	}

	return dto, nil
}
