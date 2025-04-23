package handlers

import (
	"net/http"
	"strconv"

	"github.com/Noviiich/vpn-config-generator/internal/slave/service"
	"github.com/labstack/echo/v4"
)

type HandlerConfig struct {
	vpn *service.VPNService
}

func New(vpn *service.VPNService) *HandlerConfig {
	return &HandlerConfig{vpn: vpn}
}

func (h *HandlerConfig) GetVPNService(c echo.Context) error {
	idParam := c.QueryParam("user_id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid ID format",
		})
	}
	config, err := h.vpn.GetConfig(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"config": config,
	})
}

// func (h *HandlerConfig) DeleteVPNService(c echo.Context) error {
// 	idParam := c.QueryParam("user_id")
// 	userID, err := strconv.Atoi(idParam)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, map[string]string{
// 			"error": "Invalid ID format",
// 		})
// 	}
// 	err = h.vpn.(c.Request().Context(), userID)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, map[string]string{
// 			"error": err.Error(),
// 		})
// 	}
// 	return c.NoContent(http.StatusOK)
// }
