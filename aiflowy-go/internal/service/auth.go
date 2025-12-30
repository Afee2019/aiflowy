package service

import (
	"context"
	"fmt"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"github.com/aiflowy/aiflowy-go/internal/dto"
	"github.com/aiflowy/aiflowy-go/internal/entity"
	apierrors "github.com/aiflowy/aiflowy-go/internal/errors"
	"github.com/aiflowy/aiflowy-go/internal/repository"
	"github.com/aiflowy/aiflowy-go/pkg/jwt"
	"github.com/aiflowy/aiflowy-go/pkg/logger"
	"go.uber.org/zap"
)

// AuthService handles authentication business logic
type AuthService struct {
	accountRepo *repository.AccountRepository
}

// NewAuthService creates a new AuthService
func NewAuthService() *AuthService {
	return &AuthService{
		accountRepo: repository.NewAccountRepository(),
	}
}

// Login authenticates a user and returns a token
func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// Find account by login name
	account, err := s.accountRepo.FindByLoginName(ctx, req.Account)
	if err != nil {
		logger.Error("Failed to query account", zap.Error(err))
		return nil, apierrors.InternalError("查询用户失败")
	}

	// Check if account exists
	if account == nil {
		return nil, apierrors.New(1, "用户名/密码错误")
	}

	// Check account status
	if !account.IsEnabled() {
		return nil, apierrors.New(1, "账号未启用，请联系管理员")
	}

	// Verify password using BCrypt
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(req.Password))
	if err != nil {
		return nil, apierrors.New(1, "用户名/密码错误")
	}

	// Generate JWT token
	claims := &jwt.Claims{
		UserID:      strconv.FormatInt(account.ID, 10),
		LoginName:   account.LoginName,
		Nickname:    account.Nickname,
		AccountType: account.AccountType,
		TenantID:    strconv.FormatInt(account.TenantID, 10),
		DeptID:      strconv.FormatInt(account.DeptID, 10),
	}

	token, err := jwt.GenerateToken(claims)
	if err != nil {
		logger.Error("Failed to generate token", zap.Error(err))
		return nil, apierrors.InternalError("生成Token失败")
	}

	logger.Info("User logged in successfully",
		zap.Int64("user_id", account.ID),
		zap.String("login_name", account.LoginName),
	)

	return &dto.LoginResponse{
		Token:    token,
		Nickname: account.Nickname,
		Avatar:   account.Avatar,
	}, nil
}

// GetPermissions returns permission list for current user
func (s *AuthService) GetPermissions(ctx context.Context, userID int64) ([]string, error) {
	permissions, err := s.accountRepo.GetPermissionsByAccountID(ctx, userID)
	if err != nil {
		logger.Error("Failed to get permissions",
			zap.Int64("user_id", userID),
			zap.Error(err),
		)
		return nil, apierrors.InternalError("获取权限失败")
	}

	// If no permissions found, return empty list instead of nil
	if permissions == nil {
		permissions = []string{}
	}

	return permissions, nil
}

// GetRoles returns role list for current user
func (s *AuthService) GetRoles(ctx context.Context, userID int64) ([]string, error) {
	roles, err := s.accountRepo.GetRolesByAccountID(ctx, userID)
	if err != nil {
		logger.Error("Failed to get roles",
			zap.Int64("user_id", userID),
			zap.Error(err),
		)
		return nil, apierrors.InternalError("获取角色失败")
	}

	if roles == nil {
		roles = []string{}
	}

	return roles, nil
}

// GetMenus returns menu list for current user
func (s *AuthService) GetMenus(ctx context.Context, userID int64) ([]map[string]interface{}, error) {
	menus, err := s.accountRepo.GetMenusByAccountID(ctx, userID)
	if err != nil {
		logger.Error("Failed to get menus",
			zap.Int64("user_id", userID),
			zap.Error(err),
		)
		return nil, apierrors.InternalError("获取菜单失败")
	}

	if menus == nil {
		menus = []map[string]interface{}{}
	}

	return menus, nil
}

// GetAccountByID returns account by ID
func (s *AuthService) GetAccountByID(ctx context.Context, userID int64) (*entity.SysAccount, error) {
	account, err := s.accountRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	return account, nil
}

// ValidateToken validates a JWT token and returns claims
func (s *AuthService) ValidateToken(tokenString string) (*jwt.Claims, error) {
	claims, err := jwt.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}
	return claims, nil
}
