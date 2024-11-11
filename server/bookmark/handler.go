package bookmark

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/it-ankka/bookmrk/server/config"
	"github.com/it-ankka/bookmrk/server/utils"
)

type Handler struct{ Bookmarks []Bookmark }

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Result []Bookmark `json:"result"`
	}{
		Result: h.Bookmarks,
	})
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	for _, bookmark := range h.Bookmarks {
		if bookmark.Id.String() == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(struct {
				Id string `json:"id"`
			}{
				Id: id,
			})
			return
		}
	}
	w.WriteHeader(404)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	newBookmark := &Bookmark{}
	err := json.NewDecoder(r.Body).Decode(newBookmark)
	id, err := uuid.NewUUID()
	if err != nil {
		w.WriteHeader(400)
		return
	}

	newBookmark.Id = id
	slog.Info(fmt.Sprintf("New bookmark: %+v\n", newBookmark))
	h.Bookmarks = append(h.Bookmarks, *newBookmark)

	jsonFile, err := json.Marshal(h.Bookmarks)
	if err != nil {
		slog.Error(err.Error())
	}
	err = os.WriteFile(config.Settings.JsonFile, jsonFile, 0644)
	if err != nil {
		slog.Error(err.Error())
	}
	utils.Check(err)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newBookmark)
}

// TODO
func (h *Handler) UpdateById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Id string `json:"id"`
	}{
		Id: id,
	})
}

// TODO
func (h *Handler) DeleteById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Id string `json:"id"`
	}{
		Id: id,
	})
}
