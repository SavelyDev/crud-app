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

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		httputil.NewError(c, http.StatusUnauthorized, errors.New("empty auth header"))
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		httputil.NewError(c, http.StatusUnauthorized, errors.New("invalid auth header"))
		return
	}

	userId, err := h.AuthService.ParseToken(headerParts[1])
	if err != nil {
		httputil.NewError(c, http.StatusUnauthorized, err)
		return
	}

	c.Set(userCtx, userId)
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

func (h *Handler) loggingMiddleware(c *gin.Context) {
	logrus.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"uri":    c.Request.URL,
	}).Info()
}
