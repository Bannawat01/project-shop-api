package exception

import "fmt"

type ItemAchiving struct {
	Item uint64
}

func (a *ItemAchiving) Error() string {
	return fmt.Sprintf("Failed to archive item with ID: %d", a.Item)
}
