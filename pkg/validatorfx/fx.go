package validatorfx

import (
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
)

var (
	FxModule = fx.Module("validator", fx.Provide(validator.New))
)
