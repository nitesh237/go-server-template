package ginhttp

import (
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/nitesh237/go-gin-prometheus"
	"github.com/nitesh237/go-server-template/pkg/auth"
	"github.com/nitesh237/go-server-template/pkg/log"
	"go.uber.org/fx"
)

var (
	FxGinModule = fx.Module("gin",
		fx.Provide(
			gin.New,
			GinHttRouterProvider,
			GinHttpHandlerProvider,
		),
		fx.Decorate(
			func(router *gin.Engine, zapLogger log.ZapLogger) *gin.Engine {
				router.Use(
					ginzap.GinzapWithConfig(
						zapLogger.Unwrap(),
						&ginzap.Config{
							TimeFormat: time.RFC3339,
							UTC:        true,
							SkipPaths:  []string{"/health", "/metrics"},
						},
					),
					ginzap.RecoveryWithZap(zapLogger.Unwrap(), true))
				ginprometheus.NewPrometheus(ginprometheus.WithExcludedPaths("/health")).Use(router)
				return router
			},
		),
		fx.Invoke(
			RegisterHealthCheckEndpoint,
		),
	)

	FxAuthenticationModule = fx.Module("gin-http-authentication",
		fx.Provide(
			NewStaticBearerAuthenticatorFromFileProvider,
		),
		fx.Decorate(func(router *gin.Engine, authenticator auth.Authenticator) *gin.Engine {
			router.Use(authenticator.GetGinMiddleware())
			return router
		}),
		fx.Invoke(func(authenticator auth.Authenticator) {}, func(router *gin.Engine) {}),
	)
)

func GinHttRouterProvider(e *gin.Engine) GinHttpRouter {
	return e
}

func GinHttpHandlerProvider(e *gin.Engine) http.Handler {
	return e
}

func RegisterHealthCheckEndpoint(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

type StaticBearerAuthenticatorFromFileParams struct {
	fx.In

	ConfigPath string `name:"HttpAuthConfigPath"`
	Logger     log.Logger
}

type StaticBearerAuthenticatorFromFileResult struct {
	fx.Out

	Authenticator auth.Authenticator
}

func NewStaticBearerAuthenticatorFromFileProvider(r StaticBearerAuthenticatorFromFileParams) (StaticBearerAuthenticatorFromFileResult, error) {
	authenticator, err := auth.NewStaticBearerAuthenticatorFromFile(r.ConfigPath, r.Logger)
	if err != nil {
		return StaticBearerAuthenticatorFromFileResult{}, err
	}

	return StaticBearerAuthenticatorFromFileResult{
		Authenticator: authenticator,
	}, nil
}
