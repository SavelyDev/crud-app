package rest

import (
	"net/http"

	"github.com/SavelyDev/crud-app/internal/domain"
	"github.com/SavelyDev/crud-app/pkg/httputil"
	"github.com/gin-gonic/gin"
)

// @Summary Sign Up
// @Description Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body domain.User true "User info"
// @Success 200 {integer} integer 1
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var user domain.User

	if err := c.BindJSON(&user); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	id, err := h.AuthService.CreateUser(user)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

// @Summary Sign In
// @Description Authenticate a user and return a token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body domain.SignInInput true "Sign in credentials"
// @Success 200 {string} string "token"
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var credentials domain.SignInInput

	if err := c.BindJSON(&credentials); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	token, err := h.AuthService.GenerateToken(credentials.Email, credentials.PasswordHash)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
