package post

import "go.mod/internal/user"

type CreatePostDTO struct {
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Owner       user.User `json:"owner"`
}
