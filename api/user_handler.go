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

func (h *UserHandler) HandleGetUsers(ctx *gin.Context) {

	users, err := h.store.GetUsers(ctx)
	if err != nil {
		err = types.HTTPError{
			Description: err.Error(),
			StatusCode:  http.StatusInternalServerError,
			Metadata:    "",
		}
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (h *UserHandler) HandlePostUser(ctx *gin.Context) {
	var params types.CreateUserParams

	err := ctx.ShouldBindJSON(&params)

	if err != nil {
		err = types.HTTPError{
			Description: err.Error(),
			StatusCode:  http.StatusInternalServerError,
			Metadata:    "",
		}
		ctx.Error(err)
		return
	}

	errs := params.Validate()
	if len(params.Validate()) > 0 {
		for _, err := range errs {
			err = types.HTTPError{
				Description: err.Error(),
				StatusCode:  http.StatusInternalServerError,
				Metadata:    "",
			}
			ctx.Error(err)
		}
		return
	}
	user, err := types.CreateNewUserFromParams(params)
	if err != nil {
		err = types.HTTPError{
			Description: err.Error(),
			StatusCode:  http.StatusInternalServerError,
			Metadata:    "",
		}
		ctx.Error(err)
		return
	}

	insertedUser, err := h.store.InsertUser(ctx, user)
	if err != nil {
		err = types.HTTPError{
			Description: err.Error(),
			StatusCode:  http.StatusInternalServerError,
			Metadata:    "",
		}
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, insertedUser)
}
