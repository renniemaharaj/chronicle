package migrations

import _ "embed"

//go:embed 001_initial.sql
var Foundation string

//go:embed 002_workspace.sql
var Workspace string

//go:embed 003_narratives.sql
var Narratives string
