package rest

import (
	"fmt"
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

	id, err := h.AuthService.SignUp(user)
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
// @Success 200 {string} string "acces_token"
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /auth/sign-in [get]
func (h *Handler) signIn(c *gin.Context) {
	var credentials domain.SignInInput

	if err := c.BindJSON(&credentials); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	accesToken, refreshToken, err := h.AuthService.SignIn(credentials)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	c.Header("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", refreshToken))

	c.JSON(http.StatusOK, gin.H{"acces_token": accesToken})
}

func (h *Handler) refresh(c *gin.Context) {
	token, err := c.Cookie("refresh-token")
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	accesToken, refreshToken, err := h.AuthService.RefreshToken(token)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	c.Header("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", refreshToken))

	c.JSON(http.StatusOK, gin.H{"acces_token": accesToken})
}
