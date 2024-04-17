package controller

import (
	"context"
	"fmt"
	"github.com/golang-base-template/internal/basic_module/core"
	"github.com/golang-base-template/util/config"
)

type (
	BasicController interface {
		PostData(ctx context.Context) error
		GetData(ctx context.Context) (string, error)
	}
	basicController struct{}
)

var (
	basicCore core.BasicCore
	cfg       config.Config
)

func NewBasicController() BasicController {
	cfg = config.Get()
	if basicCore == nil {
		basicCore = core.NewBasicCore()
	}
	return &basicController{}
}

func (bc *basicController) PostData(ctx context.Context) error {
	return nil
}

func (bc *basicController) GetData(ctx context.Context) (result string, err error) {
	coreResult, err := basicCore.GetData(ctx)
	httpPort := cfg.Port.GBT
	if err != nil {
		return
	}
	result = fmt.Sprintf("Hello http world %s, via port: %s", coreResult, httpPort)
	return
}
