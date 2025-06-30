package server

import (
	"github.com/ditwrd/wed/internal/server/page"
	"go.uber.org/fx"
)

var (
	Module = fx.Module("ServerModule",
		fx.Provide(NewServer),
		fx.Invoke(RegisterHooks),
	)

	PageModule = fx.Module("PageModule",
		fx.Invoke(page.Router),
		fx.Provide(page.NewHandler),
	)
)

