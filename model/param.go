package model
type Params struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
