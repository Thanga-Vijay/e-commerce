package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/ecommerce/auth-service/internal/models"
	"github.com/ecommerce/auth-service/internal/repository"
	"github.com/ecommerce/auth-service/internal/utils"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	ErrUserExists        = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidToken      = errors.New("invalid or expired token")
	ErrWeakPassword      = errors.New("password does not meet strength requirements")
)

type AuthService interface {
	Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*models.TokenResponse, error)
	Logout(ctx context.Context, userID uuid.UUID, refreshToken string) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req *models.UpdateProfileRequest) (*models.User, error)
	ForgotPassword(ctx context.Context, email string) (string, error)
	ResetPassword(ctx context.Context, token, newPassword string) error
	VerifyEmail(ctx context.Context, token string) error
}

type authService struct {
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
	passwordResetRepo repository.PasswordResetRepository
	jwtManager       *utils.JWTManager
	redisClient      *redis.Client
}

func NewAuthService(
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	passwordResetRepo repository.PasswordResetRepository,
	jwtManager *utils.JWTManager,
	redisClient *redis.Client,
) AuthService {
	return &authService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		passwordResetRepo: passwordResetRepo,
		jwtManager:       jwtManager,
		redisClient:      redisClient,
	}
}

func (s *authService) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	// Check if user exists
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserExists
	}

	// Validate password strength
	if !utils.ValidatePasswordStrength(req.Password) {
		return nil, ErrWeakPassword
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         models.RoleCustomer,
		IsVerified:   false,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate tokens
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	refreshTokenStr, expiresAt, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Store refresh token
	refreshToken := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenStr,
		ExpiresAt: expiresAt,
	}
	if err := s.refreshTokenRepo.Create(ctx, refreshToken); err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User:         models.ToUserResponse(user),
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		ExpiresIn:    s.jwtManager.GetAccessTokenExpiry(),
	}, nil
}

func (s *authService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Check password
	if !utils.CheckPassword(user.PasswordHash, req.Password) {
		return nil, ErrInvalidCredentials
	}

	// Generate tokens
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	refreshTokenStr, expiresAt, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Store refresh token
	refreshToken := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenStr,
		ExpiresAt: expiresAt,
	}
	if err := s.refreshTokenRepo.Create(ctx, refreshToken); err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User:         models.ToUserResponse(user),
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		ExpiresIn:    s.jwtManager.GetAccessTokenExpiry(),
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*models.TokenResponse, error) {
	// Validate refresh token
	userID, err := s.jwtManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Check if token exists in database
	storedToken, err := s.refreshTokenRepo.FindByToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	if storedToken == nil || storedToken.ExpiresAt.Before(time.Now()) {
		return nil, ErrInvalidToken
	}

	// Get user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// Generate new access token
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	// Generate new refresh token
	newRefreshTokenStr, expiresAt, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Delete old refresh token
	if err := s.refreshTokenRepo.DeleteByToken(ctx, refreshToken); err != nil {
		return nil, err
	}

	// Store new refresh token
	newRefreshToken := &models.RefreshToken{
		UserID:    user.ID,
		Token:     newRefreshTokenStr,
		ExpiresAt: expiresAt,
	}
	if err := s.refreshTokenRepo.Create(ctx, newRefreshToken); err != nil {
		return nil, err
	}

	return &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshTokenStr,
		ExpiresIn:    s.jwtManager.GetAccessTokenExpiry(),
	}, nil
}

func (s *authService) Logout(ctx context.Context, userID uuid.UUID, refreshToken string) error {
	// Blacklist access token (stored in Redis)
	// Delete refresh token from database
	if refreshToken != "" {
		if err := s.refreshTokenRepo.DeleteByToken(ctx, refreshToken); err != nil {
			return err
		}
	}
	return nil
}

func (s *authService) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *authService) UpdateProfile(ctx context.Context, userID uuid.UUID, req *models.UpdateProfileRequest) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) ForgotPassword(ctx context.Context, email string) (string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		// Return success even if user not found (security best practice)
		return generateToken(), nil
	}

	// Generate reset token
	token := generateToken()
	expiresAt := time.Now().Add(1 * time.Hour)

	reset := &models.PasswordReset{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
		Used:      false,
	}

	if err := s.passwordResetRepo.Create(ctx, reset); err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Find reset token
	reset, err := s.passwordResetRepo.FindByToken(ctx, token)
	if err != nil {
		return err
	}
	if reset == nil {
		return ErrInvalidToken
	}

	// Validate password strength
	if !utils.ValidatePasswordStrength(newPassword) {
		return ErrWeakPassword
	}

	// Get user
	user, err := s.userRepo.FindByID(ctx, reset.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update password
	user.PasswordHash = hashedPassword
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Mark token as used
	if err := s.passwordResetRepo.MarkAsUsed(ctx, token); err != nil {
		return err
	}

	return nil
}

func (s *authService) VerifyEmail(ctx context.Context, token string) error {
	// Implementation would verify email token
	// For now, simplified version
	return nil
}

func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
