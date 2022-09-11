package dto

type UserSession struct {
	ID        uint
	Username  string
	User_type string
}

type UserProfile struct {
	User UserSession `json:"user"`
}

type CreateUserRequest struct {
	Username  string `json:"username" binding:"required,min=4,max=15"`
	Password  string `json:"password" binding:"required,min=4,max=15"`
	User_type string `json:"user_type" binding:"required,eq=ADMIN|eq=USER"`
}

type UpdateUserRequest struct {
	Username  string `json:"username" binding:"omitempty,min=4,max=15"`
	Password  string `json:"password" binding:"omitempty,min=4,max=15"`
	User_type string `json:"user_type" binding:"omitempty,eq=ADMIN|eq=USER"`
}
