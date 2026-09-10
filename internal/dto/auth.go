package dto

type LoginBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshBody struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// TokenResponse is the exact shape @lattiq/webtk expects from
// /auth/v1/login and /auth/v1/refresh.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`  // "Bearer"
	ExpiresIn    int64  `json:"expires_in"`  // seconds
	ExpiresAt    int64  `json:"expires_at"`  // unix epoch seconds
}

// MeResponse is the shape GET /v1/users/me returns. webtk's User type reads
// `roles?: string[]` (plural, array) — not our internal model.User.Role
// (singular). This shapes the response to match that contract without
// changing the single-role data model underneath.
type MeResponse struct {
	ID    string   `json:"id"`
	Email string   `json:"email"`
	Name  string   `json:"name"`
	Roles []string `json:"roles"`
}