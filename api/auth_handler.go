package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mkabdelrahman/hotel-reservation/auth"
	"github.com/mkabdelrahman/hotel-reservation/db"
	"github.com/mkabdelrahman/hotel-reservation/errorlog"
	"github.com/mkabdelrahman/hotel-reservation/types"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	store                db.UserStore
	ErrorResponseHandler *errorlog.HTTPErrorResponseWriterAndLogger
}

func NewAuthHandler(store db.UserStore, errorLogger *log.Logger) *AuthHandler {
	return &AuthHandler{
		store:                store,
		ErrorResponseHandler: &errorlog.HTTPErrorResponseWriterAndLogger{Logger: errorLogger},
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) HandleAuthenticate(c *gin.Context) {
	var authParams AuthParams

	if err := c.BindJSON(&authParams); err != nil {
		appErr := errorlog.BadRequestError(err)
		h.ErrorResponseHandler.LogAndHandleError(c.Writer, c.Writer, appErr)
		return
	}

	user, err := h.store.GetUserByEmail(c, authParams.Email)
	if err != nil {
		appErr := errorlog.UnauthorizedError(err)
		h.ErrorResponseHandler.LogAndHandleError(c.Writer, c.Writer, appErr)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(authParams.Password)); err != nil {
		appErr := errorlog.UnauthorizedError(err)
		h.ErrorResponseHandler.LogAndHandleError(c.Writer, c.Writer, appErr)
		return
	}

	ok, err := types.IsValidPassword(user.EncryptedPassword, authParams.Password)

	if err != nil || !ok {
		appErr := errorlog.UnauthorizedError(err)
		h.ErrorResponseHandler.LogAndHandleError(c.Writer, c.Writer, appErr)
		return
	}

	token, err := auth.GenerateAuthToken(user.ID.Hex())

	if err != nil {
		appErr := errorlog.InternalServerError(err)
		h.ErrorResponseHandler.LogAndHandleError(c.Writer, c.Writer, appErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
