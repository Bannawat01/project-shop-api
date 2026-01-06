package exception

type ItemListing struct{}

func (e *ItemListing) Error() string { //over
	return "Item listing error"
}
