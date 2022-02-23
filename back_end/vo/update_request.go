package vo

type CreateUpdateRequest struct {
	Name        string `json:"name" binding:"required"`
	NewName     string `json:"new_name"`
	NewColor    string `json:"new_color"`
	NewBreed    string `json:"new_breed"`
	NewPlace    string `json:"new_place"`
	NewSpecific string `json:"new_specific"`
}
