package rest

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/SavelyDev/crud-app/internal/domain"
	"github.com/SavelyDev/crud-app/pkg/httputil"
	"github.com/gin-gonic/gin"
)

// @Summary Create List
// @Description Create a new todo list
// @Security ApiKeyAuth
// @Tags lists
// @Accept json
// @Produce json
// @Param list body domain.TodoList true "Todo list info"
// @Success 200 {integer} integer 1
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/lists [post]
func (h *Handler) createList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	var list domain.TodoList

	if err := c.BindJSON(&list); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	id, err := h.TodoListService.CreateList(userId, list)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

// @Summary Get All Lists
// @Description Get all todo lists for a user
// @Security ApiKeyAuth
// @Tags lists
// @Produce json
// @Success 200 {array} domain.TodoList
// @Failure 500 {object} httputil.HTTPError
// @Router /api/lists [get]
func (h *Handler) getAllLists(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	lists, err := h.TodoListService.GetAllLists(userId)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, lists)
}

// @Summary Get List By ID
// @Description Get a specific todo list by its ID
// @Security ApiKeyAuth
// @Tags lists
// @Produce json
// @Param id path int true "List ID"
// @Success 200 {object} domain.TodoList
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/lists/{id} [get]
func (h *Handler) getListById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, errors.New("invalid id param"))
		return
	}

	list, err := h.TodoListService.GetListById(userId, listId)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, list)
}

// @Summary Update List
// @Description Update an existing todo list
// @Security ApiKeyAuth
// @Tags lists
// @Accept json
// @Produce json
// @Param id path int true "List ID"
// @Param list body domain.UpdateListInput true "Updated list info"
// @Success 200 {string} string "ok"
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/lists/{id} [put]
func (h *Handler) updateList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, errors.New("invalid id param"))
		return
	}

	var list domain.UpdateListInput
	if err := c.BindJSON(&list); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	if err := h.TodoListService.UpdateList(userId, listId, list); err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// @Summary Delete List
// @Description Delete a specific todo list by its ID
// @Security ApiKeyAuth
// @Tags lists
// @Produce json
// @Param id path int true "List ID"
// @Success 200 {string} string "ok"
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/lists/{id} [delete]
func (h *Handler) deleteList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, errors.New("invalid id param"))
		return
	}

	if err := h.TodoListService.DeleteList(userId, listId); err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
