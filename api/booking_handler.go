package api

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mkabdelrahman/hotel-reservation/business"
	"github.com/mkabdelrahman/hotel-reservation/errorlog"
	"github.com/mkabdelrahman/hotel-reservation/types"
)

type BookingHandler struct {
	Manager              *business.Manager
	ErrorResponseHandler *errorlog.HTTPErrorResponseWriterAndLogger
}

func NewBookingHandler(m *business.Manager, errorLogger *log.Logger) *BookingHandler {
	return &BookingHandler{
		Manager:              m,
		ErrorResponseHandler: &errorlog.HTTPErrorResponseWriterAndLogger{Logger: errorLogger},
	}
}

func (h *BookingHandler) HandlePostBooking(ctx *gin.Context) {
	var params types.NewBookingParams

	err := ctx.ShouldBindJSON(&params)

	if err != nil {
		appErr := errorlog.BadRequestError(err)
		h.ErrorResponseHandler.LogAndHandleError(ctx.Writer, ctx.Writer, appErr)
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		appErr := errorlog.InternalServerError(errors.New("userID not found in context"))
		h.ErrorResponseHandler.LogAndHandleError(ctx.Writer, ctx.Writer, appErr)
		return
	}
	params.UserID = userID.(string)

	err = params.Validate()
	if err != nil {
		appErr := errorlog.BadRequestError(err)
		h.ErrorResponseHandler.LogAndHandleError(ctx.Writer, ctx.Writer, appErr)
		return
	}

	insertedBooking, err := h.Manager.AddNewBooking(ctx, params)

	if err != nil {
		appErr := errorlog.InternalServerError(err)
		h.ErrorResponseHandler.LogAndHandleError(ctx.Writer, ctx.Writer, appErr)
		return
	}

	ctx.JSON(http.StatusOK, insertedBooking)
}

func (h *BookingHandler) HandleGetBooking(ctx *gin.Context) {
	id := ctx.Param("id")

	booking, err := h.Manager.BookingStore.GetBookingByID(ctx, id)

	if err != nil {
		appErr := errorlog.InternalServerError(err)
		h.ErrorResponseHandler.LogAndHandleError(ctx.Writer, ctx.Writer, appErr)
		return
	}

	ctx.JSON(http.StatusOK, booking)
}

func (h *BookingHandler) HandleGetBookings(ctx *gin.Context) {
	bookings, err := h.Manager.ListBookings(ctx)

	if err != nil {
		appErr := errorlog.InternalServerError(err)
		h.ErrorResponseHandler.LogAndHandleError(ctx.Writer, ctx.Writer, appErr)
		return
	}

	ctx.JSON(http.StatusOK, bookings)
}

func (h *BookingHandler) HandleCancelBooking(ctx *gin.Context) {
	id := ctx.Param("id")

	err := h.Manager.CancelBooking(ctx, id)

	if err != nil {
		appErr := errorlog.InternalServerError(err)
		h.ErrorResponseHandler.LogAndHandleError(ctx.Writer, ctx.Writer, appErr)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Booking canceled successfully"})
}
