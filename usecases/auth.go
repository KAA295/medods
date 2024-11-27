package usecases

type AuthService interface {
	GenerateAccessToken()
	GenerateRefreshToken()
	GenerateTokens()
	RefreshTokens()
}
