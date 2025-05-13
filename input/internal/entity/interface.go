package entity

import "context"

type OrchestratorClientInterface interface {
	GetTemp(context.Context, CEP) (*Temp, error)
}
