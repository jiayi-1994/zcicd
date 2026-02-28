package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/zcicd/zcicd-server/internal/auth/model"
	"github.com/zcicd/zcicd-server/internal/auth/repository"
	"github.com/zcicd/zcicd-server/pkg/config"
	apperrors "github.com/zcicd/zcicd-server/pkg/errors"
	"github.com/zcicd/zcicd-server/pkg/middleware"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repository.UserRepository
	cfg      *config.Config
	rdb      *redis.Client
}

func NewAuthService(userRepo *repository.UserRepository, cfg *config.Config, rdb *redis.Client) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		cfg:      cfg,
		rdb:      rdb,
	}
}

func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*model.User, error) {
	// Check if username already exists
	if _, err := s.userRepo.FindByUsername(ctx, req.Username); err == nil {
		return nil, apperrors.New(40900, "用户名已存在")
	}

	// Check if email already exists
	if _, err := s.userRepo.FindByEmail(ctx, req.Email); err == nil {
		return nil, apperrors.New(40900, "邮箱已存在")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.Wrap(50000, "密码加密失败", err)
	}

	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		DisplayName:  req.DisplayName,
		Status:       "active",
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, apperrors.Wrap(50001, "创建用户失败", err)
	}

	// Assign default viewer role
	defaultRole := &model.UserRole{
		UserID:    user.ID,
		Role:      "viewer",
		ScopeType: "system",
		ScopeID:   "00000000-0000-0000-0000-000000000000",
	}
	if err := s.userRepo.AddUserRole(ctx, defaultRole); err != nil {
		return nil, apperrors.Wrap(50001, "分配默认角色失败", err)
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*TokenResponse, error) {
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrLoginFailed
		}
		return nil, apperrors.Wrap(50001, "查询用户失败", err)
	}

	if user.Status != "active" {
		return nil, apperrors.ErrUserDisabled
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, apperrors.ErrLoginFailed
	}

	roles, err := s.userRepo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, apperrors.Wrap(50001, "查询用户角色失败", err)
	}

	tokenResp, err := s.generateTokenPair(user, roles)
	if err != nil {
		return nil, err
	}

	// Store refresh token in Redis
	refreshKey := fmt.Sprintf("refresh_token:%s", tokenResp.RefreshToken)
	refreshExpire := time.Duration(s.cfg.JWT.RefreshExpireHours) * time.Hour
	if err := s.rdb.Set(ctx, refreshKey, user.ID, refreshExpire).Err(); err != nil {
		return nil, apperrors.Wrap(50000, "存储刷新令牌失败", err)
	}

	// Update last login time
	now := time.Now()
	user.LastLoginAt = &now
	_ = s.userRepo.Update(ctx, user)

	return tokenResp, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	refreshKey := fmt.Sprintf("refresh_token:%s", refreshToken)
	userID, err := s.rdb.Get(ctx, refreshKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, apperrors.ErrTokenInvalid
		}
		return nil, apperrors.Wrap(50000, "验证刷新令牌失败", err)
	}

	// Delete old refresh token
	s.rdb.Del(ctx, refreshKey)

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, apperrors.Wrap(50001, "查询用户失败", err)
	}

	if user.Status != "active" {
		return nil, apperrors.ErrUserDisabled
	}

	roles, err := s.userRepo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, apperrors.Wrap(50001, "查询用户角色失败", err)
	}

	tokenResp, err := s.generateTokenPair(user, roles)
	if err != nil {
		return nil, err
	}

	// Store new refresh token in Redis
	newRefreshKey := fmt.Sprintf("refresh_token:%s", tokenResp.RefreshToken)
	refreshExpire := time.Duration(s.cfg.JWT.RefreshExpireHours) * time.Hour
	if err := s.rdb.Set(ctx, newRefreshKey, user.ID, refreshExpire).Err(); err != nil {
		return nil, apperrors.Wrap(50000, "存储刷新令牌失败", err)
	}

	return tokenResp, nil
}

func (s *AuthService) GetProfile(ctx context.Context, userID string) (*model.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.Wrap(50001, "查询用户失败", err)
	}
	return user, nil
}

func (s *AuthService) UpdateProfile(ctx context.Context, userID string, req *UpdateProfileRequest) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrNotFound
		}
		return apperrors.Wrap(50001, "查询用户失败", err)
	}

	if req.DisplayName != "" {
		user.DisplayName = req.DisplayName
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}
	if req.Email != "" {
		// Check email uniqueness
		existing, err := s.userRepo.FindByEmail(ctx, req.Email)
		if err == nil && existing.ID != userID {
			return apperrors.New(40900, "邮箱已存在")
		}
		user.Email = req.Email
	}

	return s.userRepo.Update(ctx, user)
}

func (s *AuthService) ChangePassword(ctx context.Context, userID string, req *ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrNotFound
		}
		return apperrors.Wrap(50001, "查询用户失败", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		return apperrors.New(40103, "旧密码错误")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return apperrors.Wrap(50000, "密码加密失败", err)
	}

	user.PasswordHash = string(hashedPassword)
	return s.userRepo.Update(ctx, user)
}

func (s *AuthService) ListUsers(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.userRepo.List(ctx, page, pageSize)
}

func (s *AuthService) generateTokenPair(user *model.User, roles []model.UserRole) (*TokenResponse, error) {
	roleNames := make([]string, 0, len(roles))
	for _, r := range roles {
		roleNames = append(roleNames, r.Role)
	}

	expireHours := s.cfg.JWT.ExpireHours
	if expireHours == 0 {
		expireHours = 24
	}
	expireDuration := time.Duration(expireHours) * time.Hour
	now := time.Now()

	claims := &middleware.JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Roles:    roleNames,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expireDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    s.cfg.JWT.Issuer,
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return nil, apperrors.Wrap(50000, "生成访问令牌失败", err)
	}

	refreshToken := uuid.New().String()

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(expireDuration.Seconds()),
		TokenType:    "Bearer",
	}, nil
}
