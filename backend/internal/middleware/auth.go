package middleware

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/labstack/echo/v4"
	"guthub.com/imritik7303/boiler-plate-backend/internal/errs"
	"guthub.com/imritik7303/boiler-plate-backend/internal/server"
)

type AuthMiddleware struct {
	server *server.Server
}

func NewAuthMiddleware(s *server.Server) *AuthMiddleware {
	return &AuthMiddleware{
		server: s,
	}
}

func (auth *AuthMiddleware) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return echo.WrapMiddleware(
		//echo.WrapMiddleware(...) converts a net/http style middleware (function with signature func(http.Handler) http.Handler) into an Echo middleware
		clerkhttp.WithHeaderAuthorization(
			clerkhttp.AuthorizationFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				start := time.Now()

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)

				response := map[string]string{
					"code":     "UNAUTHORIZED",
					"message":  "Unauthorized",
					"override": "false",
					"status":   "401",
				}

				if err := json.NewEncoder(w).Encode(response); err != nil {
					auth.server.Logger.Error().Err(err).Str("function", "RequireAuth").Dur(
						"duration", time.Since(start)).Msg("failed to write JSON response")
				} else {
					auth.server.Logger.Error().Str("function", "RequireAuth").Dur("duration", time.Since(start)).Msg(
						"could not get session claims from context")
				}
			}))))(func(c echo.Context) error {
				//"Take the wrapper returned by echo.WrapMiddleware and immediately call it with my custom handler."
		start := time.Now()
		claims, ok := clerk.SessionClaimsFromContext(c.Request().Context())

		if !ok {
			auth.server.Logger.Error().
				Str("function", "RequireAuth").
				Str("request_id", GetRequestID(c)).
				Dur("duration", time.Since(start)).
				Msg("could not get session claims from context")
			return errs.NewUnAuthorizedError("Unauthorized", false)
		}

		c.Set("user_id", claims.Subject)
		c.Set("user_role", claims.ActiveOrganizationRole)
		c.Set("permissions", claims.Claims.ActiveOrganizationPermissions)

		auth.server.Logger.Info().
			Str("function", "RequireAuth").
			Str("user_id", claims.Subject).
			Str("request_id", GetRequestID(c)).
			Dur("duration", time.Since(start)).
			Msg("user authenticated successfully")

		return next(c)
	})
}

/*
RequireAuth is a method on AuthMiddleware that returns an Echo handler (middleware-style function). It uses clerkhttp.WithHeaderAuthorization (a Clerk-provided http middleware) with a custom failure handler (writes a JSON 401). The Clerk middleware, when successful, populates session claims into the request Context. The Echo-side handler then reads those claims, stores useful bits in echo.Context (c.Set(...)) and calls next(c) to continue the chain. If claims are missing, it logs and returns an unauthorized error.
*/


//clerkhhtp.withheaderauth Checks the Authorization header (or other header) in the incoming http.Request.

// Validates the token/session and, on success, attaches the session claims into the request Context.

// It accepts options; here the option is a custom AuthorizationFailureHandler, which tells Clerk what to do if auth fails.