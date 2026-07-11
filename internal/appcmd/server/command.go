package server

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/lwmacct/251207-go-pkg-cfgm/pkg/cfgm"
	"github.com/lwmacct/251207-go-pkg-version/pkg/version"
	"github.com/lwmacct/260605-miniport/internal/config"
)

var (
	defaults = config.DefaultConfig()
	usage    = cfgm.Schema(defaults).Command("server")
)

var Command = &cli.Command{
	Name:            "server",
	Usage:           "start HTTP server",
	HideHelpCommand: true,
	Commands:        []*cli.Command{version.Command},
	Action:          action,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "debug",
			Usage: usage.MustUsage("debug"),
			Value: defaults.Server.Debug,
		},
		&cli.StringFlag{
			Name:  "database.type",
			Usage: usage.MustUsage("database.type"),
			Value: defaults.Server.Database.Type,
		},
		&cli.StringFlag{
			Name:  "database.sqlite",
			Usage: usage.MustUsage("database.sqlite"),
			Value: defaults.Server.Database.SQLite,
		},
		&cli.StringFlag{
			Name:  "database.pgsql.host",
			Usage: usage.MustUsage("database.pgsql.host"),
			Value: defaults.Server.Database.PGSQL.Host,
		},
		&cli.StringFlag{
			Name:  "database.pgsql.port",
			Usage: usage.MustUsage("database.pgsql.port"),
			Value: defaults.Server.Database.PGSQL.Port,
		},
		&cli.StringFlag{
			Name:  "database.pgsql.user",
			Usage: usage.MustUsage("database.pgsql.user"),
			Value: defaults.Server.Database.PGSQL.User,
		},
		&cli.StringFlag{
			Name:  "database.pgsql.database",
			Usage: usage.MustUsage("database.pgsql.database"),
			Value: defaults.Server.Database.PGSQL.Database,
		},
		&cli.StringFlag{
			Name:  "database.pgsql.password",
			Usage: usage.MustUsage("database.pgsql.password"),
			Value: defaults.Server.Database.PGSQL.Password,
		},
		&cli.StringSliceFlag{
			Name:  "auth.admins",
			Usage: usage.MustUsage("auth.admins"),
			Value: defaults.Server.Auth.Admins,
		},
		&cli.BoolFlag{
			Name:  "auth.local.login-enabled",
			Usage: usage.MustUsage("auth.local.login-enabled"),
			Value: defaults.Server.Auth.Local.LoginEnabled,
		},
		&cli.BoolFlag{
			Name:  "auth.local.registration-enabled",
			Usage: usage.MustUsage("auth.local.registration-enabled"),
			Value: defaults.Server.Auth.Local.RegistrationEnabled,
		},
		&cli.StringFlag{
			Name:  "auth.challenge.provider",
			Usage: usage.MustUsage("auth.challenge.provider"),
			Value: defaults.Server.Auth.Challenge.Provider,
		},
		&cli.IntFlag{
			Name:  "auth.challenge.image.max-challenges",
			Usage: usage.MustUsage("auth.challenge.image.max-challenges"),
			Value: defaults.Server.Auth.Challenge.Image.MaxChallenges,
		},
		&cli.StringFlag{
			Name:  "auth.challenge.hcaptcha.sitekey",
			Usage: usage.MustUsage("auth.challenge.hcaptcha.sitekey"),
			Value: defaults.Server.Auth.Challenge.HCaptcha.SiteKey,
		},
		&cli.StringFlag{
			Name:  "auth.challenge.hcaptcha.secret",
			Usage: usage.MustUsage("auth.challenge.hcaptcha.secret"),
			Value: defaults.Server.Auth.Challenge.HCaptcha.Secret,
		},
		&cli.StringFlag{
			Name:  "auth.challenge.hcaptcha.verify-url",
			Usage: usage.MustUsage("auth.challenge.hcaptcha.verify-url"),
			Value: defaults.Server.Auth.Challenge.HCaptcha.VerifyURL,
		},
		&cli.StringFlag{
			Name:  "auth.challenge.turnstile.sitekey",
			Usage: usage.MustUsage("auth.challenge.turnstile.sitekey"),
			Value: defaults.Server.Auth.Challenge.Turnstile.SiteKey,
		},
		&cli.StringFlag{
			Name:  "auth.challenge.turnstile.secret",
			Usage: usage.MustUsage("auth.challenge.turnstile.secret"),
			Value: defaults.Server.Auth.Challenge.Turnstile.Secret,
		},
		&cli.StringFlag{
			Name:  "auth.challenge.turnstile.verify-url",
			Usage: usage.MustUsage("auth.challenge.turnstile.verify-url"),
			Value: defaults.Server.Auth.Challenge.Turnstile.VerifyURL,
		},
		&cli.DurationFlag{
			Name:  "auth.session.ttl",
			Usage: usage.MustUsage("auth.session.ttl"),
			Value: defaults.Server.Auth.Session.TTL,
		},
		&cli.StringFlag{
			Name:  "auth.session.cookie.name",
			Usage: usage.MustUsage("auth.session.cookie.name"),
			Value: defaults.Server.Auth.Session.Cookie.Name,
		},
		&cli.StringFlag{
			Name:  "auth.session.cookie.path",
			Usage: usage.MustUsage("auth.session.cookie.path"),
			Value: defaults.Server.Auth.Session.Cookie.Path,
		},
		&cli.BoolFlag{
			Name:  "auth.session.cookie.secure",
			Usage: usage.MustUsage("auth.session.cookie.secure"),
			Value: defaults.Server.Auth.Session.Cookie.Secure,
		},
		&cli.BoolFlag{
			Name:  "github.enabled",
			Usage: usage.MustUsage("github.enabled"),
			Value: defaults.Server.GitHub.Enabled,
		},
		&cli.Int64Flag{
			Name:  "github.app-id",
			Usage: usage.MustUsage("github.app-id"),
			Value: defaults.Server.GitHub.AppID,
		},
		&cli.StringFlag{
			Name:  "github.app-slug",
			Usage: usage.MustUsage("github.app-slug"),
			Value: defaults.Server.GitHub.AppSlug,
		},
		&cli.StringFlag{
			Name:  "github.private-key-file",
			Usage: usage.MustUsage("github.private-key-file"),
			Value: defaults.Server.GitHub.PrivateKeyFile,
		},
		&cli.StringFlag{
			Name:  "github.webhook-secret",
			Usage: usage.MustUsage("github.webhook-secret"),
			Value: defaults.Server.GitHub.WebhookSecret,
		},
		&cli.StringFlag{
			Name:  "github.setup-return-url",
			Usage: usage.MustUsage("github.setup-return-url"),
			Value: defaults.Server.GitHub.SetupReturnURL,
		},
		&cli.DurationFlag{
			Name:  "github.reconcile-interval",
			Usage: usage.MustUsage("github.reconcile-interval"),
			Value: defaults.Server.GitHub.ReconcileInterval,
		},
		&cli.StringFlag{
			Name:  "http.listen",
			Usage: usage.MustUsage("http.listen"),
			Value: defaults.Server.HTTP.Listen,
		},
		&cli.StringFlag{
			Name:  "http.web-root",
			Usage: usage.MustUsage("http.web-root"),
			Value: defaults.Server.HTTP.WebRoot,
		},
		&cli.BoolFlag{
			Name:  "http.tls.enabled",
			Usage: usage.MustUsage("http.tls.enabled"),
			Value: defaults.Server.HTTP.TLS.Enabled,
		},
		&cli.StringFlag{
			Name:  "http.tls.cert-file",
			Usage: usage.MustUsage("http.tls.cert-file"),
			Value: defaults.Server.HTTP.TLS.CertFile,
		},
		&cli.StringFlag{
			Name:  "http.tls.key-file",
			Usage: usage.MustUsage("http.tls.key-file"),
			Value: defaults.Server.HTTP.TLS.KeyFile,
		},
		&cli.DurationFlag{
			Name:  "http.tls.poll-interval",
			Usage: usage.MustUsage("http.tls.poll-interval"),
			Value: defaults.Server.HTTP.TLS.PollInterval,
		},
		&cli.DurationFlag{
			Name:  "http.read-timeout",
			Usage: usage.MustUsage("http.read-timeout"),
			Value: defaults.Server.HTTP.ReadTimeout,
		},
		&cli.DurationFlag{
			Name:  "http.write-timeout",
			Usage: usage.MustUsage("http.write-timeout"),
			Value: defaults.Server.HTTP.WriteTimeout,
		},
		&cli.DurationFlag{
			Name:  "http.idle-timeout",
			Usage: usage.MustUsage("http.idle-timeout"),
			Value: defaults.Server.HTTP.IdleTimeout,
		},
		&cli.Int64Flag{
			Name:  "http.max-api-body-bytes",
			Usage: usage.MustUsage("http.max-api-body-bytes"),
			Value: defaults.Server.HTTP.MaxAPIBodyBytes,
		},
		&cli.StringSliceFlag{
			Name:  "http.trusted-proxies",
			Usage: usage.MustUsage("http.trusted-proxies"),
			Value: defaults.Server.HTTP.TrustedProxies,
		},
	},
}

func action(ctx context.Context, cmd *cli.Command) error {
	cfg := cfgm.MustLoad(ctx, config.DefaultConfig(), cfgm.Command(cmd))
	return NewApp(cfg).Run(ctx)
}
