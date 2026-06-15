package httpauth

import (
	"net"
	"net/http"
	"net/netip"

	"github.com/lwmacct/260605-miniport/internal/domain/authsession"
)

type Service struct {
	secureCookies bool
	trusted       []netip.Prefix
}

func NewService(secureCookies bool, trustedProxies []string) *Service {
	return &Service{
		secureCookies: secureCookies,
		trusted:       parseTrustedProxies(trustedProxies),
	}
}

func (svc *Service) SecureCookies() bool {
	return svc != nil && svc.secureCookies
}

func (svc *Service) WrapHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request, ok := svc.SessionRequest(r)
		if ok {
			r = r.WithContext(ContextWithRequest(r.Context(), request))
		}
		next.ServeHTTP(w, r)
	})
}

func (svc *Service) SessionRequest(r *http.Request) (authsession.Request, bool) {
	ip, ok := svc.ClientIP(r)
	if !ok {
		return authsession.Request{}, false
	}
	return authsession.Request{
		IP:         ip.String(),
		Host:       r.Host,
		UserAgent:  r.UserAgent(),
		Method:     r.Method,
		Path:       r.URL.Path,
		RemoteAddr: r.RemoteAddr,
	}, true
}

func (svc *Service) ClientIP(r *http.Request) (netip.Addr, bool) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	remoteIP, ok := parseIP(host)
	if !ok {
		return netip.Addr{}, false
	}
	if svc == nil || len(svc.trusted) == 0 || !ipInPrefixes(remoteIP, svc.trusted) {
		return remoteIP, true
	}
	if ip, ok := parseXForwardedFor(r.Header.Get("X-Forwarded-For")); ok {
		return ip, true
	}
	if ip, ok := parseIP(r.Header.Get("X-Real-IP")); ok {
		return ip, true
	}
	return remoteIP, true
}
