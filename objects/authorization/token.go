package authorization

type Token struct {
	LicenseID   int
	AccessToken string
	Region      string
}

type TokenGetter func() *Token
