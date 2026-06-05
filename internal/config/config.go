package config

// Config is the root application config.
type Config struct {
	Server Server `json:"server" desc:"服务端配置"`
}

// Server contains backend runtime settings.
//
//nolint:tagliatelle // public config keys are kebab-case.
type Server struct {
	Database string    `json:"database" desc:"SQLite 数据库文件路径"`
	HTTP     ServerHTTP `json:"http" desc:"HTTP 服务配置"`
}

// ServerHTTP contains HTTP listener settings.
//
//nolint:tagliatelle // public config keys are kebab-case.
type ServerHTTP struct {
	Listen string `json:"listen" desc:"HTTP 服务监听地址"`
}

// DefaultConfig returns a runnable baseline config.
func DefaultConfig() Config {
	return Config{
		Server: Server{
			Database: ".local/data/app.db",
			HTTP: ServerHTTP{
				Listen: ":8080",
			},
		},
	}
}

