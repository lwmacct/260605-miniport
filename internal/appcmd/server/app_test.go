package server

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/lwmacct/260605-miniport/internal/config"
)

func TestValidateHTTPTLS(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.Config
		wantErr bool
	}{
		{
			name: "disabled by default",
		},
		{
			name: "disabled ignores empty files",
			cfg: config.Config{
				Server: config.Server{
					HTTP: config.ServerHTTP{
						TLS: config.ServerHTTPTLS{
							CertFile: "",
							KeyFile:  "",
						},
					},
				},
			},
		},
		{
			name: "files mode with cert and key",
			cfg: config.Config{
				Server: config.Server{
					HTTP: config.ServerHTTP{
						TLS: config.ServerHTTPTLS{
							Enabled:      true,
							CertFile:     "cert.pem",
							KeyFile:      "key.pem",
							PollInterval: time.Second,
						},
					},
				},
			},
		},
		{
			name: "files mode cert without key rejected",
			cfg: config.Config{
				Server: config.Server{
					HTTP: config.ServerHTTP{
						TLS: config.ServerHTTPTLS{
							Enabled:  true,
							CertFile: "cert.pem",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "files mode negative poll interval rejected",
			cfg: config.Config{
				Server: config.Server{
					HTTP: config.ServerHTTP{
						TLS: config.ServerHTTPTLS{
							Enabled:      true,
							CertFile:     "cert.pem",
							KeyFile:      "key.pem",
							PollInterval: -time.Second,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "zero poll interval is valid",
			cfg: config.Config{
				Server: config.Server{
					HTTP: config.ServerHTTP{
						TLS: config.ServerHTTPTLS{
							Enabled:      true,
							CertFile:     "cert.pem",
							KeyFile:      "key.pem",
							PollInterval: 0,
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := NewApp(&test.cfg).validateHTTPTLS()
			if test.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
