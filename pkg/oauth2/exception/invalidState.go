package exception

type InvalidState struct{}

func (e *InvalidState) Error() string {
	return "Invalid OAuth2 state parameter"
}
