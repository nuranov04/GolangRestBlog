package model

type CreatePostDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	OwnerId     int    `json:"owner_id"`
}

type UpdatePostDTO struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type Post struct {
	ID          int    `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	OwnerId     int    `json:"owner_id,omitempty"`
}
