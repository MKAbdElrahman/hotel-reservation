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

func (h *UserHandler) HandleGetUser(ctx *gin.Context) {

	id := ctx.Param("id")

	c := context.Background()

	user, err := h.store.GetUserByID(c, id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *UserHandler) HandleGetUsers(ctx *gin.Context) {

	users, err := h.store.GetUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (h *UserHandler) HandleUpdateUser(ctx *gin.Context) {
	var params types.UpdateUserParams

	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// errs := params.Validate()
	// if len(params.Validate()) > 0 {
	// 	for _, err := range errs {
	// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	}
	// 	return
	// }

	userID := ctx.Param("id")

	updatedUser, err := h.store.UpdateUser(ctx, userID, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedUser)
}

func (h *UserHandler) HandleDeleteUser(ctx *gin.Context) {
	userID := ctx.Param("id")

	if err := h.store.DeleteUser(ctx, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return

	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "user has been deleted"})
}

func (h *UserHandler) HandlePostUser(ctx *gin.Context) {
	var params types.CreateUserParams

	err := ctx.ShouldBindJSON(&params)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	errs := params.Validate()
	if len(params.Validate()) > 0 {
		for _, err := range errs {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	user, err := types.CreateNewUserFromParams(params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	insertedUser, err := h.store.InsertUser(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, insertedUser)
}
