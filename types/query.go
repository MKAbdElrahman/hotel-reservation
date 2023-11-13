package types

type QueryCriteria struct {
	Rating Rating `form:"rating" binding:"required"`
}
