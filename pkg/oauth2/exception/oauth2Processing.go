package exception

type OAuth2Processing struct{}

func (e *OAuth2Processing) Error() string {
	return "Error processing OAuth2 request"
}
