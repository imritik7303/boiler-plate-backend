package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/newrelic/go-agent/v3/integrations/nrecho-v4"
	"github.com/newrelic/go-agent/v3/integrations/nrpkgerrors"
	"github.com/newrelic/go-agent/v3/newrelic"
	"guthub.com/imritik7303/boiler-plate-backend/internal/server"
)

type TracingMiddleware struct {
	server *server.Server
	nrApp  *newrelic.Application
}

func NewTracingMiddleware (s *server.Server , nrApp *newrelic.Application) *TracingMiddleware {
    return &TracingMiddleware{
		server: s,
		nrApp: nrApp,
	}
}

//new relic middleware return the relic midlleware for echo
func (tm *TracingMiddleware) NewRelicMiddleware() echo.MiddlewareFunc {
	if tm.nrApp == nil {
		//return a no-op muddleware if new relic not initialized 
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return next
		}
	}
	return nrecho.Middleware(tm.nrApp)
}
//enhance tracing add custom attributes to new relic transactions

func(tm *TracingMiddleware) EnhanceTracing() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			//get new relic transaction from context
			txn := newrelic.FromContext(c.Request().Context())
			if txn == nil {
				return next(c)
			}

			// service.name and service.environment are already set in logger and New Relic config
			txn.AddAttribute("http.real_ip", c.RealIP())
			txn.AddAttribute("http.user_agent", c.Request().UserAgent())

			// Add request ID if available
			if requestID := GetRequestID(c); requestID != "" {
				txn.AddAttribute("request.id", requestID)
			}

			// Add user context if available
			if userID := c.Get("user_id"); userID != nil {
				if userIDStr, ok := userID.(string); ok {
					txn.AddAttribute("user.id", userIDStr)
				}
			}

			// Execute next handler
			err := next(c)
			// Record error if any with enhanced stack traces
			if err != nil {
				txn.NoticeError(nrpkgerrors.Wrap(err))
			}

			// Add response status
			txn.AddAttribute("http.status_code", c.Response().Status)

			return err
		}
	}
}