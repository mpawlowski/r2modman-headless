package r2modman

import "go.uber.org/fx"

type Config struct {
	InstallDirectory   string
	WorkDirectory      string
	ThunderstoreCDNUrl string
}

func Module(c Config) fx.Option {
	return fx.Provide(
		newExportParser,
		newModUtil,
		func() Config {
			return c
		},
	)
}
