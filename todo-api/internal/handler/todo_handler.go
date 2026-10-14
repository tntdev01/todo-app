package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tuanntran/todo-app/todo-api/internal/dto"
	"github.com/tuanntran/todo-app/todo-api/internal/service"
)

type CreateTodoRequest struct {
	Title   string `json:"title" binding:"required,min=3,max=40"`
	Comment string `json:"comment" binding:"required,min=3,max=255"`
}

type UpdateTodoRequest struct {
	Title     string `json:"title" binding:"required,min=3,max=40"`
	Comment   string `json:"comment" binding:"required,min=3,max=255"`
	Completed bool   `json:"completed"`
}

type TodoHandler struct {
	Service *service.TodoService
}

// @Summary Create a new todo item
// @Description Create a new todo item with the given title and comment
// @Tags todos
// @Accept json
// @Produce json
// @Param todo body CreateTodoRequest true "Todo to create"
// @Success 201 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /todos [post]
func (h *TodoHandler) Create(c *gin.Context) {

	var req CreateTodoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	err := h.Service.Create(req.Title, req.Comment)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{
		Success: true,
		Data:    "created",
	})
}

// @Summary Get all todo items
// @Description Retrieve a list of all todo items
// @Tags todos
// @Accept json
// @Produce json
// @Success 200 {object} dto.SuccessResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /todos [get]
func (h *TodoHandler) GetAll(c *gin.Context) {
	todos, err := h.Service.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Data:    todos,
	})
}

// @Summary Get a todo item by ID
// @Description Retrieve a single todo item by its ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /todos/{id} [get]
func (h *TodoHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	todo, err := h.Service.GetByID(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Data:    todo,
	})
}

// @Summary Update a todo item
// @Description Update an existing todo item by its ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Param todo body UpdateTodoRequest true "Todo to update"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /todos/{id} [put]
func (h *TodoHandler) Update(c *gin.Context) {

	id := c.Param("id")

	var req UpdateTodoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	err := h.Service.Update(
		id,
		req.Title,
		req.Comment,
		req.Completed,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Data:    "updated",
	})
}

// @Summary Delete a todo item
// @Description Delete an existing todo item by its ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /todos/{id} [delete]
func (h *TodoHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.Service.Delete(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Data:    "deleted",
	})
}
