package middlewares

import (
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tokens"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/spf13/cast"
)

// LoadAuthContext middleware reads the pb_auth cookie
// and loads the token related record or admin instance into the
// request's context.
func LoadAuthContextFromCookie(app core.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := c.Cookie("pb_auth")
			if err != nil {
				app.Logger().Error(err.Error())
				return next(c)
			}
			if token == nil || token.Value == "" {
				return next(c)
			}

			claims, _ := security.ParseUnverifiedJWT(token.Value)
			tokenType := cast.ToString(claims["type"])

			switch tokenType {
			case tokens.TypeAdmin:
				admin, err := app.Dao().FindAdminByToken(
					token.Value,
					app.Settings().AdminAuthToken.Secret,
				)
				if err == nil && admin != nil {
					c.Set(apis.ContextAdminKey, admin)
				}
			case tokens.TypeAuthRecord:
				record, err := app.Dao().FindAuthRecordByToken(
					token.Value,
					app.Settings().RecordAuthToken.Secret,
				)
				if err == nil && record != nil {
					c.Set(apis.ContextAuthRecordKey, record)
				}
			}

			return next(c)
		}
	}
}

func RequireAuth(app core.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return LoadAuthContextFromCookie(app)(func(c echo.Context) error {
			record, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
			if record == nil {
				return c.Redirect(303, "/login")
			}
			return next(c)
		})
	}
}
