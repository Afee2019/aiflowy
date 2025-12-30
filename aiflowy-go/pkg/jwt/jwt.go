package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/aiflowy/aiflowy-go/internal/config"
)

// Claims represents the JWT claims structure
type Claims struct {
	UserID       string `json:"userId"`
	LoginName    string `json:"loginName"`
	Nickname     string `json:"nickname"`
	AccountType  int    `json:"accountType"`
	TenantID     string `json:"tenantId"`
	DeptID       string `json:"deptId"`
	jwt.RegisteredClaims
}

var (
	ErrTokenExpired     = errors.New("token已过期")
	ErrTokenInvalid     = errors.New("无效的token")
	ErrTokenNotProvided = errors.New("未提供token")
)

// GenerateToken generates a new JWT token for the user
func GenerateToken(claims *Claims) (string, error) {
	cfg := config.Get()

	// Set expiration time
	expireTime := time.Now().Add(time.Duration(cfg.JWT.Expire) * time.Second)

	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expireTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    cfg.JWT.Issuer,
	}

	// Create token with HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken parses and validates a JWT token
func ParseToken(tokenString string) (*Claims, error) {
	cfg := config.Get()

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalid
		}
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// RefreshToken refreshes an existing token if it's still valid
func RefreshToken(tokenString string) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	// Generate new token with refreshed expiration
	return GenerateToken(&Claims{
		UserID:      claims.UserID,
		LoginName:   claims.LoginName,
		Nickname:    claims.Nickname,
		AccountType: claims.AccountType,
		TenantID:    claims.TenantID,
		DeptID:      claims.DeptID,
	})
}
