package sqlite

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewDBConnection),
	fx.Provide(NewSQLiteRSVPRepository),
	fx.Invoke(RegisterHooks),
)
