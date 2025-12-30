package jwt

import (
	"testing"
	"time"

	"github.com/aiflowy/aiflowy-go/internal/config"
)

// initTestConfig initializes a test configuration
func initTestConfig() {
	// Create a minimal test config file path
	// For testing, we'll use Load with a temp config
}

// setupTestConfig sets up the config for testing
func setupTestConfig() {
	// Load config from the test configs directory
	cfg, err := config.Load("../../configs/config.dev.yaml")
	if err != nil {
		// If config file doesn't exist, manually set defaults
		// This is a workaround for testing
		panic("config file not found, please ensure configs/config.dev.yaml exists: " + err.Error())
	}
	_ = cfg
}

func init() {
	// Try to load config, ignore errors for CI environments
	_, _ = config.Load("../../configs/config.dev.yaml")
}

func TestGenerateToken(t *testing.T) {
	cfg := config.Get()
	if cfg == nil {
		t.Skip("config not loaded, skipping test")
	}

	claims := &Claims{
		UserID:      "123",
		LoginName:   "testuser",
		Nickname:    "Test User",
		AccountType: 1,
		TenantID:    "1000000",
		DeptID:      "1",
	}

	token, err := GenerateToken(claims)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("expected non-empty token")
	}
}

func TestParseToken(t *testing.T) {
	cfg := config.Get()
	if cfg == nil {
		t.Skip("config not loaded, skipping test")
	}

	// Generate a token first
	originalClaims := &Claims{
		UserID:      "456",
		LoginName:   "parsetest",
		Nickname:    "Parse Test",
		AccountType: 2,
		TenantID:    "2000000",
		DeptID:      "2",
	}

	token, err := GenerateToken(originalClaims)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	// Parse the token
	parsedClaims, err := ParseToken(token)
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}

	// Verify claims
	if parsedClaims.UserID != originalClaims.UserID {
		t.Errorf("expected UserID '%s', got '%s'", originalClaims.UserID, parsedClaims.UserID)
	}
	if parsedClaims.LoginName != originalClaims.LoginName {
		t.Errorf("expected LoginName '%s', got '%s'", originalClaims.LoginName, parsedClaims.LoginName)
	}
	if parsedClaims.Nickname != originalClaims.Nickname {
		t.Errorf("expected Nickname '%s', got '%s'", originalClaims.Nickname, parsedClaims.Nickname)
	}
	if parsedClaims.AccountType != originalClaims.AccountType {
		t.Errorf("expected AccountType %d, got %d", originalClaims.AccountType, parsedClaims.AccountType)
	}
	if parsedClaims.TenantID != originalClaims.TenantID {
		t.Errorf("expected TenantID '%s', got '%s'", originalClaims.TenantID, parsedClaims.TenantID)
	}
	if parsedClaims.DeptID != originalClaims.DeptID {
		t.Errorf("expected DeptID '%s', got '%s'", originalClaims.DeptID, parsedClaims.DeptID)
	}
}

func TestParseToken_Invalid(t *testing.T) {
	cfg := config.Get()
	if cfg == nil {
		t.Skip("config not loaded, skipping test")
	}

	tests := []struct {
		name  string
		token string
	}{
		{"empty token", ""},
		{"invalid format", "not-a-jwt"},
		{"malformed jwt", "a.b.c"},
		{"tampered jwt", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIxMjMifQ.tamperedSignature"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseToken(tt.token)
			if err == nil {
				t.Error("expected error for invalid token")
			}
		})
	}
}

func TestRefreshToken(t *testing.T) {
	cfg := config.Get()
	if cfg == nil {
		t.Skip("config not loaded, skipping test")
	}

	// Generate a token first
	originalClaims := &Claims{
		UserID:      "789",
		LoginName:   "refreshtest",
		Nickname:    "Refresh Test",
		AccountType: 3,
		TenantID:    "3000000",
		DeptID:      "3",
	}

	originalToken, err := GenerateToken(originalClaims)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	// Wait a bit to ensure different timestamp
	time.Sleep(time.Millisecond * 10)

	// Refresh the token
	newToken, err := RefreshToken(originalToken)
	if err != nil {
		t.Fatalf("failed to refresh token: %v", err)
	}

	if newToken == "" {
		t.Error("expected non-empty refreshed token")
	}

	// Verify the new token has the same claims
	newClaims, err := ParseToken(newToken)
	if err != nil {
		t.Fatalf("failed to parse refreshed token: %v", err)
	}

	if newClaims.UserID != originalClaims.UserID {
		t.Errorf("expected UserID '%s', got '%s'", originalClaims.UserID, newClaims.UserID)
	}
	if newClaims.LoginName != originalClaims.LoginName {
		t.Errorf("expected LoginName '%s', got '%s'", originalClaims.LoginName, newClaims.LoginName)
	}
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	cfg := config.Get()
	if cfg == nil {
		t.Skip("config not loaded, skipping test")
	}

	_, err := RefreshToken("invalid-token")
	if err == nil {
		t.Error("expected error for invalid token")
	}
}

func TestClaims_Fields(t *testing.T) {
	claims := &Claims{
		UserID:      "user123",
		LoginName:   "testlogin",
		Nickname:    "Test Nick",
		AccountType: 99,
		TenantID:    "tenant456",
		DeptID:      "dept789",
	}

	if claims.UserID != "user123" {
		t.Error("UserID field incorrect")
	}
	if claims.LoginName != "testlogin" {
		t.Error("LoginName field incorrect")
	}
	if claims.Nickname != "Test Nick" {
		t.Error("Nickname field incorrect")
	}
	if claims.AccountType != 99 {
		t.Error("AccountType field incorrect")
	}
	if claims.TenantID != "tenant456" {
		t.Error("TenantID field incorrect")
	}
	if claims.DeptID != "dept789" {
		t.Error("DeptID field incorrect")
	}
}

func TestErrorConstants(t *testing.T) {
	if ErrTokenExpired.Error() != "token已过期" {
		t.Errorf("ErrTokenExpired message incorrect: %s", ErrTokenExpired.Error())
	}
	if ErrTokenInvalid.Error() != "无效的token" {
		t.Errorf("ErrTokenInvalid message incorrect: %s", ErrTokenInvalid.Error())
	}
	if ErrTokenNotProvided.Error() != "未提供token" {
		t.Errorf("ErrTokenNotProvided message incorrect: %s", ErrTokenNotProvided.Error())
	}
}

func BenchmarkGenerateToken(b *testing.B) {
	cfg := config.Get()
	if cfg == nil {
		b.Skip("config not loaded, skipping benchmark")
	}

	claims := &Claims{
		UserID:      "bench123",
		LoginName:   "benchuser",
		Nickname:    "Bench User",
		AccountType: 1,
		TenantID:    "1000000",
		DeptID:      "1",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GenerateToken(claims)
	}
}

func BenchmarkParseToken(b *testing.B) {
	cfg := config.Get()
	if cfg == nil {
		b.Skip("config not loaded, skipping benchmark")
	}

	claims := &Claims{
		UserID:      "bench123",
		LoginName:   "benchuser",
		Nickname:    "Bench User",
		AccountType: 1,
		TenantID:    "1000000",
		DeptID:      "1",
	}

	token, err := GenerateToken(claims)
	if err != nil {
		b.Fatalf("failed to generate token: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseToken(token)
	}
}
