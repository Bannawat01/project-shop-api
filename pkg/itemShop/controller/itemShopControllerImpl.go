package controller

import (
	"net/http"

	"github.com/Bannawat101/project-shop-api/pkg/custom"
	_itemShopModel "github.com/Bannawat101/project-shop-api/pkg/itemShop/model"
	_itemShopService "github.com/Bannawat101/project-shop-api/pkg/itemShop/service"
	"github.com/Bannawat101/project-shop-api/pkg/validation"
	"github.com/labstack/echo/v4"
)

type itemShopControllerImpl struct {
	itemShopService _itemShopService.ItemShopService
}

func NewItemShopControllerImpl(itemShopService _itemShopService.ItemShopService) ItemShopController {
	return &itemShopControllerImpl{itemShopService}
}

func (c *itemShopControllerImpl) Listing(pctx echo.Context) error { //override เมธอด Listing จาก interface
	itemFilter := new(_itemShopModel.ItemFilter)
	customEchoRequest := custom.NewCustomEchoRequest(pctx)

	if err := customEchoRequest.Bind(itemFilter); err != nil {
		return custom.CustomError(pctx, http.StatusBadRequest, err) // Handle validation error with custom error response
	}

	// Set default pagination values if not provided
	if itemFilter.Page == 0 {
		itemFilter.Page = 1
	}
	if itemFilter.Size == 0 {
		itemFilter.Size = 10
	}

	itemModelist, err := c.itemShopService.Listing(itemFilter)
	if err != nil {
		return custom.CustomError(pctx, http.StatusInternalServerError, err) // Handle error with custom error response
	}
	return pctx.JSON(http.StatusOK, itemModelist) //ส่งกลับไปเป็น json
	// return custom.CustomError(pctx, http.StatusInternalServerError, (&_itemShopException.Itemisting{}).Error()) // Return custom error for item listing
}

func (c *itemShopControllerImpl) Buying(pctx echo.Context) error {
	playerID, err := validation.PlayerIDGetting(pctx)
	if err != nil {
		return custom.CustomError(pctx, http.StatusBadRequest, err) // Handle validation error with custom error response
	}

	buyingReq := new(_itemShopModel.BuyingReq)
	customEchoRequest := custom.NewCustomEchoRequest(pctx)

	if err := customEchoRequest.Bind(buyingReq); err != nil {
		return custom.CustomError(pctx, http.StatusBadRequest, err) // Handle validation error with custom error response
	}

	buyingReq.PlayerID = playerID

	playerCoin, err := c.itemShopService.Buying(buyingReq)
	if err != nil {
		return custom.CustomError(pctx, http.StatusInternalServerError, err) // Handle error with custom error response
	}

	return pctx.JSON(http.StatusOK, playerCoin)
}

func (c *itemShopControllerImpl) Selling(pctx echo.Context) error {
	playerID, err := validation.PlayerIDGetting(pctx)
	if err != nil {
		return custom.CustomError(pctx, http.StatusBadRequest, err) // Handle validation error with custom error response
	}

	sellingReq := new(_itemShopModel.SellingReq)
	customEchoRequest := custom.NewCustomEchoRequest(pctx)

	if err := customEchoRequest.Bind(sellingReq); err != nil {
		return custom.CustomError(pctx, http.StatusBadRequest, err) // Handle validation error with custom error response
	}

	sellingReq.PlayerID = playerID
	playerCoin, err := c.itemShopService.Selling(sellingReq)
	if err != nil {
		return custom.CustomError(pctx, http.StatusInternalServerError, err) // Handle error with custom error response
	}

	return pctx.JSON(http.StatusOK, playerCoin)
}
