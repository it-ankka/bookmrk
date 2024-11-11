package bookmark

import "github.com/google/uuid"

type Bookmark struct {
	Id          uuid.UUID `json:"id"`
	Url         string    `json:"url"`
	Tags        []string  `json:"tags"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}
