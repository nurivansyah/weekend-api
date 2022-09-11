package dto

type JWT struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SignUpRequest struct {
	Username string `json:"username" binding:"required,min=4,max=15"`
	Password string `json:"password" binding:"required,min=4,max=15"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,min=4,max=15"`
	Password string `json:"password" binding:"required,min=4,max=15"`
}

type RenewTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required,jwt"`
}
