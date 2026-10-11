package config

import (
	_ "embed"
	"fmt"
	"os"
	"strconv"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

// System configuration, corresponding to the YAML file
// Viper includes mapstructure; YAML uses hyphens, which are converted to camel case

// Global configuration variable
var Conf = new(config)

//go:embed go-ldap-admin-priv.pem
var priv []byte

//go:embed go-ldap-admin-pub.pem
var pub []byte

type config struct {
	System   *SystemConfig `mapstructure:"system" json:"system"`
	Logs     *LogsConfig   `mapstructure:"logs" json:"logs"`
	Database *Database     `mapstructure:"database" json:"database"`
	Mysql    *MysqlConfig  `mapstructure:"mysql" json:"mysql"`
	// Casbin    *CasbinConfig    `mapstructure:"casbin" json:"casbin"`
	Jwt       *JwtConfig       `mapstructure:"jwt" json:"jwt"`
	RateLimit *RateLimitConfig `mapstructure:"rate-limit" json:"rateLimit"`
	I18n      *I18nConfig      `mapstructure:"i18n" json:"i18n"`
	Ldap      *LdapConfig      `mapstructure:"ldap" json:"ldap"`
	Email     *EmailConfig     `mapstructure:"email" json:"email"`
	DingTalk  *DingTalkConfig  `mapstructure:"dingtalk" json:"dingTalk"`
	WeCom     *WeComConfig     `mapstructure:"wecom" json:"weCom"`
	FeiShu    *FeiShuConfig    `mapstructure:"feishu" json:"feiShu"`
}

// Set up configuration loading
func InitConfig() {
	workDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("failed to read application directory:%s", err))
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/")
	// Read configuration information
	err = viper.ReadInConfig()

	// Reload configuration dynamically
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// Save the loaded configuration in the global variable Conf
		if err := viper.Unmarshal(Conf); err != nil {
			panic(fmt.Errorf("failed to initialize configuration file:%s", err))
		}
		normalizeConfig()
		// Load the RSA key
		Conf.System.RSAPublicBytes = pub
		Conf.System.RSAPrivateBytes = priv
	})

	if err != nil {
		panic(fmt.Errorf("failed to read configuration file:%s", err))
	}
	// Save the loaded configuration in the global variable Conf
	if err := viper.Unmarshal(Conf); err != nil {
		panic(fmt.Errorf("failed to initialize configuration file:%s", err))
	}
	normalizeConfig()
	// Load the RSA key
	Conf.System.RSAPublicBytes = pub
	Conf.System.RSAPrivateBytes = priv

	// Some settings are loaded through environment variables
	dbDriver := os.Getenv("DB_DRIVER")
	if dbDriver != "" {
		Conf.Database.Driver = dbDriver
	}
	mysqlHost := os.Getenv("MYSQL_HOST")
	if mysqlHost != "" {
		Conf.Mysql.Host = mysqlHost
	}
	mysqlUsername := os.Getenv("MYSQL_USERNAME")
	if mysqlUsername != "" {
		Conf.Mysql.Username = mysqlUsername
	}
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	if mysqlPassword != "" {
		Conf.Mysql.Password = mysqlPassword
	}
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")
	if mysqlDatabase != "" {
		Conf.Mysql.Database = mysqlDatabase
	}
	mysqlPort := os.Getenv("MYSQL_PORT")
	if mysqlPort != "" {
		Conf.Mysql.Port, _ = strconv.Atoi(mysqlPort)
	}

	ldapUrl := os.Getenv("LDAP_URL")
	if ldapUrl != "" {
		Conf.Ldap.Url = ldapUrl
	}
	ldapBaseDN := os.Getenv("LDAP_BASE_DN")
	if ldapBaseDN != "" {
		Conf.Ldap.BaseDN = ldapBaseDN
	}
	ldapAdminDN := os.Getenv("LDAP_ADMIN_DN")
	if ldapAdminDN != "" {
		Conf.Ldap.AdminDN = ldapAdminDN
	}
	ldapAdminPass := os.Getenv("LDAP_ADMIN_PASS")
	if ldapAdminPass != "" {
		Conf.Ldap.AdminPass = ldapAdminPass
	}
	ldapUserDN := os.Getenv("LDAP_USER_DN")
	if ldapUserDN != "" {
		Conf.Ldap.UserDN = ldapUserDN
	}
	ldapUserInitPassword := os.Getenv("LDAP_USER_INIT_PASSWORD")
	if ldapUserInitPassword != "" {

		Conf.Ldap.UserInitPassword = ldapUserInitPassword
	}
	ldapDefaultEmailSuffix := os.Getenv("LDAP_DEFAULT_EMAIL_SUFFIX")
	if ldapDefaultEmailSuffix != "" {
		Conf.Ldap.DefaultEmailSuffix = ldapDefaultEmailSuffix
	}
	ldapUserPasswordEncryptionType := os.Getenv("LDAP_USER_PASSWORD_ENCRYPTION_TYPE")
	if ldapUserPasswordEncryptionType != "" {
		Conf.Ldap.UserPasswordEncryptionType = ldapUserPasswordEncryptionType
	}
}

