package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/it-ankka/bookmrk/server/bookmark"
	"github.com/it-ankka/bookmrk/server/config"
	"github.com/it-ankka/bookmrk/server/middleware"
	"github.com/it-ankka/bookmrk/server/utils"
)

type bookmarksJson struct{ bookmarks []bookmark.Bookmark }

func readBookmarks() bookmarksJson {
	jsonFile, err := os.ReadFile(config.Settings.JsonFile)
	var data bookmarksJson = bookmarksJson{bookmarks: []bookmark.Bookmark{}}
	if os.IsNotExist(err) {
		jsonFile, err := json.Marshal(&data)
		utils.Check(err)
		err = os.WriteFile(config.Settings.JsonFile, jsonFile, 0644)
		utils.Check(err)
		return data
	}
	err = json.Unmarshal(jsonFile, &data)
	utils.Check(err)
	return data
}

func main() {

	router := http.NewServeMux()
	logger := slog.Default()
	v1 := http.NewServeMux()
	bookmarkHandler := bookmark.Handler{Bookmarks: []bookmark.Bookmark{}}

	v1.HandleFunc("GET /bookmarks", bookmarkHandler.GetAll)
	v1.HandleFunc("POST /bookmarks", bookmarkHandler.Create)
	v1.HandleFunc("GET /bookmarks/{id}", bookmarkHandler.GetById)
	v1.HandleFunc("PATCH /bookmarks/{id}", bookmarkHandler.UpdateById)
	v1.HandleFunc("DELETE /bookmarks/{id}", bookmarkHandler.DeleteById)

	router.Handle("/v1/", http.StripPrefix("/v1", v1))

	logger.Info("Server started. Listening on " + config.Settings.Address)
	http.ListenAndServe(config.Settings.Address, middleware.Logger(router, logger))
}
