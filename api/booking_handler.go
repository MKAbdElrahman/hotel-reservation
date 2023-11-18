package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mkabdelrahman/hotel-reservation/business"
	"github.com/mkabdelrahman/hotel-reservation/types"
)

type BookingHandler struct {
	Manager *business.Manager
}

func NewBookingHandler(m *business.Manager) *BookingHandler {
	return &BookingHandler{
		Manager: m,
	}
}

func (h *BookingHandler) HandlePostBooking(ctx *gin.Context) {
	var params types.NewBookingParams

	err := ctx.ShouldBindJSON(&params)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	booking := types.NewBookingFromParams(params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "userID not found in context"})
		return
	}
	booking.UserID = userID.(string)

	insertedBooking, err := h.Manager.BookingStore.InsertBooking(ctx, booking)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, insertedBooking)
}

func (h *BookingHandler) HandleGetBooking(ctx *gin.Context) {

	id := ctx.Param("id")

	booking, err := h.Manager.BookingStore.GetBookingByID(ctx, id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, booking)
}
