package locales

import "embed"

// Files keeps locale resources available in single-binary deployments.
//
//go:embed *.yaml
var Files embed.FS
