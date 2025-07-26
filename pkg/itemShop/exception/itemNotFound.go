package exception

type ItemNotFound struct {
	ItemID uint64
}

func (e *ItemNotFound) Error() string {
	return "Item with ID " + string(e.ItemID) + " not found"
}
