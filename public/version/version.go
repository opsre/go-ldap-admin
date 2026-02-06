package version

var (
	Version   string // 版本号（如v1.0.0）
	GitCommit string // Git提交哈希（短格式）
	BuildTime string // 构建时间
	GoVersion string // Go编译器版本
	Platform  string // 编译平台（如linux/amd64）
)

// GetVersion 获取版本信息
func GetVersion() map[string]string {
	return map[string]string{
		"version":   Version,
		"gitCommit": GitCommit,
		"buildTime": BuildTime,
		"goVersion": GoVersion,
		"platform":  Platform,
	}
}
