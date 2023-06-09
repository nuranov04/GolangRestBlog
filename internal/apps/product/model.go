package product

type CreateProductDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	OwnerId     int    `json:"owner_id"`
	CategoryId  int    `json:"category_id"`
}

type UpdateProductDTO struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type Product struct {
	ID          int    `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	OwnerId     int    `json:"owner_id,omitempty"`
	CategoryId  int    `json:"category_id,omitempty"`
}
