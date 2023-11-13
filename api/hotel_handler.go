package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mkabdelrahman/hotel-reservation/business"
	"github.com/mkabdelrahman/hotel-reservation/types"
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

func (h *HotelHandler) HandleHotelSearch(ctx *gin.Context) {
	var q types.QueryCriteria

	if err := ctx.ShouldBindQuery(&q); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hotels, err := h.Manager.QueryHotels(ctx, q)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Example response
	ctx.JSON(http.StatusOK, hotels)
}

func (h *HotelHandler) HandleGetHotelRooms(ctx *gin.Context) {

	id := ctx.Param("id")

	hotel, err := h.Manager.HotelStore.GetHotel(ctx, id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var rooms []*types.Room

	for _, roomID := range hotel.Rooms {
		room, err := h.Manager.RoomStore.GetRoom(ctx, roomID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rooms = append(rooms, room)

	}
	ctx.JSON(http.StatusOK, rooms)

}
