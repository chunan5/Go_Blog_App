package dto

type ListPostDto struct {
	Title      *string `json:"title"`
	AuthorName *string `json:"author"`
	CreatedAt  *string `json:"createdAt"`
}
