package exception

type NotFoundException struct {
	Message string
}

func (impl *NotFoundException) Error() string {
	return impl.Message
}
