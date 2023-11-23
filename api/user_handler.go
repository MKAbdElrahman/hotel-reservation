package api

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mkabdelrahman/hotel-reservation/db"
	"github.com/mkabdelrahman/hotel-reservation/errorlog"
	"github.com/mkabdelrahman/hotel-reservation/types"
)

type UserHandler struct {
	store                db.UserStore
	ErrorResponseHandler *errorlog.HTTPErrorResponseWriterAndLogger
}

func NewUserHandler(store db.UserStore, errorLogger *log.Logger) *UserHandler {
	return &UserHandler{
		store:                store,
		ErrorResponseHandler: &errorlog.HTTPErrorResponseWriterAndLogger{Logger: errorLogger},
	}
}

func (h *UserHandler) HandleGetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := h.store.GetUserByID(ctx, id)

	if err != nil {
		httpError := errorlog.InternalServerError(err)
		h.ErrorResponseHandler.LogAndHandleError(os.Stdout, ctx.Writer, httpError)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *UserHandler) HandleGetUsers(ctx *gin.Context) {

	filter := types.NewUsersPaginationFilter()

	if err := ctx.ShouldBindQuery(&filter); err != nil {
		httpError := errorlog.BadRequestError(err)
		h.ErrorResponseHandler.LogAndHandleError(os.Stdout, ctx.Writer, httpError)
		return
	}

	if err := filter.Validate(); err != nil {
		httpError := errorlog.BadRequestError(err)
		h.ErrorResponseHandler.LogAndHandleError(os.Stdout, ctx.Writer, httpError)
		return
	}

	users, err := h.store.GetUsersWithPagination(ctx, filter)
	if err != nil {
		httpError := errorlog.InternalServerError(err)
		h.ErrorResponseHandler.LogAndHandleError(os.Stdout, ctx.Writer, httpError)
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (h *UserHandler) HandleUpdateUser(ctx *gin.Context) {
	var params types.UpdateUserParams

	if err := ctx.ShouldBindJSON(&params); err != nil {
		httpError := errorlog.BadRequestError(err)
		h.ErrorResponseHandler.LogAndHandleError(os.Stdout, ctx.Writer, httpError)
		return
	}

	userID := ctx.Param("id")

	updatedUser, err := h.store.UpdateUser(ctx, userID, params)
	if err != nil {
		httpError := errorlog.InternalServerError(err)
		h.ErrorResponseHandler.LogAndHandleError(os.Stdout, ctx.Writer, httpError)
		return
	}

	ctx.JSON(http.StatusOK, updatedUser)
}

func (h *UserHandler) HandleDeleteUser(ctx *gin.Context) {
	userID := ctx.Param("id")

	if err := h.store.DeleteUser(ctx, userID); err != nil {
		httpError := errorlog.InternalServerError(err)
		h.ErrorResponseHandler.LogAndHandleError(os.Stdout, ctx.Writer, httpError)
		return

	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "user has been deleted"})
}

func (h *UserHandler) HandlePostUser(ctx *gin.Context) {
	var params types.UserParams

	err := ctx.ShouldBindJSON(&params)

	if err != nil {
		httpError := errorlog.BadRequestError(err)
		h.ErrorResponseHandler.LogAndHandleError(os.Stdout, ctx.Writer, httpError)
		return
	}

	errs := params.Validate()
	if len(params.Validate()) > 0 {
		for _, err := range errs {
			httpError := errorlog.BadRequestError(err)
			h.ErrorResponseHandler.LogAndHandleError(os.Stdout, ctx.Writer, httpError)
		}
		return
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		httpError := errorlog.BadRequestError(err)
		h.ErrorResponseHandler.LogAndHandleError(os.Stdout, ctx.Writer, httpError)
		return
	}

	insertedUser, err := h.store.InsertUser(ctx, user)
	if err != nil {
		httpError := errorlog.InternalServerError(err)
		h.ErrorResponseHandler.LogAndHandleError(os.Stdout, ctx.Writer, httpError)
		return
	}

	ctx.JSON(http.StatusOK, insertedUser)
}
