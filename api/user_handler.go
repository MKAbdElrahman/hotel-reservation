package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mkabdelrahman/hotel-reservation/db"
	"github.com/mkabdelrahman/hotel-reservation/types"
)

type UserHandler struct {
	store db.UserStore
}

func NewUserHandler(store db.UserStore) *UserHandler {
	return &UserHandler{
		store: store,
	}
}

func (h *UserHandler) HandleGetUserByID(ctx *gin.Context) {

	id := ctx.Param("id")

	c := context.Background()

	user, err := h.store.GetUserByID(c, id)

	if err != nil {

		err := types.HTTPError{
			Description: err.Error(),
			StatusCode:  http.StatusInternalServerError,
			Metadata:    "",
		}
		ctx.Error(err)
		return
	}

	ctx.JSON(200, user)
}

func HandleGetUsers(ctx *gin.Context) {

	ctx.JSON(200, "mk")
}
