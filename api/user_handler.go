package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mkabdelrahman/hotel-reservation/types"
)

func HandleGetUsers(ctx *gin.Context) {

	ctx.JSON(200, "mk")
}

func HandleGetUserById(ctx *gin.Context) {
	u := types.User{
		FirstName: "Mohamed",
		LastName:  "Kamal",
	}
	ctx.JSON(200, u)
}
