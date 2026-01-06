package exception

type CoinNotEnough struct{}

func (e *CoinNotEnough) Error() string {
	return "Not enough coins"
}
