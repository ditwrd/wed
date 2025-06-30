package db

import (
	"github.com/ditwrd/wed/internal/db/sqlite"
	"go.uber.org/fx"
)

var Module = fx.Options(
	sqlite.Module,
)
