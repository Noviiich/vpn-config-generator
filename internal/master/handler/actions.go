package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/Noviiich/vpn-config-generator/internal/master/service"
	"github.com/labstack/echo/v4"
)

type ActionHandler struct {
	service *service.VPNService
}

func NewActionHandler(service *service.VPNService) *ActionHandler {
	return &ActionHandler{service: service}
}

func (h *ActionHandler) GetActionsHandler(c echo.Context) error {
	query := c.QueryParam("since")
	if query == "" {
		return c.String(http.StatusBadRequest, "Отсутствует параметр 'since'")
	}

	parseTime, err := time.Parse(time.RFC3339, query)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный формат даты")
	}

	actions, err := h.service.GetActionsByTime(context.Background(), parseTime)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Ошибка при получении actions")
	}

	return c.JSON(http.StatusOK, actions)
}
