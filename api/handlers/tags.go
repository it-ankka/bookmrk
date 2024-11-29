package handlers

import (
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/template"
)

func TagsViewHandler(registry *template.Registry, app *pocketbase.PocketBase) echo.HandlerFunc {
	return func(c echo.Context) error {
		record, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)

		html, err := registry.LoadFiles(
			"views/layout.html",
			"views/navbar.html",
			"views/main.html",
		).Render(map[string]any{
			"username":      record.Username(),
			"authenticated": true,
		})

		if err != nil {
			return apis.NewNotFoundError("", err)
		}

		return c.HTML(200, html)
	}
}
