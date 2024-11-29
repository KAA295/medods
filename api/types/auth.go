package types

type TokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type GenerateTokensRequest struct {
	UserID string `json:"guid"`
}

type RefreshTokensRequest struct {
	UserID       string `json:"guid"`
	RefreshToken string `json:"refresh_token"`
}
