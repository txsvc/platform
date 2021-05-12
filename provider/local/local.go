package local

import (
	"context"
	"fmt"
	"log"
	h "net/http"

	"go.uber.org/zap"

	"github.com/txsvc/platform/v2"
	"github.com/txsvc/platform/v2/auth"
	"github.com/txsvc/platform/v2/errorreporting"
	"github.com/txsvc/platform/v2/http"
	"github.com/txsvc/platform/v2/logging"
	"github.com/txsvc/platform/v2/metrics"
	"github.com/txsvc/platform/v2/pkg/account"
	"github.com/txsvc/platform/v2/pkg/timestamp"
	"github.com/txsvc/platform/v2/tasks"
)

type (
	LocalProviderImpl struct {
	}

	LocalLoggingProviderImpl struct {
		lvl logging.Severity
		log *zap.SugaredLogger
	}

	LocalErrorReportingProviderImpl struct {
		log *zap.SugaredLogger
	}
)

var (
	DefaultLoggingConfig        platform.PlatformOpts = platform.WithProvider("platform.default.logger", platform.ProviderTypeLogger, NewLocalLoggingProvider)
	DefaultErrorReportingConfig platform.PlatformOpts = platform.WithProvider("platform.default.errorreporting", platform.ProviderTypeErrorReporter, NewLocalErrorReportingProvider)
	DefaultContextConfig        platform.PlatformOpts = platform.WithProvider("platform.default.context", platform.ProviderTypeHttpContext, NewLocalProvider)
	DefaultTaskConfig           platform.PlatformOpts = platform.WithProvider("platform.default.task", platform.ProviderTypeTask, NewLocalProvider)
	DefaultMetricsConfig        platform.PlatformOpts = platform.WithProvider("platform.default.metrics", platform.ProviderTypeMetrics, NewLocalProvider)
	DefaultAuthConfig           platform.PlatformOpts = platform.WithProvider("platform.default.auth", platform.ProviderTypeAuth, auth.NewAuthorizationProvider)

	errorReportingClient *LocalErrorReportingProviderImpl

	// Interface guards
	_ platform.GenericProvider        = (*LocalProviderImpl)(nil)
	_ http.HttpRequestContextProvider = (*LocalProviderImpl)(nil)
	_ metrics.MetricsProvider         = (*LocalProviderImpl)(nil)
	_ tasks.HttpTaskProvider          = (*LocalProviderImpl)(nil)
	_ auth.AuthorizationProvider      = (*LocalProviderImpl)(nil)

	_ platform.GenericProvider              = (*LocalErrorReportingProviderImpl)(nil)
	_ errorreporting.ErrorReportingProvider = (*LocalErrorReportingProviderImpl)(nil)
	_ platform.GenericProvider              = (*LocalLoggingProviderImpl)(nil)
	_ logging.LoggingProvider               = (*LocalLoggingProviderImpl)(nil)
	_ metrics.MetricsProvider               = (*LocalProviderImpl)(nil)
	_ tasks.HttpTaskProvider                = (*LocalProviderImpl)(nil)
)

func init() {
	callerSkipConf := zap.AddCallerSkip(2)
	l, err := zap.NewProduction(callerSkipConf)

	if err != nil {
		log.Fatal(err)
	}
	er := LocalErrorReportingProviderImpl{
		log: l.Sugar(),
	}

	errorReportingClient = &er
}

func InitLocalProviders() {
	p, err := platform.InitPlatform(context.Background(), DefaultLoggingConfig, DefaultErrorReportingConfig, DefaultContextConfig, DefaultTaskConfig, DefaultMetricsConfig)
	if err != nil {
		log.Fatal(err)
	}
	platform.RegisterPlatform(p)
}

func NewLocalProvider(ID string) interface{} {
	return &LocalProviderImpl{}
}

func (c *LocalProviderImpl) Close() error {
	return nil
}

func (c *LocalProviderImpl) NewHttpContext(req *h.Request) context.Context {
	return context.Background()
}

func NewLocalLoggingProvider(ID string) interface{} {
	callerSkipConf := zap.AddCallerSkip(1)

	l, err := zap.NewProduction(callerSkipConf)
	if err != nil {
		return nil
	}

	logger := LocalLoggingProviderImpl{
		lvl: logging.Info,
		log: l.Sugar(),
	}

	return &logger
}

func (l *LocalLoggingProviderImpl) Close() error {
	return nil
}

func (l *LocalLoggingProviderImpl) Log(msg string, keyValuePairs ...string) {
	l.LogWithLevel(l.lvl, msg, keyValuePairs...)
}

func (l *LocalLoggingProviderImpl) LogWithLevel(lvl logging.Severity, msg string, keyValuePairs ...string) {

	if len(keyValuePairs) > 0 {
		params := make([]interface{}, len(keyValuePairs))
		for i := range keyValuePairs {
			params[i] = keyValuePairs[i]
		}

		switch lvl {
		case logging.Info:
			l.log.Infow(msg, params...)
		case logging.Warn:
			l.log.Warnw(msg, params...)
		case logging.Error:
			l.log.Errorw(msg, params...)
		case logging.Debug:
			l.log.Debugw(msg, params...)
		}
	} else {
		switch lvl {
		case logging.Info:
			l.log.Infow(msg)
		case logging.Warn:
			l.log.Warnw(msg)
		case logging.Error:
			l.log.Errorw(msg)
		case logging.Debug:
			l.log.Debugw(msg)
		}
	}
}

func (er *LocalErrorReportingProviderImpl) Close() error {
	return nil
}

func NewLocalErrorReportingProvider(ID string) interface{} {
	return errorReportingClient
}

func (er *LocalErrorReportingProviderImpl) ReportError(e error) {
	er.log.Error(e)
}

func (m *LocalProviderImpl) Meter(ctx context.Context, metric string, args ...string) {
	// actually does nothing right now
}

func (t *LocalProviderImpl) CreateHttpTask(ctx context.Context, task tasks.HttpTask) error {
	return fmt.Errorf("not implemented")
}

func NewLocalAuthorizationProvider(ID string) interface{} {
	return &LocalProviderImpl{}
}

// SendAccountChallenge sends a notification to the user promting to confirm the account
func (a *LocalProviderImpl) SendAccountChallenge(ctx context.Context, account *account.Account) error {
	return nil
}

// SendAuthToken sends a notification to the user with the current authentication token
func (a *LocalProviderImpl) SendAuthToken(ctx context.Context, account *account.Account) error {
	return nil
}

func (a *LocalProviderImpl) CreateAuthorization(account *account.Account, req *auth.AuthorizationRequest) *auth.Authorization {
	now := timestamp.Now()
	scope := auth.DefaultScope
	if req.Scope != "" {
		scope = req.Scope
	}

	aa := auth.Authorization{
		ClientID:  account.ClientID,
		Realm:     req.Realm,
		Token:     auth.CreateSimpleToken(),
		TokenType: auth.DefaultTokenType,
		UserID:    req.UserID,
		Scope:     scope,
		Revoked:   false,
		Expires:   now + (auth.DefaultAuthorizationExpiration * 86400),
		Created:   now,
		Updated:   now,
	}
	return &aa
}

func (a *LocalProviderImpl) Scope() string {
	return auth.DefaultScope
}

func (a *LocalProviderImpl) Endpoint() string {
	return auth.DefaultEndpoint
}

func (a *LocalProviderImpl) AuthenticationExpiration() int {
	return auth.DefaultAuthenticationExpiration
}

func (a *LocalProviderImpl) AuthorizationExpiration() int {
	return auth.DefaultAuthorizationExpiration
}
