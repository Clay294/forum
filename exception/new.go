package exception

import "net/http"

func NewNonForumException(err error) *ForumException {
	return NewForumException(http.StatusServiceUnavailable, err).SetForumExceptionCode(NonForumException)
}

func NewErrInvalidRequest(err error) *ForumException {
	return NewForumException(http.StatusBadRequest, err).SetForumExceptionCode(ErrInvalidRequest)
}

func NewErrMiddlewareAuthenticate(err error) *ForumException {
	return NewForumException(http.StatusUnauthorized, err).SetForumExceptionCode(ErrMiddlewareAuthenticate)
}

func NewErrLogin(err error) *ForumException {
	return NewForumException(http.StatusInternalServerError, err).SetForumExceptionCode(ErrLogin)
}

func NewErrGetUserInfo(err error) *ForumException {
	return NewForumException(http.StatusInternalServerError, err).SetForumExceptionCode(ErrGetUserInfo)
}

func NewErrCreateThrad(err error) *ForumException {
	return NewForumException(http.StatusInternalServerError, err).SetForumExceptionCode(ErrCreateThread)
}

func NewErrUpdateUserInfo(err error) *ForumException {
	return NewForumException(http.StatusInternalServerError, err).SetForumExceptionCode(ErrUpdateUserInfo)
}

func NewErrNameExists(err error) *ForumException {
	return NewForumException(http.StatusConflict, err).SetForumExceptionCode(ErrNameExists)
}

func NewErrMailExists(err error) *ForumException {
	return NewForumException(http.StatusConflict, err).SetForumExceptionCode(ErrMailExists)
}

func NewErrTitleExists(err error) *ForumException {
	return NewForumException(http.StatusConflict, err).SetForumExceptionCode(ErrTitleExists)
}

func NewErrUserNotFound(err error) *ForumException {
	return NewForumException(http.StatusNotFound, err).SetForumExceptionCode(ErrUserNotFound)
}

func NewErrMySQLInternalError(err error) *ForumException {
	return NewForumException(http.StatusInternalServerError, err).SetForumExceptionCode(ErrMySQLInternalError)
}

func NewErrInternalServerError(err error) *ForumException {
	return NewForumException(http.StatusInternalServerError, err)
}

func NewErrUnauthorized(err error) *ForumException {
	return NewForumException(http.StatusUnauthorized, err)
}
