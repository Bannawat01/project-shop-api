package exception

type ItemNotFound struct {
	ItemID string
}

func (e *ItemNotFound) Error() string {
	return "Item with ID " + string(e.ItemID) + " not found"
}
