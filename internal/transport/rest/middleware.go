package rest

import (
	"errors"
	"net/http"
	"strings"

	"github.com/SavelyDev/crud-app/pkg/httputil"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) loggingMiddleware(c *gin.Context) {
	logrus.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"uri":    c.Request.URL,
	}).Info()
}

func (h *Handler) userIdentity(c *gin.Context) {
	token, err := getTokenFromRequest(c)
	if err != nil {
		httputil.NewError(c, http.StatusUnauthorized, err)
		return
	}

	userId, err := h.AuthService.ParseToken(token)
	if err != nil {
		httputil.NewError(c, http.StatusUnauthorized, err)
		return
	}

	c.Set(userCtx, userId)
}

func getTokenFromRequest(c *gin.Context) (string, error) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")

	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return headerParts[1], nil
}

func getUserId(c *gin.Context) (int, error) {
	userId, ok := c.Get("userId")
	if !ok {
		return 0, errors.New("user id not found")
	}

	userIdInt, ok := userId.(int)
	if !ok {
		return 0, errors.New("user id is of invalid type")
	}

	return userIdInt, nil
}
