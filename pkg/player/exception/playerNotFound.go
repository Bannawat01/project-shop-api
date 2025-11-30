package exception

import "fmt"

type PlayerNotFound struct {
	PlayerID string
}

func (e *PlayerNotFound) Error() string {
	return fmt.Sprintf("Player with ID %s not found", e.PlayerID)
}
