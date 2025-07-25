package exception

type Itemisting struct{}

func (e *Itemisting) Error() string {
	return "Item listing error"
}
