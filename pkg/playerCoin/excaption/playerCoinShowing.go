package excaption

type PlayerCoinShowing struct{}

func (e *PlayerCoinShowing) Error() string {
	return "Error showing player coins"
}
