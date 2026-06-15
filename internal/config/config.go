package config

import (
	"errors"
	"time"
)

type Config struct {
	Server Server `json:"server" desc:"服务端配置"`
}

type Server struct {
	Debug    bool           `json:"debug" desc:"启用调试日志和诊断信息"`
	Database ServerDatabase `json:"database" desc:"数据库配置"`
	Auth     ServerAuth     `json:"auth" desc:"认证配置"`
	HTTP     ServerHTTP     `json:"http" desc:"HTTP 服务配置"`
}

type ServerDatabase struct {
	Type   string              `json:"type" desc:"数据库类型：sqlite、pgsql"`
	SQLite string              `json:"sqlite" desc:"SQLite 数据库文件路径"`
	PGSQL  ServerDatabasePGSQL `json:"pgsql" desc:"PostgreSQL 连接参数"`
}

type ServerDatabasePGSQL struct {
	Host     string `json:"host" desc:"PostgreSQL 主机"`
	Port     string `json:"port" desc:"PostgreSQL 端口"`
	User     string `json:"user" desc:"PostgreSQL 用户名"`
	Database string `json:"database" desc:"PostgreSQL 数据库名"`
	Password string `json:"password" desc:"PostgreSQL 密码"`
}

type ServerHTTP struct {
	Listen          string        `json:"listen" desc:"HTTP 服务监听地址"`
	WebRoot         string        `json:"web-root" desc:"静态 Web 根目录，留空则不托管前端"`
	TLS             ServerHTTPTLS `json:"tls" desc:"HTTPS TLS 配置"`
	SessionTTL      time.Duration `json:"session-ttl" desc:"HTTP 登录会话有效期"`
	TrustedProxies  []string      `json:"trusted-proxies" desc:"可信 HTTP 反向代理 CIDR/IP 列表"`
	ReadTimeout     time.Duration `json:"read-timeout" desc:"HTTP 读取超时时间"`
	WriteTimeout    time.Duration `json:"write-timeout" desc:"HTTP 写入超时时间"`
	IdleTimeout     time.Duration `json:"idle-timeout" desc:"HTTP 空闲连接超时时间"`
	MaxAPIBodyBytes int64         `json:"max-api-body-bytes" desc:"HTTP API 最大请求体字节数，0 表示不限制"`
}

type ServerAuth struct {
	Admins []string        `json:"admins" desc:"运行时管理员用户名列表"`
	Local  ServerAuthLocal `json:"local" desc:"本地账号认证配置"`
}

type ServerAuthLocal struct {
	LoginEnabled        bool `json:"login-enabled" desc:"是否启用本地账号登录"`
	RegistrationEnabled bool `json:"registration-enabled" desc:"是否允许用户名密码公开注册"`
}

type ServerHTTPTLS struct {
	CertFile string `json:"cert-file" desc:"TLS 证书文件路径"`
	KeyFile  string `json:"key-file" desc:"TLS 私钥文件路径"`
}

func (cfg ServerHTTPTLS) Enabled() bool {
	return cfg.CertFile != "" || cfg.KeyFile != ""
}

func (cfg ServerHTTPTLS) Validate() error {
	if (cfg.CertFile == "") != (cfg.KeyFile == "") {
		return errors.New("http tls.cert-file and tls.key-file must be configured together")
	}
	return nil
}

func DefaultConfig() Config {
	return Config{
		Server: Server{
			Database: ServerDatabase{
				Type:   "sqlite",
				SQLite: "${APP_DATA:-.local/data}/sqlite.db",
				PGSQL: ServerDatabasePGSQL{
					Host:     "${PGHOST}",
					Port:     "${PGPORT}",
					User:     "${PGUSER}",
					Database: "${PGDATABASE}",
					Password: "${PGPASSWORD}",
				},
			},
			Auth: ServerAuth{
				Admins: []string{"admin"},
				Local: ServerAuthLocal{
					LoginEnabled:        true,
					RegistrationEnabled: true,
				},
			},
			HTTP: ServerHTTP{
				Listen:  ":40238",
				WebRoot: "${WEB_ROOT:-dist}",
				TLS: ServerHTTPTLS{
					CertFile: "",
					KeyFile:  "",
				},
				SessionTTL:      7 * 24 * time.Hour,
				TrustedProxies:  nil,
				ReadTimeout:     15 * time.Second,
				WriteTimeout:    30 * time.Second,
				IdleTimeout:     2 * time.Minute,
				MaxAPIBodyBytes: 1 << 20,
			},
		},
	}
}
