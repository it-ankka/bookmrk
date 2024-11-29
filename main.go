package main

import (
	"embed"
	"log"
	"os"
	"strings"

	"github.com/it-ankka/bookmrk/api/handlers"
	"github.com/it-ankka/bookmrk/api/middlewares"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/pocketbase/pocketbase/tools/template"

	_ "github.com/it-ankka/bookmrk/migrations"
)

// content holds our static web server content.
//
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

		// Root page redirects to bookmarks page
		e.Router.GET("/", func(c echo.Context) error {
			return c.Redirect(303, "/bookmarks")
		}, apis.ActivityLogger(app), middlewares.RequireAuth(app))

		e.Router.GET("/*",
			handlers.NotFoundViewHandler(registry, app),
			apis.ActivityLogger(app),
			middlewares.RequireAuth(app))

		e.Router.GET("/bookmarks",
			handlers.BookmarksViewHandler(registry, app),
			apis.ActivityLogger(app),
			middlewares.RequireAuth(app))

		e.Router.POST("/bookmarks",
			handlers.BookmarksPostHandler(app),
			apis.ActivityLogger(app),
			middlewares.RequireAuth(app))

		e.Router.GET("/tags",
			handlers.TagsViewHandler(registry, app),
			apis.ActivityLogger(app),
			middlewares.RequireAuth(app))

		e.Router.GET("/login",
			handlers.LoginViewHandler(registry, app),
			apis.ActivityLogger(app),
			middlewares.LoadAuthContextFromCookie(app))

		e.Router.POST("/login",
			handlers.LoginPostHandler(app),
			apis.ActivityLogger(app))

		e.Router.POST("/logout",
			handlers.LogoutPostHandler(app),
			apis.ActivityLogger(app))

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
