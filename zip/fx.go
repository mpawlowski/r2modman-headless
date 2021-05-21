package zip

import "go.uber.org/fx"

type Config struct {
}

func Module(c Config) fx.Option {
	return fx.Provide(
		newExtractor,
		func() Config {
			return c
		},
	)
}