func normalizeConfig() {
	if Conf.I18n == nil {
		Conf.I18n = &I18nConfig{}
	}
	if Conf.I18n.DefaultLocale == "" {
		Conf.I18n.DefaultLocale = "zh-CN"
	}
	if len(Conf.I18n.SupportedLocales) == 0 {
		Conf.I18n.SupportedLocales = []string{"zh-CN", "en-US", "ja-JP", "es-ES", "ko-KR"}
	}
}

type SystemConfig struct {
	Mode            string `mapstructure:"mode" json:"mode"`
	UrlPathPrefix   string `mapstructure:"url-path-prefix" json:"urlPathPrefix"`
	Port            int    `mapstructure:"port" json:"port"`
	InitData        bool   `mapstructure:"init-data" json:"initData"`
	RSAPublicBytes  []byte `mapstructure:"-" json:"-"`
	RSAPrivateBytes []byte `mapstructure:"-" json:"-"`
}

type LogsConfig struct {
	Level      zapcore.Level `mapstructure:"level" json:"level"`
	Path       string        `mapstructure:"path" json:"path"`
	MaxSize    int           `mapstructure:"max-size" json:"maxSize"`
	MaxBackups int           `mapstructure:"max-backups" json:"maxBackups"`
	MaxAge     int           `mapstructure:"max-age" json:"maxAge"`
	Compress   bool          `mapstructure:"compress" json:"compress"`
}

type Database struct {
	Driver string `mapstructure:"driver" json:"driver"`
	Source string `mapstructure:"source" json:"source"`
}

type MysqlConfig struct {
	Username    string `mapstructure:"username" json:"username"`
	Password    string `mapstructure:"password" json:"password"`
	Database    string `mapstructure:"database" json:"database"`
	Host        string `mapstructure:"host" json:"host"`
	Port        int    `mapstructure:"port" json:"port"`
	Query       string `mapstructure:"query" json:"query"`
	LogMode     bool   `mapstructure:"log-mode" json:"logMode"`
	TablePrefix string `mapstructure:"table-prefix" json:"tablePrefix"`
	Charset     string `mapstructure:"charset" json:"charset"`
	Collation   string `mapstructure:"collation" json:"collation"`
}

// type CasbinConfig struct {
// 	ModelPath string `mapstructure:"model-path" json:"modelPath"`
// }

type JwtConfig struct {
	Realm      string `mapstructure:"realm" json:"realm"`
	Key        string `mapstructure:"key" json:"key"`
	Timeout    int    `mapstructure:"timeout" json:"timeout"`
	MaxRefresh int    `mapstructure:"max-refresh" json:"maxRefresh"`
}

type RateLimitConfig struct {
	FillInterval int64 `mapstructure:"fill-interval" json:"fillInterval"`
	Capacity     int64 `mapstructure:"capacity" json:"capacity"`
}

