package users

type LOGINBY byte

const (
	CreateUserUrl     = "/user"
	GetUserInfoUrl    = "/user"
	LoginUrl          = "/login"
	LogoutUrl         = "/logout"
	LogoutMethod      = "/get"
	ValidateTokenUrl  = "/is_valid_token"
	UpdateUserInfoUrl = "user"
	CreateUserMethod  = "post"
	LoginMethod
	GetUserInfoMethod = "get"
	ValidateTokenMethod
	UpdateUserInfoMethod = "put"
)

const (
	LoginByName LOGINBY = iota + 1
	LoginByMail
)

const (
	UnitName                 = "users"
	TokensTableName          = "tokens"
	JsonTokenCookieName      = "json_token"
	TokensBlacklistTableName = "tokens_blacklist"
	TokenMetaContextName     = "token_meta"
)
