package entity

type LoginRequest struct {
	Username    string `form:"username"`
	Email       string `form:"email" binding:"email"`
	PhoneNumber string `form:"phone_number" binding:"e164"`
	Password    string `form:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
