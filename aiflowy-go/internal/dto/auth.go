package dto

// LoginRequest represents the login request body
type LoginRequest struct {
	Account  string `json:"account" validate:"required"` // Login name
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the login response data
type LoginResponse struct {
	Token    string `json:"token"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// CaptchaResponse represents the captcha response
type CaptchaResponse struct {
	CaptchaID   string `json:"captchaId"`
	CaptchaData string `json:"captchaData"` // Base64 encoded image
}

// GetPermissionsResponse is just a []string, handled directly
