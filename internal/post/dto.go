package post

type CreatePostDTO struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	OwnerId     int    `json:"owner_id"`
}

type UpdatePostDTO struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}
