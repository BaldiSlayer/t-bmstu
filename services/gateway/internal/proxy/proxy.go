package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"go.uber.org/zap"
)

type Route struct {
	Prefix string
	Target string
}

type ProxyRouter struct {
	routes []Route
	logger *zap.Logger
}

func NewProxyRouter(routes []Route, logger *zap.Logger) *ProxyRouter {
	validated := make([]Route, 0, len(routes))

	for _, r := range routes {
		_, err := url.Parse(r.Target)
		if err != nil {
			logger.Error("invalid route target", zap.String("target", r.Target), zap.Error(err))

			continue
		}

		validated = append(validated, r)
	}

	return &ProxyRouter{
		routes: validated,
		logger: logger,
	}
}

func (p *ProxyRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range p.routes {
		if strings.HasPrefix(r.URL.Path, route.Prefix) {
			targetURL, err := url.Parse(route.Target)
			if err != nil {
				p.logger.Error("failed to parse route target", zap.String("target", route.Target), zap.Error(err))
				http.Error(w, "bad gateway", http.StatusBadGateway)

				return
			}

			proxy := httputil.NewSingleHostReverseProxy(targetURL)
			r.Host = targetURL.Host

			p.logger.Info("proxying",
				zap.String("from", r.URL.Path),
				zap.String("to", targetURL.String()),
			)

			proxy.ServeHTTP(w, r)

			return
		}
	}

	http.NotFound(w, r)
}
