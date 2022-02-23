package vo

type DeleteRequest struct {
	CatName string `json:"cat_name" binding:"required"`
}
