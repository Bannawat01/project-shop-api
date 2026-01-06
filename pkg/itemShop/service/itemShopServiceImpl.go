package service

import (
	"github.com/Bannawat101/project-shop-api/entities"
	_inventoryRepository "github.com/Bannawat101/project-shop-api/pkg/inventory/repository"
	_itemShopException "github.com/Bannawat101/project-shop-api/pkg/itemShop/exception"
	_itemShopModel "github.com/Bannawat101/project-shop-api/pkg/itemShop/model"
	_itemShopRepository "github.com/Bannawat101/project-shop-api/pkg/itemShop/repository"
	_playerCoinModel "github.com/Bannawat101/project-shop-api/pkg/playerCoin/model"
	_playerCoinRepository "github.com/Bannawat101/project-shop-api/pkg/playerCoin/repository"
	"github.com/labstack/echo/v4"
)

type itemShopServiceImpl struct {
	ItemShopRepository   _itemShopRepository.ItemShopRepository //inject repository เข้ามาใช้ใน service
	playerCoinRepository _playerCoinRepository.PlayerCoinRepository
	inventoryRepository  _inventoryRepository.InventoryRepository
	logger               echo.Logger
}

func NewItemShopServiceImpl(
	itemShopRepository _itemShopRepository.ItemShopRepository,
	playerCoinRepository _playerCoinRepository.PlayerCoinRepository,
	inventoryRepository _inventoryRepository.InventoryRepository,
	logger echo.Logger,
) ItemShopService { //สร้าง instance ของ service implementation
	return &itemShopServiceImpl{itemShopRepository, playerCoinRepository, inventoryRepository, logger}
}

func (s *itemShopServiceImpl) Listing(itemFilter *_itemShopModel.ItemFilter) (*_itemShopModel.ItemResult, error) { //ดึงรายการสินค้า
	itemList, err := s.ItemShopRepository.Listing(itemFilter) //เรียกใช้ method Listing จาก repository
	if err != nil {
		return nil, err
	}

	itemCounting, err := s.ItemShopRepository.Counting(itemFilter) //เรียกใช้ method Counting จาก repository
	if err != nil {
		return nil, err
	}

	size := itemFilter.Size //ขนาดของหน้าที่ส่งมาจาก client เพื่อใช้ในการคำนวณหน้าทั้งหมด
	page := itemFilter.Page //หน้าปัจจุบันที่ส่งมาจาก client เพื่อใช้ในการสร้าง response

	totalPage := s.totalPageCalculation(itemCounting, size)     //คำนวณหาจำนวนหน้าทั้งหมด
	result := s.toItemResultResponse(itemList, page, totalPage) //แปลงข้อมูลเป็นรูปแบบ response ที่ต้องการ
	return result, nil

}

func (s *itemShopServiceImpl) totalPageCalculation(totalItem int64, size int64) int64 {
	if size == 0 {
		return 1 // or handle as needed
	}
	totalPage := totalItem / size

	if totalItem%size != 0 {
		totalPage++
	}
	return totalPage
}

// find item by id
// total price calculation
// check player coin enough or not
// record purchase history
// Coin deduction
// Inventory Felling
// return player coin
func (s *itemShopServiceImpl) Buying(buyingReq *_itemShopModel.BuyingReq) (*_playerCoinModel.PlayerCoin, error) {
	itemEntity, err := s.ItemShopRepository.FindByID(buyingReq.ItemID)
	if err != nil {
		return nil, err
	}

	totalPrice := s.totalPriceCalculation(itemEntity.ToItemModel(), buyingReq.Quantity)

	if err := s.playerCoinChecking(buyingReq.PlayerID, totalPrice); err != nil {
		return nil, err
	}

	tx := s.ItemShopRepository.BeginTransaction()

	purchaseRecording, err := s.ItemShopRepository.PurchaseHistoryRecording(&entities.PurchaseHistory{
		PlayerID:        buyingReq.PlayerID,
		ItemID:          itemEntity.ID,
		ItemName:        itemEntity.Name,
		ItemDescription: itemEntity.Description,
		ItemPrice:       itemEntity.Price,
		ItemPicture:     itemEntity.Picture,
		Quantity:        buyingReq.Quantity,
		IsBuying:        true,
	}, tx)
	if err != nil {
		s.ItemShopRepository.RollbackTransaction(tx)
		return nil, err
	}
	s.logger.Infof("Purchase history recorded: %d", purchaseRecording.ID)

	coinRecording, err := s.playerCoinRepository.CoinAdding(&entities.PlayerCoin{
		PlayerID: buyingReq.PlayerID,
		Amount:   -totalPrice,
	}, tx)
	if err != nil {
		s.ItemShopRepository.RollbackTransaction(tx)
		return nil, err
	}
	s.logger.Infof("Player coins reduced for: %d", totalPrice)

	inventoryRecording, err := s.inventoryRepository.Filling(
		buyingReq.PlayerID,
		buyingReq.ItemID,
		int(buyingReq.Quantity),
		tx,
	)
	if err != nil {
		s.ItemShopRepository.RollbackTransaction(tx)
		return nil, err
	}
	s.logger.Infof("Items recorded into player inventory: %d", len(inventoryRecording))

	if err := s.ItemShopRepository.CommitTransaction(tx); err != nil {
		s.ItemShopRepository.RollbackTransaction(tx)
		return nil, err
	}

	return coinRecording.ToPlayerCoinModel(), nil
}

