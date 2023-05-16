package build

import (
	_ "embed"
)

// Version вервия сервиса.
//
//go:embed .version
var Version string
