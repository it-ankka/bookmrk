package handlers

import (
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/template"
)

func NotFoundViewHandler(registry *template.Registry, app *pocketbase.PocketBase) echo.HandlerFunc {
	return func(c echo.Context) error {
		record, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)

		html, _ := registry.LoadFiles(
			"views/layout.html",
			"views/404.html",
		).Render(map[string]any{
			"name":          record.Username(),
			"authenticated": true,
		})
		return c.HTML(404, html)
	}
}