// find item by id
// total price calculation
// check player itrm
// record purchase history
// Coin adding
// Inventory removing
// return player coin
func (s *itemShopServiceImpl) Selling(sellingReq *_itemShopModel.SellingReq) (*_playerCoinModel.PlayerCoin, error) {
	itemEntity, err := s.ItemShopRepository.FindByID(sellingReq.ItemID)
	if err != nil {
		return nil, err
	}

	totalPrice := s.totalPriceCalculation(itemEntity.ToItemModel(), sellingReq.Quantity)
	totalPrice = totalPrice / 2 //ขายคืนได้ครึ่งราคา

	if err := s.playerItemChecking(sellingReq.PlayerID, sellingReq.ItemID, sellingReq.Quantity); err != nil { //ตรวจสอบว่าผู้เล่นมีไอเท็มเพียงพอหรือไม่
		return nil, err
	}

	tx := s.ItemShopRepository.BeginTransaction() //เริ่ม transaction

	purchaseRecording, err := s.ItemShopRepository.PurchaseHistoryRecording(&entities.PurchaseHistory{ //บันทึกประวัติการขายไอเท็ม
		PlayerID:        sellingReq.PlayerID,
		ItemID:          itemEntity.ID,
		ItemName:        itemEntity.Name,
		ItemDescription: itemEntity.Description,
		ItemPrice:       itemEntity.Price,
		ItemPicture:     itemEntity.Picture,
		Quantity:        sellingReq.Quantity,
		IsBuying:        false, //จริงแล้วควรจะเป็น enum มากกว่ามันซื่อความหมาย //FIXME: เปลี่ยน IsBuying เป็น enum จะดีกว่า

	}, tx)
	if err != nil {
		s.ItemShopRepository.RollbackTransaction(tx)
		return nil, err
	}
	s.logger.Infof("Purchase history recorded: %d", purchaseRecording.ID)

	coinRecording, err := s.playerCoinRepository.CoinAdding(&entities.PlayerCoin{
		PlayerID: sellingReq.PlayerID,
		Amount:   totalPrice,
	}, tx)
	if err != nil {
		s.ItemShopRepository.RollbackTransaction(tx)
		return nil, err
	}
	s.logger.Infof("Player coins reduced for: %d", totalPrice)

	if err := s.inventoryRepository.Removing(
		sellingReq.PlayerID,
		sellingReq.ItemID,
		int(sellingReq.Quantity),
		tx,
	); err != nil {
		s.ItemShopRepository.RollbackTransaction(tx)
		return nil, err
	}
	s.logger.Infof("Items removed from player inventory: %d", int(sellingReq.Quantity))

	if err := s.ItemShopRepository.CommitTransaction(tx); err != nil {
		s.ItemShopRepository.RollbackTransaction(tx)
		return nil, err
	}

	return coinRecording.ToPlayerCoinModel(), nil
}

func (s *itemShopServiceImpl) toItemResultResponse(itemEntityList []*entities.Item, page, totalPage int64) *_itemShopModel.ItemResult {
	itemModelList := make([]*_itemShopModel.Item, 0, len(itemEntityList)) //ได้มาเป็น entity เราจึงต้องแปลงกับเป็น model ด้วยการ loop
	for _, item := range itemEntityList {
		itemModelList = append(itemModelList, item.ToItemModel()) //แปลงเป็น model ด้วยการเรียกใช้ function ToItemModel
	}

	return &_itemShopModel.ItemResult{
		Items: itemModelList,
		Paginate: _itemShopModel.PaginateResult{
			Page:      page,
			TotalPage: totalPage,
		},
	}
}

func (s *itemShopServiceImpl) playerCoinChecking(playerID string, totalPrice int64) error {
	playerCoin, err := s.playerCoinRepository.Showing(playerID)
	if err != nil {
		return err
	}

	if playerCoin.Coin < totalPrice {
		s.logger.Errorf("Player %s has not enough coin", playerID)
		return &_itemShopException.CoinNotEnough{}
	}

	return nil
}

func (s *itemShopServiceImpl) totalPriceCalculation(item *_itemShopModel.Item, qty uint) int64 {
	return int64(item.Price) * int64(qty)
}

func (s *itemShopServiceImpl) playerItemChecking(playerID string, itemID uint64, quantity uint) error {
	inventoryItem := s.inventoryRepository.PlayerItemCounting(playerID, itemID)

	if int(inventoryItem) < int(quantity) {
		s.logger.Errorf("Player %s does not have enough items", playerID)
		return &_itemShopException.ItemNotEnough{}
	}

	return nil
}
