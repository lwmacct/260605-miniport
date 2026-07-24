package config

import (
	"time"

	"github.com/lwmacct/251207-go-pkg-cfgm/pkg/cfgm"
	"github.com/lwmacct/260614-go-pkg-tlsreload/pkg/tlsreload"
	"github.com/lwmacct/260711-go-pkg-authme/pkg/authme"
	"github.com/lwmacct/260711-go-pkg-authme/pkg/authme/adapters/dexgithub"
	"github.com/lwmacct/260711-go-pkg-authme/pkg/authme/adapters/statictoken"
)

type Config struct {
	Server Server `json:"server" desc:"服务端配置"`
}

type Server struct {
	Debug    bool           `json:"debug" desc:"启用调试日志和诊断信息"`
	Database ServerDatabase `json:"database" desc:"数据库配置"`
	GitHub   ServerGitHub   `json:"github" desc:"GitHub App 集成配置"`
	HTTP     ServerHTTP     `json:"http" desc:"HTTP 服务配置"`
}

type ServerGitHub struct {
	Enabled           bool          `json:"enabled" desc:"启用 GitHub App 仓库同步"`
	AppID             int64         `json:"app-id" desc:"GitHub App ID"`
	AppSlug           string        `json:"app-slug" desc:"GitHub App slug"`
	PrivateKeyFile    string        `json:"private-key-file" desc:"GitHub App 私钥 PEM 文件"`
	WebhookSecret     string        `json:"webhook-secret" desc:"GitHub App Webhook 签名密钥"`
	SetupReturnURL    string        `json:"setup-return-url" desc:"GitHub 安装完成后的站内跳转地址"`
	ReconcileInterval time.Duration `json:"reconcile-interval" desc:"GitHub 仓库全量校准周期，0 表示不启用定时校准"`
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
	AuthMe          ServerAuthMe  `json:"authme" desc:"单用户访问认证配置"`
	TLS             ServerHTTPTLS `json:"tls" desc:"HTTPS TLS 配置"`
	TrustedProxies  []string      `json:"trusted-proxies" desc:"可信 HTTP 反向代理 CIDR/IP 列表"`
	ReadTimeout     time.Duration `json:"read-timeout" desc:"HTTP 读取超时时间"`
	WriteTimeout    time.Duration `json:"write-timeout" desc:"HTTP 写入超时时间"`
	IdleTimeout     time.Duration `json:"idle-timeout" desc:"HTTP 空闲连接超时时间"`
	MaxAPIBodyBytes int64         `json:"max-api-body-bytes" desc:"HTTP API 最大请求体字节数，0 表示不限制"`
}

type ServerAuthMe struct {
	PathPrefix        string               `json:"path-prefix" desc:"认证 HTTP 路由前缀"`
	Origins           []string             `json:"origins" desc:"允许浏览器访问应用的可信 origin"`
	Session           authme.SessionConfig `json:"session" desc:"加密浏览器 Session 配置"`
	StaticToken       statictoken.Config   `json:"statictoken" desc:"静态 Access Token 认证配置"`
	DexGitHub         dexgithub.Config     `json:"dexgithub" desc:"Dex GitHub OIDC 认证配置"`
	AllowedGitHubUser string               `json:"allowed-github-user" desc:"允许访问应用的唯一 GitHub 用户名"`
}

type ServerHTTPTLS struct {
	Enabled            bool                          `json:"enabled" desc:"是否启用 HTTPS TLS"`
	Certificates       []tlsreload.CertificateSource `json:"certificates" desc:"TLS 证书来源列表"`
	DefaultCertificate string                        `json:"default-certificate" desc:"未匹配 SNI 时使用的证书 ID"`
	PollInterval       time.Duration                 `json:"poll-interval" desc:"证书来源兜底轮询间隔"`
	RetryInterval      time.Duration                 `json:"retry-interval" desc:"证书重载失败后的重试间隔"`
}

func (c ServerHTTPTLS) ReloadConfig() tlsreload.Config {
	return tlsreload.Config{
		Certificates:       c.Certificates,
		DefaultCertificate: c.DefaultCertificate,
		PollInterval:       c.PollInterval,
		RetryInterval:      c.RetryInterval,
	}
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
			GitHub: ServerGitHub{ //nolint:gosec // WebhookSecret is an environment variable placeholder, not a credential.
				Enabled:           false,
				WebhookSecret:     "${GITHUB_APP_WEBHOOK_SECRET}",
				PrivateKeyFile:    "${APP_DATA:-.local/data}/github-app.pem",
				SetupReturnURL:    "/#/settings/github?github=connected",
				ReconcileInterval: 6 * time.Hour,
			},
			HTTP: ServerHTTP{
				Listen:  ":40238",
				WebRoot: "${WEB_ROOT:-dist}",
				AuthMe: ServerAuthMe{
					Origins: []string{
						"http://localhost:40238",
						"http://localhost:40239",
					},
					Session: authme.SessionConfig{
						Keys: []authme.SessionKey{{ID: "default", Secret: "${AUTHME_SESSION_KEY}"}}, //nolint:gosec // Environment variable placeholder.
						TTL:  24 * time.Hour,
					},
					StaticToken: func() statictoken.Config {
						cfg := statictoken.DefaultConfig()
						cfg.Enabled = true
						cfg.Credentials = []statictoken.Credential{{ID: "operator", Name: "Operator", Token: "${AUTHME_ACCESS_TOKEN}"}} //nolint:gosec // Environment variable placeholder.
						return cfg
					}(),
					DexGitHub: func() dexgithub.Config {
						config := dexgithub.DefaultConfig()
						config.Enabled = true
						config.Issuer = "https://2008.s.lwmacct.com:20088"
						config.ClientID = "260605-miniport"
						return config
					}(),
					AllowedGitHubUser: "lwmacct",
				},
				TLS: ServerHTTPTLS{
					Enabled:            false,
					DefaultCertificate: "default",
					Certificates: []tlsreload.CertificateSource{
						{
							ID:          "default",
							Certificate: "${APP_DATA:-.local/data}/ssl/fullchain.pem",
							PrivateKey:  "${APP_DATA:-.local/data}/ssl/privkey.pem",
						},
					},
					PollInterval: 3 * time.Second,
				},
				TrustedProxies:  nil,
				ReadTimeout:     15 * time.Second,
				WriteTimeout:    30 * time.Second,
				IdleTimeout:     2 * time.Minute,
				MaxAPIBodyBytes: 1 << 20,
			},
		},
	}
}

var Manager = cfgm.New(DefaultConfig(), cfgm.AppName("app"))
