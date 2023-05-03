package post

import "go.mod/internal/user"

type Post struct {
	ID          string    `json:"id,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	OwnerId     string    `json:"owner_id,omitempty"`
	Owner       user.User `json:"owner"`
}
