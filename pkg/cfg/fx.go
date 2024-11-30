package cfg

import "go.uber.org/fx"

var (
	FxEnvModule = fx.Module("environment", fx.Provide(
		fx.Annotate(
			GetEnvironment,
			fx.ResultTags(`name:"Env"`),
		),
	))

	FxConfigModule = fx.Module("config",
		fx.Provide(
			fx.Annotate(
				GetConfigDir,
				fx.ResultTags(`name:"ConfigDir"`),
			),
		),
		fx.Supply(
			fx.Annotate(
				ConfigTypeYaml,
				fx.ResultTags(`name:"ConfigType"`),
			),
		),
	)
)
