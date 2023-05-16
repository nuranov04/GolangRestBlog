package post

type CreatePostDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	OwnerId     int    `json:"owner_id"`
}

type UpdatePostDTO struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}
