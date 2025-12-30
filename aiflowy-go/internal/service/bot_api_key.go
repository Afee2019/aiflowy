package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"

	"github.com/aiflowy/aiflowy-go/internal/config"
	"github.com/aiflowy/aiflowy-go/internal/entity"
	"github.com/aiflowy/aiflowy-go/internal/repository"
)

// BotApiKeyService Bot API 密钥服务
type BotApiKeyService struct {
	repo    *repository.BotApiKeyRepository
	botRepo *repository.BotRepository
}

// NewBotApiKeyService 创建 BotApiKeyService
func NewBotApiKeyService() *BotApiKeyService {
	return &BotApiKeyService{
		repo:    repository.NewBotApiKeyRepository(),
		botRepo: repository.NewBotRepository(repository.GetDB()),
	}
}

// getMasterKey 获取主密钥
func (s *BotApiKeyService) getMasterKey() string {
	cfg := config.GetConfig()
	if cfg != nil && cfg.Security.ApiKeyMasterKey != "" {
		return cfg.Security.ApiKeyMasterKey
	}
	// 默认主密钥 (32字节)
	return "Kj9#mP2$nQ4&rT6*uY8@wE1!zX3%vC5^"
}

// GenerateByBotID 根据 BotID 生成 API 密钥
func (s *BotApiKeyService) GenerateByBotID(ctx context.Context, botID, userID int64) (string, error) {
	// 验证 Bot 存在
	bot, err := s.botRepo.GetBotByID(ctx, botID)
	if err != nil {
		return "", fmt.Errorf("查询 Bot 失败: %w", err)
	}
	if bot == nil {
		return "", fmt.Errorf("Bot 不存在")
	}

	masterKey := s.getMasterKey()
	if len(masterKey) != 32 {
		return "", fmt.Errorf("主密钥长度必须为32字节")
	}

	// 生成随机盐 (16字节)
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("生成随机盐失败: %w", err)
	}

	// 创建 AES 加密器
	block, err := aes.NewCipher([]byte(masterKey))
	if err != nil {
		return "", fmt.Errorf("创建加密器失败: %w", err)
	}

	// 将 botID 转为 bytes
	botIDBytes := big.NewInt(botID).Bytes()

	// PKCS7 填充
	blockSize := block.BlockSize()
	padding := blockSize - len(botIDBytes)%blockSize
	padText := make([]byte, len(botIDBytes)+padding)
	copy(padText, botIDBytes)
	for i := len(botIDBytes); i < len(padText); i++ {
		padText[i] = byte(padding)
	}

	// CBC 加密
	mode := cipher.NewCBCEncrypter(block, salt)
	cipherText := make([]byte, len(padText))
	mode.CryptBlocks(cipherText, padText)

	// Base64 编码
	apiKey := base64.StdEncoding.EncodeToString(cipherText)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)

	// 保存到数据库
	key := &entity.BotApiKey{
		ApiKey:    apiKey,
		BotID:     botID,
		Salt:      saltBase64,
		CreatedBy: &userID,
	}
	if err := s.repo.Create(ctx, key); err != nil {
		return "", fmt.Errorf("保存 API 密钥失败: %w", err)
	}

	return apiKey, nil
}

// DecryptApiKey 解密 API 密钥获取 BotID
func (s *BotApiKeyService) DecryptApiKey(ctx context.Context, apiKey string) (int64, error) {
	if apiKey == "" {
		return 0, fmt.Errorf("API 密钥不能为空")
	}

	// 从数据库获取盐值
	key, err := s.repo.GetByApiKey(ctx, apiKey)
	if err != nil {
		return 0, fmt.Errorf("查询 API 密钥失败: %w", err)
	}
	if key == nil {
		return 0, fmt.Errorf("API 密钥不存在")
	}

	// 解码盐值
	salt, err := base64.StdEncoding.DecodeString(key.Salt)
	if err != nil {
		return 0, fmt.Errorf("解码盐值失败: %w", err)
	}

	// 解码 apiKey
	cipherText, err := base64.StdEncoding.DecodeString(apiKey)
	if err != nil {
		return 0, fmt.Errorf("解码 API 密钥失败: %w", err)
	}

	masterKey := s.getMasterKey()

	// 创建 AES 解密器
	block, err := aes.NewCipher([]byte(masterKey))
	if err != nil {
		return 0, fmt.Errorf("创建解密器失败: %w", err)
	}

	// CBC 解密
	mode := cipher.NewCBCDecrypter(block, salt)
	plainText := make([]byte, len(cipherText))
	mode.CryptBlocks(plainText, cipherText)

	// 去除 PKCS7 填充
	if len(plainText) > 0 {
		padding := int(plainText[len(plainText)-1])
		if padding <= len(plainText) {
			plainText = plainText[:len(plainText)-padding]
		}
	}

	// 转回 int64
	botID := new(big.Int).SetBytes(plainText).Int64()
	return botID, nil
}

// Delete 删除 API 密钥
func (s *BotApiKeyService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// ListByBotID 获取 Bot 的所有 API 密钥
func (s *BotApiKeyService) ListByBotID(ctx context.Context, botID int64) ([]*entity.BotApiKey, error) {
	return s.repo.ListByBotID(ctx, botID)
}
