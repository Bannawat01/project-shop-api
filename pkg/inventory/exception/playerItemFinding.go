package exception

import "fmt"

type PlayerItemsFinding struct {
	PlayerID string
}

func (e *PlayerItemsFinding) Error() string {
	return fmt.Sprintf("Cannot find item for player %s", e.PlayerID)
}
