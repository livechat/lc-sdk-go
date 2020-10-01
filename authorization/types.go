package authorization

// TokenType represents Bearer or Basic authentication scheme.
type TokenType int

// Possible values of TokenType.
const (
	BearerToken TokenType = iota
	BasicToken
)
