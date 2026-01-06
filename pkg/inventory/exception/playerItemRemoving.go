package exception

import "fmt"

type PlayerItemRemoving struct {
	ItemID uint64
}

func (e *PlayerItemRemoving) Error() string {
	return fmt.Sprintf("Cannot remove item with ID %d", e.ItemID)
}
