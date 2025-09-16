package migrations

import "embed"

//go:embed assets/*.sql
var Migrations embed.FS
