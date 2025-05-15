package entity

import "context"

type OrchestratorInterface interface {
	GetTemp(context.Context, CEP) (*Temp, error)
}
