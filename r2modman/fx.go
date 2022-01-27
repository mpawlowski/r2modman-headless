package r2modman

import (
	"time"

	"go.uber.org/fx"
)

type Config struct {
	InstallDirectory          string
	WorkDirectory             string
	ThunderstoreForceDownload bool
	ThunderstoreCDN           string
	ThunderstoreCDNTimeout    time.Duration
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
