package category

type Category struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	ChildId int    `json:"child_id"`
}

type CreateUpdateCategory struct {
	Title   string `json:"title"`
	ChildId int    `json:"child_id"`
}
