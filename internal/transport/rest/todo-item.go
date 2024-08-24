package rest

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/SavelyDev/crud-app/internal/domain"
	"github.com/SavelyDev/crud-app/pkg/httputil"
	"github.com/gin-gonic/gin"
)

// @Summary Create Item
// @Description Create a new todo item
// @Security ApiKeyAuth
// @Tags items
// @Accept json
// @Produce json
// @Param id path int true "List ID"
// @Param item body domain.TodoItem true "Todo item info"
// @Success 200 {integer} integer 1
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/lists/{id}/items [post]
func (h *Handler) createItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	var item domain.TodoItem

	if err := c.BindJSON(&item); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	id, err := h.TodoItemService.CreateItem(userId, listId, item)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

// @Summary Get All Items
// @Description Get all todo items for a specific list
// @Security ApiKeyAuth
// @Tags items
// @Produce json
// @Param id path int true "List ID"
// @Success 200 {array} domain.TodoItem
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/lists/{id}/items [get]
func (h *Handler) getAllItems(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	items, err := h.TodoItemService.GetAllItems(userId, listId)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, items)
}

// @Summary Get Item By ID
// @Description Get a specific todo item by its ID
// @Security ApiKeyAuth
// @Tags items
// @Produce json
// @Param id path int true "Item ID"
// @Success 200 {object} domain.TodoItem
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/items/{id} [get]
func (h *Handler) getItemById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	item, err := h.TodoItemService.GetItemById(userId, itemId)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, item)
}

// @Summary Update Item
// @Description Update an existing todo item
// @Security ApiKeyAuth
// @Tags items
// @Accept json
// @Produce json
// @Param id path int true "Item ID"
// @Param item body domain.UpdateItemInput true "Updated item info"
// @Success 200 {string} string "ok"
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/items/{id} [put]
func (h *Handler) updateItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, errors.New("invalid id param"))
		return
	}

	var item domain.UpdateItemInput
	if err := c.BindJSON(&item); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	if err := h.TodoItemService.UpdateItem(userId, itemId, item); err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// @Summary Delete Item
// @Description Delete a specific todo item by its ID
// @Security ApiKeyAuth
// @Tags items
// @Produce json
// @Param id path int true "Item ID"
// @Success 200 {string} string "ok"
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /api/items/{id} [delete]
func (h *Handler) deleteItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	err = h.TodoItemService.DeleteItem(userId, itemId)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
