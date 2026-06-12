package models

// Request DTOs
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"firstName" binding:"required,min=2"`
	LastName  string `json:"lastName" binding:"required,min=2"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=8"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

type UpdateProfileRequest struct {
	FirstName string `json:"firstName" binding:"min=2"`
	LastName  string `json:"lastName" binding:"min=2"`
}

// Response DTOs
type AuthResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"accessToken"`
	RefreshToken string        `json:"refreshToken"`
	ExpiresIn    int           `json:"expiresIn"`
}

type UserResponse struct {
	ID         string   `json:"id"`
	Email      string   `json:"email"`
	FirstName  string   `json:"firstName"`
	LastName   string   `json:"lastName"`
	Role       UserRole `json:"role"`
	IsVerified bool     `json:"isVerified"`
	CreatedAt  string   `json:"createdAt"`
	UpdatedAt  string   `json:"updatedAt"`
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
}

func ToUserResponse(user *User) *UserResponse {
	return &UserResponse{
		ID:         user.ID.String(),
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Role:       user.Role,
		IsVerified: user.IsVerified,
		CreatedAt:  user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
