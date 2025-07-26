package exception

type ItemCreating struct{}

func (c *ItemCreating) Error() string {
	return "Failed to edit item"
}
