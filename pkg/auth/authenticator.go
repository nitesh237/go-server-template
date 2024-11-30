package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nitesh237/go-server-template/pkg/errors"
	"github.com/nitesh237/go-server-template/pkg/log"
	"github.com/shaj13/go-guardian/auth"
	"github.com/shaj13/go-guardian/auth/strategies/token"
	"go.uber.org/zap"
)

type Authenticator interface {
	GetHTTPMiddleware() func(next http.Handler) http.HandlerFunc
	GetGinMiddleware() gin.HandlerFunc
}

type staticBearerAuthenticator struct {
	authStrategy auth.Strategy
	logger       log.Logger
}

func NewStaticBearerAuthenticatorFromFile(filePath string, logger log.Logger) (Authenticator, error) {
	// TODO: add api scope in options check https://pkg.go.dev/github.com/shaj13/go-guardian/v2/auth/strategies/token#NewStatic
	strategy, err := token.NewStaticFromFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialise authenticator")
	}

	return &staticBearerAuthenticator{
		authStrategy: strategy,
		logger:       logger,
	}, nil
}

func (a *staticBearerAuthenticator) GetHTTPMiddleware() func(next http.Handler) http.HandlerFunc {
	return func(next http.Handler) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			a.logger.Debug(r.Context(), "Executing Auth Middleware")
			_, err := a.authStrategy.Authenticate(r.Context(), r)
			if err != nil {
				a.logger.Error(r.Context(), "authentication failed", zap.Error(err))
				code := http.StatusUnauthorized
				http.Error(w, http.StatusText(code), code)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (a *staticBearerAuthenticator) GetGinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// By pass auth for health check
		if c.Request.URL.Path == "/health" || c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		_, err := a.authStrategy.Authenticate(c.Request.Context(), c.Request)
		if err != nil {
			code := http.StatusUnauthorized
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.NewErrorResponseWithCode(http.StatusText(code), err.Error(), http.StatusUnauthorized))
			return
		}
		c.Next()
	}
}
