package post

type Post struct {
	ID          int    `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	OwnerId     int    `json:"owner_id,omitempty"`
}
