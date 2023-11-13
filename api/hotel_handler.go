package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mkabdelrahman/hotel-reservation/business"
)

type HotelHandler struct {
	Manager *business.Manager
}

func NewHotelHandler(m *business.Manager) *HotelHandler {
	return &HotelHandler{
		Manager: m,
	}
}

func (h *HotelHandler) HandleGetHotels(ctx *gin.Context) {
	hotels, err := h.Manager.ListHotels(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, hotels)
}

func (h *HotelHandler) HandleGetHotel(ctx *gin.Context) {

	id := ctx.Param("id")

	hotel, err := h.Manager.HotelStore.GetHotel(ctx, id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, hotel)
}
