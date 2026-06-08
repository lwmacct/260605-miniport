package config

import "errors"

// Config 是应用根配置。
type Config struct {
	Server Server `json:"server" desc:"服务端配置"`
}

// Server 包含后端运行配置。
//
//nolint:tagliatelle // 对外配置键使用 kebab-case。
type Server struct {
	DB      ServerDB      `json:"db" desc:"数据库配置"`
	HTTP    ServerHTTP    `json:"http" desc:"HTTP 服务配置"`
	Control ServerControl `json:"control" desc:"本地控制面配置"`
}

// ServerDB 包含数据库配置。
//
//nolint:tagliatelle // 对外配置键使用 kebab-case。
type ServerDB struct {
	Type   string        `json:"type" desc:"数据库类型: sqlite 或 pgsql"`
	SQLite string        `json:"sqlite" desc:"SQLite 数据库文件路径"`
	PGSQL  ServerDBPGSQL `json:"pgsql" desc:"PostgreSQL 数据库连接参数"`
}

// ServerDBPGSQL 包含 PostgreSQL 连接参数。
type ServerDBPGSQL struct {
	Host     string `json:"host" desc:"PostgreSQL 主机"`
	Port     string `json:"port" desc:"PostgreSQL 端口"`
	User     string `json:"user" desc:"PostgreSQL 用户"`
	Database string `json:"database" desc:"PostgreSQL 数据库名"`
	Password string `json:"password" desc:"PostgreSQL 密码"`
}

// ServerHTTP 包含 HTTP 监听配置。
//
//nolint:tagliatelle // 对外配置键使用 kebab-case。
type ServerHTTP struct {
	Listen      string `json:"listen" desc:"HTTP 服务监听地址"`
	SSLCertFile string `json:"ssl-cert-file" desc:"SSL 证书文件路径"`
	SSLKeyFile  string `json:"ssl-key-file" desc:"SSL 密钥文件路径"`
}

// ServerControl 包含本地控制面 socket 配置。
//
//nolint:tagliatelle // 对外配置键使用 kebab-case。
type ServerControl struct {
	Listen string `json:"listen" desc:"本地控制命令 Unix socket 路径"`
}

// Validate 检查 HTTP 配置内部一致性。
func (h ServerHTTP) Validate() error {
	if (h.SSLCertFile == "") != (h.SSLKeyFile == "") {
		return errors.New("http ssl-cert-file and ssl-key-file must be configured together")
	}
	return nil
}

// UsesTLS 判断是否启用 HTTPS。
func (h ServerHTTP) UsesTLS() bool {
	return h.SSLCertFile != "" && h.SSLKeyFile != ""
}

// DefaultConfig 返回可运行的默认配置。
func DefaultConfig() Config {
	return Config{
		Server: Server{
			DB: ServerDB{
				Type:   "sqlite",
				SQLite: ".local/data/app.db",
				PGSQL: ServerDBPGSQL{
					Host:     "${PGHOST}",
					Port:     "${PGPORT}",
					User:     "${PGUSER}",
					Database: "${PGDATABASE}",
					// #nosec G101 -- shell 模板占位符引用环境变量，不是硬编码密码。
					Password: "${PGPASSWORD}",
				},
			},
			HTTP: ServerHTTP{
				Listen: ":40238",
			},
			Control: ServerControl{
				Listen: ".local/run/control.sock",
			},
		},
	}
}
