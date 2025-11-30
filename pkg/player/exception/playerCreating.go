package exception

type PlayerCreating struct {
	PlayerID string
}

func (e *PlayerCreating) Error() string {
	return "Error occurred while creating player"
}
