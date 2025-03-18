package auth

// BaseAuth interface for authentication
type BaseAuth interface {
	AuthenticateByToken(token string) (uid string, err error)
}
