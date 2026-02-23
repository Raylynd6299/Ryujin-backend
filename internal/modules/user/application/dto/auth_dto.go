package dto

// RegisterRequest contains the data required to register a new user.
type RegisterRequest struct {
	Email     string `json:"email"     binding:"required,email"`
	Password  string `json:"password"  binding:"required,min=8"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName"  binding:"required"`
	Locale    string `json:"locale"`
}

// LoginRequest contains the credentials for authentication.
type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest contains the refresh token to issue a new access token.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// AuthResponse is returned after a successful register or login.
// Contains the authenticated user data and both tokens.
type AuthResponse struct {
	User   UserResponse `json:"user"`
	Tokens TokenPair    `json:"tokens"`
}

// TokenPair holds the access and refresh JWT tokens.
type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
