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
		fx.Provide(page.NewHandler),
		fx.Invoke(page.Router),
	)
)
