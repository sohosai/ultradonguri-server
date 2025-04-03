package error_handling

type InternalErrors string

const (
	CastError     InternalErrors = "Type casting error"
	NotFoundError InternalErrors = "Not found error"
)

func (ie InternalErrors) Error() string {
	return string(ie)
}
