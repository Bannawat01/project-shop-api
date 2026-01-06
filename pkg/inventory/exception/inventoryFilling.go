package exception

import "fmt"

type InventoryFilling struct {
	PlayerID string
	ItemID   uint64
}

func (e *InventoryFilling) Error() string {
	return fmt.Sprintf("Cannot fill inventory for player %s with item ID %d", e.PlayerID, e.ItemID)
}
