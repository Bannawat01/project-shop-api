package exception

type AdminCreating struct {
	AdminID string
}

func (e *AdminCreating) Error() string {
	return "Error occurred while creating admin"
}
