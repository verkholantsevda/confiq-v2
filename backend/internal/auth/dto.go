package auth

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	TOTPCode string `json:"totp_code"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
