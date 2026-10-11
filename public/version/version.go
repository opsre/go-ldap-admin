package version

var (
	Version   string // Version (for example, v1.0.0)
	GitCommit string // Git commit hash (short format)
	BuildTime string // Build time
	GoVersion string // Go compiler version
	Platform  string // Build platform (for example, linux/amd64)
)

// GetVersion gets version information
func GetVersion() map[string]string {
	return map[string]string{
		"version":   Version,
		"gitCommit": GitCommit,
		"buildTime": BuildTime,
		"goVersion": GoVersion,
		"platform":  Platform,
	}
}
