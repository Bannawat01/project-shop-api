package controller

import (
	"net/http"
	"strconv"

	"github.com/Bannawat101/project-shop-api/pkg/custom"
	_itemManagingModel "github.com/Bannawat101/project-shop-api/pkg/itemManaging/model"
	_itemManagingService "github.com/Bannawat101/project-shop-api/pkg/itemManaging/service"
	"github.com/Bannawat101/project-shop-api/pkg/validation"
	"github.com/labstack/echo/v4"
)

type itemManagingControllerImpl struct {
	itemManagingService _itemManagingService.ItemManagingService
}

func NewItemManagingControllerImpl(
	itemManagingService _itemManagingService.ItemManagingService,
) ItemManagingController {
	return &itemManagingControllerImpl{itemManagingService}
}

func (c *itemManagingControllerImpl) Creating(pctx echo.Context) error {
	adminID, err := validation.AdminIDGetting(pctx)
	if err != nil {
		return custom.CustomError(pctx, http.StatusBadRequest, err) //
	}

	itemCreatingReq := new(_itemManagingModel.ItemCreateReq)

	customEchoRequest := custom.NewCustomEchoRequest(pctx)

	if err := customEchoRequest.Bind(itemCreatingReq); err != nil {
		return custom.CustomError(pctx, http.StatusBadRequest, err) // Handle validation error with custom error response
	}
	itemCreatingReq.AdminID = adminID

	item, err := c.itemManagingService.Creating(itemCreatingReq)
	if err != nil {
		return custom.CustomError(pctx, http.StatusInternalServerError, err) // Handle error with custom error response

	}

	return pctx.JSON(http.StatusCreated, item) // Return the created item as JSON response

}

func (c *itemManagingControllerImpl) Editing(pctx echo.Context) error {
	itemID, err := c.getItemID(pctx)

	if err != nil {
		return custom.CustomError(pctx, http.StatusBadRequest, err) //
	}

	itemEditingReq := new(_itemManagingModel.ItemEditingReq)

	customEchoRequest := custom.NewCustomEchoRequest(pctx)
	if err := customEchoRequest.Bind(itemEditingReq); err != nil {
		return custom.CustomError(pctx, http.StatusBadRequest, err) // Handle validation error with custom error response
	}

	item, err := c.itemManagingService.Editing(itemID, itemEditingReq)
	if err != nil {
		return custom.CustomError(pctx, http.StatusInternalServerError, err) // Handle error with custom error response
	}
	return pctx.JSON(http.StatusOK, item) // Return the updated item as JSON response

}

func (c *itemManagingControllerImpl) Archiving(pctx echo.Context) error {
	itemID, err := c.getItemID(pctx)

	if err != nil {
		return custom.CustomError(pctx, http.StatusBadRequest, err) //
	}

	if err = c.itemManagingService.Archive(itemID); err != nil {
		return custom.CustomError(pctx, http.StatusInternalServerError, err) // Handle error with custom error response
	}
	return pctx.NoContent(http.StatusNoContent) // Return no content response for successful archiving

}

func (c *itemManagingControllerImpl) getItemID(pctx echo.Context) (uint64, error) {
	itemIDStr := pctx.Param("itemID")
	itemIDuint, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint64(itemIDuint), nil
}
