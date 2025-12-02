package excaption

type CoinAdding struct{}

func (e *CoinAdding) Error() string {
	return "Error adding coins to player account"
}
