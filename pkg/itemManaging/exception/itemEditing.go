package exception

import "fmt"

type ItemEditing struct {
	ItemID uint64
}

func (e *ItemEditing) Error() string {
	return fmt.Sprintf("Failed to edit item with ID: %d", e.ItemID)
}
