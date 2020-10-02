package authorization

// TokenType represents authentication scheme.
type TokenType int

// Supported values of TokenType.
const (
	BearerToken TokenType = iota
	BasicToken
)

func (t TokenType) String() string {
	if t == BasicToken {
		return "Basic"
	} else if t == BearerToken {
		return "Bearer"
	}
	return "Unknown"
}
