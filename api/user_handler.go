package api

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mkabdelrahman/hotel-reservation/business"
	"github.com/mkabdelrahman/hotel-reservation/errorlog"
	"github.com/mkabdelrahman/hotel-reservation/types"
)

type UserHandler struct {
	Manager              *business.Manager
	ErrorResponseHandler *errorlog.HTTPErrorResponseWriterAndLogger
}

func NewUserHandler(m *business.Manager, errorLogger *log.Logger) *UserHandler {
	return &UserHandler{
		Manager:              m,
		ErrorResponseHandler: &errorlog.HTTPErrorResponseWriterAndLogger{Logger: errorLogger},
	}
}

func (h *UserHandler) HandleGetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := h.Manager.GetUserByID(ctx, id)

	if err != nil {
		httpError := errorlog.InternalServerError(err)
		h.ErrorResponseHandler.LogAndHandleError(os.Stdout, ctx.Writer, httpError)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *UserHandler) HandleGetUserBookings(ctx *gin.Context) {
	id := ctx.Param("id")

	// Make sure the user exists
	user, err := h.Manager.GetUserByID(ctx, id)
	if err != nil {
		httpError := errorlog.InternalServerError(err)
		h.ErrorResponseHandler.LogAndHandleError(os.Stdout, ctx.Writer, httpError)
		return
	}

	if user == nil {
		httpError := errorlog.NotFoundError(errors.New("user not found"))
		h.ErrorResponseHandler.LogAndHandleError(os.Stdout, ctx.Writer, httpError)
		return
	}

	// Get the user's bookings
	bookings, err := h.Manager.ListUserBookings(ctx, id)
	if err != nil {
		httpError := errorlog.InternalServerError(err)
		h.ErrorResponseHandler.LogAndHandleError(os.Stdout, ctx.Writer, httpError)
		return
	}

	ctx.JSON(http.StatusOK, bookings)
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

	users, err := h.Manager.ListUsers(ctx, filter)
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

	updatedUser, err := h.Manager.UpdateUser(ctx, userID, params)
	if err != nil {
		httpError := errorlog.InternalServerError(err)
		h.ErrorResponseHandler.LogAndHandleError(os.Stdout, ctx.Writer, httpError)
		return
	}

	ctx.JSON(http.StatusOK, updatedUser)
}

func (h *UserHandler) HandleDeleteUser(ctx *gin.Context) {
	userID := ctx.Param("id")

	if err := h.Manager.DeleteUser(ctx, userID); err != nil {
		httpError := errorlog.InternalServerError(err)
		h.ErrorResponseHandler.LogAndHandleError(os.Stdout, ctx.Writer, httpError)
		return

	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "user has been deleted"})
}

func (h *UserHandler) HandlePostUser(ctx *gin.Context) {
	var params types.NewUserParams

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

	insertedUser, err := h.Manager.AddNewUser(ctx, params)
	if err != nil {
		httpError := errorlog.InternalServerError(err)
		h.ErrorResponseHandler.LogAndHandleError(os.Stdout, ctx.Writer, httpError)
		return
	}

	ctx.JSON(http.StatusOK, insertedUser)
}
