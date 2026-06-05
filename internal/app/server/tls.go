package server

import (
	"crypto/sha256"
	"crypto/tls"
	"fmt"
	"log/slog"
	"os"
	"sync"
)

type tlsCertificateReloader struct {
	certPath string
	keyPath  string

	mu    sync.Mutex
	hash  [32]byte
	cert  tls.Certificate
	ready bool
}

func newTLSCertificateReloader(certPath, keyPath string) (*tlsCertificateReloader, error) {
	r := &tlsCertificateReloader{
		certPath: certPath,
		keyPath:  keyPath,
	}
	if err := r.Reload(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *tlsCertificateReloader) GetCertificate(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.ready {
		return nil, fmt.Errorf("tls certificate not loaded")
	}

	return &r.cert, nil
}

func (r *tlsCertificateReloader) Reload() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	certPEM, keyPEM, hash, err := readTLSFiles(r.certPath, r.keyPath)
	if err != nil {
		return err
	}

	if r.ready && hash == r.hash {
		return nil
	}

	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return fmt.Errorf("load tls certificate: %w", err)
	}

	r.hash = hash
	r.cert = cert
	r.ready = true
	slog.Info("tls certificate loaded", "cert_file", r.certPath, "key_file", r.keyPath)
	return nil
}

func readTLSFiles(certPath, keyPath string) ([]byte, []byte, [32]byte, error) {
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return nil, nil, [32]byte{}, fmt.Errorf("read tls cert file: %w", err)
	}

	keyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, nil, [32]byte{}, fmt.Errorf("read tls key file: %w", err)
	}

	sumInput := make([]byte, 0, len(certPEM)+len(keyPEM)+1)
	sumInput = append(sumInput, certPEM...)
	sumInput = append(sumInput, 0)
	sumInput = append(sumInput, keyPEM...)
	sum := sha256.Sum256(sumInput)
	return certPEM, keyPEM, sum, nil
}
