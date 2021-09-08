package blogServer

type AuthorizationHeader struct {
	Token string `header:"authorization" binding:"required"`
}
