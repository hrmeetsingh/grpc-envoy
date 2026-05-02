package port

type TokenSigner interface {
	Sign(userID string, tenant string) (string, error)
	Verify(tokenString string) (claims map[string]interface{}, err error)
}
