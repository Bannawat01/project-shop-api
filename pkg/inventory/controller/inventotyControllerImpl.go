package controller

import (
	"net/http"

	"github.com/Bannawat101/project-shop-api/pkg/custom"
	_inventoryService "github.com/Bannawat101/project-shop-api/pkg/inventory/service"
	"github.com/Bannawat101/project-shop-api/pkg/validation"
	"github.com/labstack/echo/v4"
)

type inventoryControllerImpl struct {
	inventoryService _inventoryService.InventoryService
	logger           echo.Logger
}

func NewInventoryControllerImpl(
	inventoryService _inventoryService.InventoryService,
	logger echo.Logger,
) InventoryController {
	return &inventoryControllerImpl{
		inventoryService,
		logger,
	}
}

func (c *inventoryControllerImpl) Listing(pctx echo.Context) error {
	playerID, err := validation.PlayerIDGetting(pctx)
	if err != nil {
		c.logger.Error("InventoryController - Listing: ", err.Error())
		return custom.CustomError(pctx, http.StatusBadRequest, err)
	}

	inventoryListing, err := c.inventoryService.Listing(playerID)
	if err != nil {
		c.logger.Error("InventoryController - Listing: ", err.Error())
		return custom.CustomError(pctx, http.StatusInternalServerError, err)
	}
	return pctx.JSON(http.StatusOK, inventoryListing)
}
