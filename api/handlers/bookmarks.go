package handlers

import (
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/template"
)

func BookmarksViewHandler(registry *template.Registry, app *pocketbase.PocketBase) echo.HandlerFunc {
	return func(c echo.Context) error {
		record, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
		q := c.QueryParam("q")
		bookmarks, err := app.Dao().FindRecordsByFilter("bookmarks", "user.id = {:user_id}", "-created", -1, 0, dbx.Params{"user_id": record.Id})

		if err != nil {
			println(err.Error())
		}

		matchingBookmarks := []*models.Record{}

		if len(q) > 0 {
			for _, bookmark := range bookmarks {
				bm := bookmark.PublicExport()
				fields := []string{"url", "name", "description"}
				for _, field := range fields {
					if strings.Contains(bm[field].(string), q) {
						matchingBookmarks = append(matchingBookmarks, bookmark)
					}
				}
			}
		} else {
			matchingBookmarks = bookmarks
		}

		html, err := registry.LoadFiles(
			"views/layout.html",
			"views/navbar.html",
			"views/bookmarks.html",
			"views/bookmark.html",
		).Render(map[string]any{
			"username":      record.Username(),
			"authenticated": true,
			"bookmarks":     matchingBookmarks,
			"query":         q,
		})

		if err != nil {
			println(err.Error())
			return apis.NewNotFoundError("", err)
		}

		return c.HTML(200, html)
	}
}

func BookmarksPostHandler(app *pocketbase.PocketBase) echo.HandlerFunc {
	return func(c echo.Context) error {
		userRecord, _ := c.Get(apis.ContextAuthRecordKey).(*models.Record)
		data := &struct {
			Id          string `form:"id" json:"id"`
			Url         string `form:"url" json:"url"`
			Name        string `form:"name" json:"name"`
			Description string `form:"description" json:"description"`
		}{}

		contentType := c.Request().Header.Get("Content-Type")

		// read the request data
		if err := c.Bind(data); err != nil {
			if strings.EqualFold(contentType, "Application/json") {
				return apis.NewBadRequestError("Failed to read request data", err)
			}
			return c.Redirect(303, "/bookmarks?error=1")
		}

		var record *models.Record
		var err error

		if len(data.Id) > 0 {
			// Update existing record
			record, err = app.Dao().FindRecordById("bookmarks", data.Id)
			if err != nil {
				app.Logger().Error(err.Error())
				if strings.EqualFold(contentType, "Application/json") {
					return apis.NewNotFoundError("Bookmark not found", nil)
				}
				return c.Redirect(303, "/bookmarks?error=1")
			}
		} else {
			// Create new record
			collection, err := app.Dao().FindCollectionByNameOrId("bookmarks")
			if err != nil {
				app.Logger().Error(err.Error())
				if strings.EqualFold(contentType, "Application/json") {
					return apis.NewApiError(500, "Failed to access bookmarks collection", nil)
				}
				return c.Redirect(303, "/bookmarks?error=1")
			}
			record = models.NewRecord(collection)
			record.IsNew()
			record.Set("user", userRecord.Id)
		}

		record.Set("url", data.Url)
		record.Set("name", data.Name)
		record.Set("description", data.Description)

		err = app.Dao().Save(record)

		if err != nil {
			app.Logger().Error(err.Error())
			if strings.EqualFold(contentType, "Application/json") {
				return apis.NewApiError(500, "Operation failed", nil)
			}
			return c.Redirect(303, "/bookmarks?error=1")
		}

		return c.Redirect(303, "/bookmarks?error=0")
	}
}
