package core

import (
	"context"
)

type (
	BasicCore interface {
		GetData(ctx context.Context) (string, error)
	}
	basicCore struct{}
)

func NewBasicCore() BasicCore {
	return &basicCore{}
}

func (bc *basicCore) GetData(ctx context.Context) (string, error) {
	return "from core", nil
}
