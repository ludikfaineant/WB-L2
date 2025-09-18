package handler

import (
	"net/http"
	"time"
	"wb_l12/18/internal/service"

	"github.com/gin-gonic/gin"
)

type eventHandler struct {
	service *service.Service
}

func NewEventHandler(service *service.Service) *eventHandler {
	return &eventHandler{service: service}
}

func (h *eventHandler) CreateEvent(c *gin.Context) {
	var req struct {
		UserID int    `json:"user_id" binding:"required"`
		Date   string `json:"date" binding:"required"`
		Title  string `json:"title" binding:"required"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	parsedDate, err := parsedDate(req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
		return
	}

	id, err := h.service.CreateEvent(req.UserID, parsedDate, req.Title)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": id})
}

func (h *eventHandler) UpdateEvent(c *gin.Context) {
	var req struct {
		ID     int    `json:"id" binding:"required"`
		UserID int    `json:"user_id" binding:"required"`
		Date   string `json:"date" binding:"required"`
		Title  string `json:"title" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	parsedDate, err := parsedDate(req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
		return
	}
	err = h.service.UpdateEvent(req.ID, req.UserID, parsedDate, req.Title)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "successfully update"})
}

func (h *eventHandler) DeleteEvent(c *gin.Context) {
	var req struct {
		ID int `json:"id" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	err := h.service.DeleteEvent(req.ID)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "successfully delete"})
}

func (h *eventHandler) GetByDay(c *gin.Context) {
	var req struct {
		UserID int    `json:"user_id" binding:"required"`
		Date   string `json:"date" binding:"required"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}

	parsedDate, err := parsedDate(req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
		return
	}
	events, err := h.service.GetByDay(req.UserID, parsedDate)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": events})
}

func (h *eventHandler) GetByWeek(c *gin.Context) {
	var req struct {
		UserID int    `json:"user_id" binding:"required"`
		Date   string `json:"date" binding:"required"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}

	parsedDate, err := parsedDate(req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
		return
	}
	events, err := h.service.GetByWeek(req.UserID, parsedDate)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": events})
}

func (h *eventHandler) GetByMonth(c *gin.Context) {
	var req struct {
		UserID int    `json:"user_id" binding:"required"`
		Date   string `json:"date" binding:"required"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}
	parsedDate, err := parsedDate(req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
		return
	}
	events, err := h.service.GetByMonth(req.UserID, parsedDate)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": events})
}

func parsedDate(date string) (time.Time, error) {
	return time.Parse("2006-01-02", date)
}
