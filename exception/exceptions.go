package exception

type FORUMEXCEPTIONCODE int

const (
	NonForumException         FORUMEXCEPTIONCODE = 1000
	ErrInvalidRequest         FORUMEXCEPTIONCODE = 4001
	ErrMiddlewareAuthenticate FORUMEXCEPTIONCODE = 4002
	ErrGetUserInfo            FORUMEXCEPTIONCODE = 5002
	ErrLogin                  FORUMEXCEPTIONCODE = 5003
	ErrCreateThread           FORUMEXCEPTIONCODE = 5005
	ErrUpdateUserInfo         FORUMEXCEPTIONCODE = 5009
	ErrNameExists             FORUMEXCEPTIONCODE = 4091
	ErrMailExists             FORUMEXCEPTIONCODE = 4092
	ErrTitleExists            FORUMEXCEPTIONCODE = 4093
	ErrUserNotFound           FORUMEXCEPTIONCODE = 4041
	ErrMySQLInternalError     FORUMEXCEPTIONCODE = 5001
)
