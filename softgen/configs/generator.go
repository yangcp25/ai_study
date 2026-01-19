package configs

import (
	"embed"
)

//go:embed prompts template
var Config embed.FS
