package config

import "time"

type Config struct {
	Server Server `json:"server" desc:"服务端配置"`
}

type Server struct {
	Debug    bool           `json:"debug" desc:"启用调试日志和诊断信息"`
	Database ServerDatabase `json:"database" desc:"数据库配置"`
	Auth     ServerAuth     `json:"auth" desc:"认证配置"`
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
	TLS             ServerHTTPTLS `json:"tls" desc:"HTTPS TLS 配置"`
	TrustedProxies  []string      `json:"trusted-proxies" desc:"可信 HTTP 反向代理 CIDR/IP 列表"`
	ReadTimeout     time.Duration `json:"read-timeout" desc:"HTTP 读取超时时间"`
	WriteTimeout    time.Duration `json:"write-timeout" desc:"HTTP 写入超时时间"`
	IdleTimeout     time.Duration `json:"idle-timeout" desc:"HTTP 空闲连接超时时间"`
	MaxAPIBodyBytes int64         `json:"max-api-body-bytes" desc:"HTTP API 最大请求体字节数，0 表示不限制"`
}

type ServerAuth struct {
	Admins    []string            `json:"admins" desc:"运行时管理员用户名列表"`
	Local     ServerAuthLocal     `json:"local" desc:"本地账号认证配置"`
	Challenge ServerAuthChallenge `json:"challenge" desc:"认证挑战配置"`
	Session   ServerAuthSession   `json:"session" desc:"认证会话配置"`
}

type ServerAuthLocal struct {
	LoginEnabled        bool `json:"login-enabled" desc:"是否启用本地账号登录"`
	RegistrationEnabled bool `json:"registration-enabled" desc:"是否允许用户名密码公开注册"`
}

type ServerAuthChallenge struct {
	Provider  string                    `json:"provider" desc:"认证挑战提供方：image、hcaptcha、turnstile"`
	Image     ServerAuthChallengeImage  `json:"image" desc:"图片验证码配置"`
	HCaptcha  ServerAuthChallengeRemote `json:"hcaptcha" desc:"hCaptcha 挑战配置"`
	Turnstile ServerAuthChallengeRemote `json:"turnstile" desc:"Cloudflare Turnstile 挑战配置"`
}

type ServerAuthChallengeImage struct {
	MaxChallenges int `json:"max-challenges" desc:"内存中图片验证码最大数量，0 表示不限制"`
}

type ServerAuthChallengeRemote struct {
	SiteKey   string `json:"sitekey" desc:"认证挑战站点公钥"`
	Secret    string `json:"secret" desc:"认证挑战服务端密钥"`
	VerifyURL string `json:"verify-url" desc:"认证挑战服务端验证地址"`
}

type ServerAuthSession struct {
	TTL    time.Duration           `json:"ttl" desc:"登录会话有效期"`
	Cookie ServerAuthSessionCookie `json:"cookie" desc:"登录会话 Cookie 策略"`
}

type ServerAuthSessionCookie struct {
	Name   string `json:"name" desc:"登录会话 Cookie 名称"`
	Path   string `json:"path" desc:"登录会话 Cookie Path"`
	Secure bool   `json:"secure" desc:"是否给登录会话 Cookie 设置 Secure 属性"`
}

type ServerHTTPTLS struct {
	Enabled      bool          `json:"enabled" desc:"是否启用 HTTPS TLS"`
	CertFile     string        `json:"cert-file" desc:"TLS 证书文件路径或 URI"`
	KeyFile      string        `json:"key-file" desc:"TLS 私钥文件路径或 URI"`
	PollInterval time.Duration `json:"poll-interval" desc:"TLS 证书文件重载兜底轮询间隔，未配置时使用默认间隔"`
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
				Challenge: ServerAuthChallenge{
					Provider: "image",
					Image: ServerAuthChallengeImage{
						MaxChallenges: 1024,
					},
					HCaptcha: ServerAuthChallengeRemote{
						VerifyURL: "https://api.hcaptcha.com/siteverify",
					},
					Turnstile: ServerAuthChallengeRemote{
						VerifyURL: "https://challenges.cloudflare.com/turnstile/v0/siteverify",
					},
				},
				Session: ServerAuthSession{
					TTL: 7 * 24 * time.Hour,
					Cookie: ServerAuthSessionCookie{
						Name: "web_session",
						Path: "/api",
					},
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
				TLS: ServerHTTPTLS{
					Enabled:      false,
					CertFile:     "${APP_DATA:-.local/data}/ssl/fullchain.pem",
					KeyFile:      "${APP_DATA:-.local/data}/ssl/privkey.pem",
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
