package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mkabdelrahman/hotel-reservation/business"
	"github.com/mkabdelrahman/hotel-reservation/errorlog"
	"github.com/mkabdelrahman/hotel-reservation/types"
)

type HotelHandler struct {
	Manager              *business.Manager
	ErrorResponseHandler *errorlog.HTTPErrorResponseWriterAndLogger
}

func NewHotelHandler(m *business.Manager, errorLogger *log.Logger) *HotelHandler {
	return &HotelHandler{
		Manager:              m,
		ErrorResponseHandler: &errorlog.HTTPErrorResponseWriterAndLogger{Logger: errorLogger},
	}
}

func (h *HotelHandler) HandleGetHotels(ctx *gin.Context) {
	hotels, err := h.Manager.ListHotels(ctx)

	if err != nil {
		appErr := errorlog.InternalServerError(err)
		h.ErrorResponseHandler.LogAndHandleError(ctx.Writer, ctx.Writer, appErr)
		return
	}

	ctx.JSON(http.StatusOK, hotels)
}

func (h *HotelHandler) HandleGetHotel(ctx *gin.Context) {
	id := ctx.Param("id")

	hotel, err := h.Manager.HotelStore.GetHotel(ctx, id)

	if err != nil {
		appErr := errorlog.InternalServerError(err)
		h.ErrorResponseHandler.LogAndHandleError(ctx.Writer, ctx.Writer, appErr)
		return
	}

	ctx.JSON(http.StatusOK, hotel)
}

func (h *HotelHandler) HandleHotelSearch(ctx *gin.Context) {
	var q types.QueryCriteria

	if err := ctx.ShouldBindQuery(&q); err != nil {
		appErr := errorlog.BadRequestError(err)
		h.ErrorResponseHandler.LogAndHandleError(ctx.Writer, ctx.Writer, appErr)
		return
	}

	hotels, err := h.Manager.QueryHotels(ctx, q)
	if err != nil {
		appErr := errorlog.InternalServerError(err)
		h.ErrorResponseHandler.LogAndHandleError(ctx.Writer, ctx.Writer, appErr)
		return
	}

	ctx.JSON(http.StatusOK, hotels)
}

func (h *HotelHandler) HandleGetHotelRooms(ctx *gin.Context) {
	id := ctx.Param("id")

	hotel, err := h.Manager.HotelStore.GetHotel(ctx, id)

	if err != nil {
		appErr := errorlog.InternalServerError(err)
		h.ErrorResponseHandler.LogAndHandleError(ctx.Writer, ctx.Writer, appErr)
		return
	}

	var rooms []*types.Room

	for _, roomID := range hotel.Rooms {
		room, err := h.Manager.RoomStore.GetRoomByID(ctx, roomID)
		if err != nil {
			appErr := errorlog.InternalServerError(err)
			h.ErrorResponseHandler.LogAndHandleError(ctx.Writer, ctx.Writer, appErr)
			return
		}

		rooms = append(rooms, room)
	}
	ctx.JSON(http.StatusOK, rooms)
}
