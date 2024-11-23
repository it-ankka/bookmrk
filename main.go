package main

import (
	"embed"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/it-ankka/bookmrk/api"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/pocketbase/pocketbase/tokens"
	"github.com/pocketbase/pocketbase/tools/template"

	_ "github.com/it-ankka/bookmrk/migrations"
)

// content holds our static web server content.

//go:embed views static migrations
var content embed.FS

func main() {
	app := pocketbase.New()

	// loosely check if it was executed using "go run"
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Admin UI
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// this is safe to be used by multiple goroutines
		// (it acts as store for the parsed templates)
		registry := template.NewRegistry()

		e.Router.Static("/static", "./static")

		// registes a new "GET /*" handler
		e.Router.GET("/*", func(c echo.Context) error {
			record, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
			if record == nil {
				return c.Redirect(303, "/login")
			}

			html, err := registry.LoadFiles(
				"views/layout.html",
				"views/navbar.html",
				"views/default.html",
			).Render(map[string]any{
				"name":          record.Username(),
				"authenticated": true,
			})

			if err != nil {
				return apis.NewNotFoundError("", err)
			}

			// TODO Check why is it not authenticating properly?
			return c.HTML(200, html)
			// return c.NoContent(http.StatusNoContent)
		}, apis.ActivityLogger(app), api.LoadAuthContextFromCookie(app))

		e.Router.GET("/login", func(c echo.Context) error {
			record, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
			if record != nil {
				return c.Redirect(303, "/")
			}

			html, err := registry.LoadFiles(
				"views/layout.html",
				"views/login.html",
			).Render(map[string]any{
				"error": c.QueryParam("error"),
			})

			if err != nil {
				return apis.NewNotFoundError("", err)
			}

			if c.QueryParam("error") != "" {
				return c.HTML(400, html)
			}

			return c.HTML(200, html)
		}, apis.ActivityLogger(app), api.LoadAuthContextFromCookie(app))

		e.Router.POST("/login", func(c echo.Context) error {
			data := &struct {
				Email    string `form:"email" json:"email"`
				Password string `form:"password" json:"password"`
			}{}

			contentType := c.Request().Header.Get("Content-Type")

			// read the request data
			if err := c.Bind(data); err != nil {
				if strings.EqualFold(contentType, "Application/json") {
					return apis.NewBadRequestError("Failed to read request data", err)
				}
				return c.Redirect(303, "/login?error=1")
			}

			// fetch the user and check the provided password
			record, err := app.Dao().FindAuthRecordByEmail("users", data.Email)
			if err != nil || !record.ValidatePassword(data.Password) {
				if strings.EqualFold(contentType, "Application/json") {
					return apis.NewBadRequestError("Invalid login credentials", err)
				}
				return c.Redirect(303, "/login?error=1")
			}

			// generate a new auth token for the found user record
			token, err := tokens.NewRecordAuthToken(app, record)
			if err != nil {
				if strings.EqualFold(contentType, "Application/json") {
					return apis.NewBadRequestError("Failed to create auth token", err)
				}
				return c.Redirect(303, "/login?error=1")
			}

			// set it as cookie
			cookie := &http.Cookie{
				Name:     "pb_auth", // rename with the name of your cookie
				Value:    token,
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
				// you can use the token duration or any other expiration (eg. only 1 day)
				Expires: time.Now().Add(time.Duration(app.Settings().RecordAuthToken.Duration) * time.Second),
			}

			c.SetCookie(cookie)

			if strings.EqualFold(contentType, "Application/json") {
				return c.NoContent(http.StatusNoContent)
			}
			return c.Redirect(303, "/")
		}, apis.ActivityLogger(app))

		e.Router.POST("/logout", func(c echo.Context) error {
			// Clear auth cookie
			c.SetCookie(&http.Cookie{
				Name:     "pb_auth",
				Value:    "",
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
				MaxAge:   -1,
				Expires:  time.Unix(0, 0),
			})

			return c.Redirect(303, "/login")
		}, apis.ActivityLogger(app))

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
