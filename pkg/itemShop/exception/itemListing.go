package exception

type Itemisting struct{}

func (e *Itemisting) Error() string { //over
	return "Item listing error"
}
