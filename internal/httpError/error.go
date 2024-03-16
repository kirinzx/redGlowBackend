package httpError

type HTTPError interface{
	Error() string
	Status() int
}

type forbiddenError struct {
	message string
	generalMessage *string
	status int
}

func NewForbiddenError(msg string) *forbiddenError{
	return &forbiddenError{
		message:msg,
		status: 403,
	}
}

func (e *forbiddenError) Error() string{
	return e.message
}


func (e *forbiddenError) Status() int{
	return e.status
}


type badRequestError struct {
	message string
	status int
}

func NewBadRequestError(msg string) *badRequestError{
	return &badRequestError{
		message:msg,
		status:400,
	}
}

func (e *badRequestError) Error() string{
	return e.message
}

func (e *badRequestError) Status() int{
	return e.status
}

type notFoundError struct {
	message string
	status int
}

func NewNotFoundError(msg string) *notFoundError{
	return &notFoundError{
		message:msg,
		status:404,
	}
}

func (e *notFoundError) Error() string{
	return e.message
}

func (e *notFoundError) Status() int{
	return e.status
}

type internalServerError struct {
	message string
	status int
}

func NewInternalServerError(msg string) *internalServerError{
	return &internalServerError{
		message:msg,
		status:500,
	}
}

func (e *internalServerError) Error() string{
	return e.message
}

func (e *internalServerError) Status() int{
	return e.status
}