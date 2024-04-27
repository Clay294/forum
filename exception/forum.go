package exception

type ForumException struct {
	HttpCode           int                `json:"http_code"`
	ForumExceptionCode FORUMEXCEPTIONCODE `json:"forum_exception_code"`
	Message            string             `json:"message"`
}

func NewForumException(hc int, err error) *ForumException {
	return &ForumException{
		HttpCode:           hc,
		ForumExceptionCode: FORUMEXCEPTIONCODE(hc),
		Message:            err.Error(),
	}
}

func (fe *ForumException) Error() string {
	return fe.Message
}

func (fe *ForumException) SetForumExceptionCode(feCode FORUMEXCEPTIONCODE) *ForumException {
	fe.ForumExceptionCode = feCode
	return fe
}
