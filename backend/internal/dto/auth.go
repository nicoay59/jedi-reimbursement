package dto

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID             int64  `json:"id"`
	EmployeeNumber string `json:"employee_number"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	Position       string `json:"position"`
	Division       string `json:"division"`
	Role           string `json:"role"`
}

type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	ExpiresIn   int64        `json:"expires_in"`
	User        UserResponse `json:"user"`
}
