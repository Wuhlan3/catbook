package vo

type CreatePostRequest struct {
	Name string `json:"name" binding:"required"`
	Color string `json:"color"`
	Breed string `json:"breed"`
	Place string `json:"place"`
	Specific string `json:"specific"`
	ImgUrl string 	`json:"img_url"`
}
