// Package migrations embeds all SQL migration files into the binary at compile time.
// The binary is fully self-contained — no external migrations directory needed at runtime.
package migrations

import "embed"

// FS contains all *.sql files from this directory.
//
//go:embed *.sql
var FS embed.FS

// Dir is the directory within FS passed to iofs.New.
// Since the *.sql files live at the root of this package, use ".".
const Dir = "."