type I18nConfig struct {
	DefaultLocale    string   `mapstructure:"default-locale" json:"defaultLocale"`
	SupportedLocales []string `mapstructure:"supported-locales" json:"supportedLocales"`
}

type LdapConfig struct {
	Url                        string `mapstructure:"url" json:"url"`
	MaxConn                    int    `mapstructure:"max-conn" json:"maxConn"`
	BaseDN                     string `mapstructure:"base-dn" json:"baseDN"`
	AdminDN                    string `mapstructure:"admin-dn" json:"adminDN"`
	AdminPass                  string `mapstructure:"admin-pass" json:"adminPass"`
	UserDN                     string `mapstructure:"user-dn" json:"userDN"`
	UserInitPassword           string `mapstructure:"user-init-password" json:"userInitPassword"`
	GroupNameModify            bool   `mapstructure:"group-name-modify" json:"groupNameModify"`
	UserNameModify             bool   `mapstructure:"user-name-modify" json:"userNameModify"`
	DefaultEmailSuffix         string `mapstructure:"default-email-suffix" json:"defaultEmailSuffix"`
	UserPasswordEncryptionType string `mapstructure:"user-password-encryption-type" json:"userPasswordEncryptionType"`
	EnableSync                 bool   `mapstructure:"enable-sync" json:"enableSync"`
}
type EmailConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port string `mapstructure:"port" json:"port"`
	User string `mapstructure:"user" json:"user"`
	Pass string `mapstructure:"pass" json:"pass"`
	From string `mapstructure:"from" json:"from"`
}

type DingTalkConfig struct {
	AppKey        string   `mapstructure:"app-key" json:"appKey"`
	AppSecret     string   `mapstructure:"app-secret" json:"appSecret"`
	AgentId       string   `mapstructure:"agent-id" json:"agentId"`
	RootOuName    string   `mapstructure:"root-ou-name" json:"rootOuName"`
	Flag          string   `mapstructure:"flag" json:"flag"`
	EnableSync    bool     `mapstructure:"enable-sync" json:"enableSync"`
	DeptSyncTime  string   `mapstructure:"dept-sync-time" json:"deptSyncTime"`
	UserSyncTime  string   `mapstructure:"user-sync-time" json:"userSyncTime"`
	DeptList      []string `mapstructure:"dept-list" json:"deptList"`
	IsUpdateSyncd bool     `mapstructure:"is-update-syncd" json:"isUpdateSyncd"`
	ULeaveRange   uint     `mapstructure:"user-leave-range" json:"userLevelRange"`
}

type WeComConfig struct {
	Flag          string `mapstructure:"flag" json:"flag"`
	CorpID        string `mapstructure:"corp-id" json:"corpId"`
	AgentID       int    `mapstructure:"agent-id" json:"agentId"`
	CorpSecret    string `mapstructure:"corp-secret" json:"corpSecret"`
	EnableSync    bool   `mapstructure:"enable-sync" json:"enableSync"`
	DeptSyncTime  string `mapstructure:"dept-sync-time" json:"deptSyncTime"`
	UserSyncTime  string `mapstructure:"user-sync-time" json:"userSyncTime"`
	IsUpdateSyncd bool   `mapstructure:"is-update-syncd" json:"isUpdateSyncd"`
}

type FeiShuConfig struct {
	Flag          string   `mapstructure:"flag" json:"flag"`
	AppID         string   `mapstructure:"app-id" json:"appId"`
	AppSecret     string   `mapstructure:"app-secret" json:"appSecret"`
	EnableSync    bool     `mapstructure:"enable-sync" json:"enableSync"`
	DeptSyncTime  string   `mapstructure:"dept-sync-time" json:"deptSyncTime"`
	UserSyncTime  string   `mapstructure:"user-sync-time" json:"userSyncTime"`
	DeptList      []string `mapstructure:"dept-list" json:"deptList"`
	IsUpdateSyncd bool     `mapstructure:"is-update-syncd" json:"isUpdateSyncd"`
}
