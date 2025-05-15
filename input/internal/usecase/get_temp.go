package usecase

import (
	"context"

	"github.com/rafaelmascaro/weather-api-otel/input/internal/entity"
)

type TempInput struct {
	Zipcode string `json:"cep"`
}

type TempOutput struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type GetTempUseCase struct {
	Orchestrator entity.OrchestratorInterface
}

func NewGetTempUseCase(
	Orchestrator entity.OrchestratorInterface,
) *GetTempUseCase {
	return &GetTempUseCase{
		Orchestrator: Orchestrator,
	}
}

func (g *GetTempUseCase) Execute(ctx context.Context, input TempInput) (TempOutput, error) {
	cep, err := entity.NewCEP(input.Zipcode)
	if err != nil {
		return TempOutput{}, err
	}

	temp, err := g.Orchestrator.GetTemp(ctx, cep)
	if err != nil {
		return TempOutput{}, err
	}

	dto := TempOutput{
		City:  temp.City,
		TempC: temp.TempC,
		TempF: temp.TempF,
		TempK: temp.TempK,
	}

	return dto, nil
}
