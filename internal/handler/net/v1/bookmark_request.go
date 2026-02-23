package v1

type CreateBookmarkRequest struct {
	Title string `json:"title" validate:"required"`
	Value string `json:"value" validate:"required"`
}
