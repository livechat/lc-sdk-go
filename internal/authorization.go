package internal

// Token represents SSO token from Chat API's perspective.
type Token struct {
	// LicenseID specifies ID of license which owns the token.
	LicenseID int
	// AccessToken is a customer access token returned by LiveChat OAuth Server.
	AccessToken string
	// Region is a datacenter for LicenseID (`dal` or `fra`).
	Region string
}

// TokenGetter is called by each API method to obtain valid Token.
// If TokenGetter returns nil, the method won't be executed on API.
type TokenGetter func() *Token
