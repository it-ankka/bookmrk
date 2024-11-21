package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/it-ankka/bookmrk/api"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tokens"
)

// content holds our static web server content.
//
//go:embed pb_public/*
var content embed.FS

var loginFormHtml = `<form action="/auth" method="POST" style="display: flex; flex-direction: column; gap: 8px; max-width: max-content;"> 
	<input name="email" type="email" placeholder="email" required/>
	<input name="password" type="password" placeholder="password" required/>
	<button type="submit">Log in</button>
</form>`

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// registes a new "GET /*" handler
		e.Router.GET("/*", func(c echo.Context) error {
			record, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
			if record == nil {
				return c.Redirect(303, "/login")
			}
			// TODO Check why is it not authenticating properly?
			return c.HTML(200, "<h1>YOU ARE NOW LOGGED IN "+record.Username()+"</h1>")
			// return c.NoContent(http.StatusNoContent)
		}, apis.ActivityLogger(app), api.LoadAuthContextFromCookie(app))

		e.Router.GET("/login", func(c echo.Context) error {
			if c.QueryParam("error") != "" {
				return c.HTML(400, fmt.Sprintf(loginFormHtml+`<strong style="color: red">%s</strong>`, c.QueryParam("error")))
			}
			return c.HTML(200, loginFormHtml)
		}, apis.ActivityLogger(app), api.LoadAuthContextFromCookie(app))

		// registes a new "POST /auth" handler
		e.Router.POST("/auth", func(c echo.Context) error {
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
				return c.Redirect(303, "/login?error=Invalid login credentials")
			}

			// fetch the user and check the provided password
			record, err := app.Dao().FindAuthRecordByEmail("users", data.Email)
			if err != nil || !record.ValidatePassword(data.Password) {
				if strings.EqualFold(contentType, "Application/json") {
					return apis.NewBadRequestError("Invalid login credentials", err)
				}
				return c.Redirect(303, "/login?error=Invalid login credentials")
			}

			// generate a new auth token for the found user record
			token, err := tokens.NewRecordAuthToken(app, record)
			if err != nil {
				if strings.EqualFold(contentType, "Application/json") {
					return apis.NewBadRequestError("Failed to create auth token", err)
				}
				return c.Redirect(303, "/login?error=Invalid login credentials")
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
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
