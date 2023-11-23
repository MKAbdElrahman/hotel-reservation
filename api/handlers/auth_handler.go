package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mkabdelrahman/hotel-reservation/business"
	"github.com/mkabdelrahman/hotel-reservation/errorlog"
)

type AuthHandler struct {
	Manager              *business.Manager
	ErrorResponseHandler *errorlog.HTTPErrorResponseWriterAndLogger
}

func NewAuthHandler(Manager *business.Manager, errorLogger *log.Logger) *AuthHandler {
	return &AuthHandler{
		Manager:              Manager,
		ErrorResponseHandler: &errorlog.HTTPErrorResponseWriterAndLogger{Logger: errorLogger},
	}
}

func (h *AuthHandler) HandleAuthenticate(c *gin.Context) {
	var authParams business.AuthParams

	if err := c.BindJSON(&authParams); err != nil {
		appErr := errorlog.BadRequestError(err)
		h.ErrorResponseHandler.LogAndHandleError(c.Writer, c.Writer, appErr)
		return
	}

	token, err := h.Manager.GetUserToken(c, authParams)

	if err != nil {
		appErr := errorlog.UnauthorizedError(err)
		h.ErrorResponseHandler.LogAndHandleError(c.Writer, c.Writer, appErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
