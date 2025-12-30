package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aiflowy/aiflowy-go/internal/dto"
	"github.com/aiflowy/aiflowy-go/internal/entity"
	apierrors "github.com/aiflowy/aiflowy-go/internal/errors"
	"github.com/aiflowy/aiflowy-go/internal/repository"
	"github.com/aiflowy/aiflowy-go/pkg/snowflake"
)

// BotService handles bot-related business logic
type BotService struct {
	botRepo   *repository.BotRepository
	modelRepo *repository.ModelRepository
}

// NewBotService creates a new BotService
func NewBotService() *BotService {
	return &BotService{
		botRepo:   repository.GetBotRepository(),
		modelRepo: repository.GetModelRepository(),
	}
}

// ========== Bot Operations ==========

// GetBot retrieves a bot by ID
func (s *BotService) GetBot(ctx context.Context, id int64) (*entity.Bot, error) {
	return s.botRepo.GetBotByID(ctx, id)
}

// GetBotByAlias retrieves a bot by alias
func (s *BotService) GetBotByAlias(ctx context.Context, alias string) (*entity.Bot, error) {
	return s.botRepo.GetBotByAlias(ctx, alias)
}

// GetBotDetail retrieves a bot with related data
func (s *BotService) GetBotDetail(ctx context.Context, id int64) (*dto.BotDetailResponse, error) {
	bot, err := s.botRepo.GetBotByID(ctx, id)
	if err != nil {
		return nil, apierrors.InternalError("获取机器人失败")
	}
	if bot == nil {
		return nil, apierrors.NotFound("机器人不存在")
	}

	resp := &dto.BotDetailResponse{
		ID:          bot.ID,
		Alias:       bot.Alias,
		DeptID:      bot.DeptID,
		TenantID:    bot.TenantID,
		CategoryID:  bot.CategoryID,
		Title:       bot.Title,
		Description: bot.Description,
		Icon:        bot.Icon,
		ModelID:     bot.ModelID,
		Status:      bot.Status,
	}

	// Parse model options
	if bot.ModelOptions != "" {
		var modelOptions interface{}
		if err := json.Unmarshal([]byte(bot.ModelOptions), &modelOptions); err == nil {
			resp.ModelOptions = modelOptions
		}
	}

	// Parse options
	if bot.Options != "" {
		var options interface{}
		if err := json.Unmarshal([]byte(bot.Options), &options); err == nil {
			resp.Options = options
		}
	}

	// Load category
	if bot.CategoryID > 0 {
		category, err := s.botRepo.GetCategoryByID(ctx, bot.CategoryID)
		if err == nil && category != nil {
			resp.Category = &dto.CategoryInfo{
				ID:           category.ID,
				CategoryName: category.CategoryName,
			}
		}
	}

	// Load model info
	if bot.ModelID > 0 {
		model, err := s.modelRepo.GetModelWithProvider(ctx, bot.ModelID)
		if err == nil && model != nil {
			resp.Model = &dto.ModelInfo{
				ID:           model.ID,
				Title:        model.Title,
				ModelName:    model.ModelName,
				ProviderType: model.ProviderType,
			}
		}
	}

	return resp, nil
}

// ListBots lists bots
func (s *BotService) ListBots(ctx context.Context, req *dto.BotListRequest) ([]*entity.Bot, error) {
	return s.botRepo.ListBots(ctx, req)
}

// PageBots returns paginated bots
func (s *BotService) PageBots(ctx context.Context, req *dto.BotListRequest) ([]*entity.Bot, int64, error) {
	return s.botRepo.PageBots(ctx, req)
}

// SaveBot creates or updates a bot
func (s *BotService) SaveBot(ctx context.Context, req *dto.BotSaveRequest, userID, tenantID, deptID int64) (*entity.Bot, error) {
	now := time.Now()

	if req.ID > 0 {
		// Update existing
		existing, err := s.botRepo.GetBotByID(ctx, req.ID)
		if err != nil {
			return nil, apierrors.InternalError("获取机器人失败")
		}
		if existing == nil {
			return nil, apierrors.NotFound("机器人不存在")
		}

		existing.Alias = req.Alias
		existing.CategoryID = req.CategoryID
		existing.Title = req.Title
		existing.Description = req.Description
		existing.Icon = req.Icon
		existing.ModelID = req.ModelID
		existing.ModelOptions = req.ModelOptions
		existing.Status = req.Status
		existing.Options = req.Options
		existing.Modified = now
		existing.ModifiedBy = userID

		if err := s.botRepo.UpdateBot(ctx, existing); err != nil {
			return nil, apierrors.InternalError("更新机器人失败")
		}
		return existing, nil
	}

	// Create new
	bot := &entity.Bot{
		ID:           snowflake.MustGenerateID(),
		Alias:        req.Alias,
		DeptID:       deptID,
		TenantID:     tenantID,
		CategoryID:   req.CategoryID,
		Title:        req.Title,
		Description:  req.Description,
		Icon:         req.Icon,
		ModelID:      req.ModelID,
		ModelOptions: req.ModelOptions,
		Status:       req.Status,
		Options:      req.Options,
		Created:      now,
		CreatedBy:    userID,
		Modified:     now,
		ModifiedBy:   userID,
	}

	if err := s.botRepo.CreateBot(ctx, bot); err != nil {
		return nil, apierrors.InternalError("创建机器人失败")
	}
	return bot, nil
}

// UpdateBot updates a bot
func (s *BotService) UpdateBot(ctx context.Context, req *dto.BotSaveRequest, userID int64) error {
	existing, err := s.botRepo.GetBotByID(ctx, req.ID)
	if err != nil {
		return apierrors.InternalError("获取机器人失败")
	}
	if existing == nil {
		return apierrors.NotFound("机器人不存在")
	}

	existing.Alias = req.Alias
	existing.CategoryID = req.CategoryID
	existing.Title = req.Title
	existing.Description = req.Description
	existing.Icon = req.Icon
	existing.ModelID = req.ModelID
	existing.ModelOptions = req.ModelOptions
	existing.Status = req.Status
	existing.Options = req.Options
	existing.Modified = time.Now()
	existing.ModifiedBy = userID

	if err := s.botRepo.UpdateBot(ctx, existing); err != nil {
		return apierrors.InternalError("更新机器人失败")
	}
	return nil
}

// UpdateLlmOptions updates bot's LLM options
func (s *BotService) UpdateLlmOptions(ctx context.Context, req *dto.BotUpdateLlmOptionsRequest, userID int64) error {
	existing, err := s.botRepo.GetBotByID(ctx, req.ID)
	if err != nil {
		return apierrors.InternalError("获取机器人失败")
	}
	if existing == nil {
		return apierrors.NotFound("机器人不存在")
	}

	if err := s.botRepo.UpdateBotLlmOptions(ctx, req.ID, req.ModelID, req.ModelOptions, userID); err != nil {
		return apierrors.InternalError("更新LLM配置失败")
	}
	return nil
}

// DeleteBot deletes a bot
func (s *BotService) DeleteBot(ctx context.Context, id int64) error {
	// Delete related data first
	if err := s.botRepo.DeleteMessagesByBotID(ctx, id); err != nil {
		return apierrors.InternalError("删除消息失败")
	}
	if err := s.botRepo.DeleteConversationsByBotID(ctx, id); err != nil {
		return apierrors.InternalError("删除会话失败")
	}
	if err := s.botRepo.DeleteBot(ctx, id); err != nil {
		return apierrors.InternalError("删除机器人失败")
	}
	return nil
}

// ========== Category Operations ==========

// GetCategory retrieves a category by ID
func (s *BotService) GetCategory(ctx context.Context, id int64) (*entity.BotCategory, error) {
	return s.botRepo.GetCategoryByID(ctx, id)
}

// ListCategories lists all categories
func (s *BotService) ListCategories(ctx context.Context, req *dto.BotCategoryListRequest) ([]*entity.BotCategory, error) {
	return s.botRepo.ListCategories(ctx, req)
}

// SaveCategory creates or updates a category
func (s *BotService) SaveCategory(ctx context.Context, req *dto.BotCategorySaveRequest, userID int64) (*entity.BotCategory, error) {
	now := time.Now()

	if req.ID > 0 {
		// Update existing
		existing, err := s.botRepo.GetCategoryByID(ctx, req.ID)
		if err != nil {
			return nil, apierrors.InternalError("获取分类失败")
		}
		if existing == nil {
			return nil, apierrors.NotFound("分类不存在")
		}

		existing.CategoryName = req.CategoryName
		existing.SortNo = req.SortNo
		existing.Status = req.Status
		existing.Modified = now
		existing.ModifiedBy = userID

		if err := s.botRepo.UpdateCategory(ctx, existing); err != nil {
			return nil, apierrors.InternalError("更新分类失败")
		}
		return existing, nil
	}

	// Create new
	category := &entity.BotCategory{
		ID:           snowflake.MustGenerateID(),
		CategoryName: req.CategoryName,
		SortNo:       req.SortNo,
		Status:       req.Status,
		Created:      now,
		CreatedBy:    userID,
		Modified:     now,
		ModifiedBy:   userID,
	}

	if err := s.botRepo.CreateCategory(ctx, category); err != nil {
		return nil, apierrors.InternalError("创建分类失败")
	}
	return category, nil
}

// UpdateCategory updates a category
func (s *BotService) UpdateCategory(ctx context.Context, req *dto.BotCategorySaveRequest, userID int64) error {
	existing, err := s.botRepo.GetCategoryByID(ctx, req.ID)
	if err != nil {
		return apierrors.InternalError("获取分类失败")
	}
	if existing == nil {
		return apierrors.NotFound("分类不存在")
	}

	existing.CategoryName = req.CategoryName
	existing.SortNo = req.SortNo
	existing.Status = req.Status
	existing.Modified = time.Now()
	existing.ModifiedBy = userID

	if err := s.botRepo.UpdateCategory(ctx, existing); err != nil {
		return apierrors.InternalError("更新分类失败")
	}
	return nil
}

// DeleteCategory deletes a category
func (s *BotService) DeleteCategory(ctx context.Context, id int64) error {
	// Check if any bots use this category
	count, err := s.botRepo.CountBotsByCategory(ctx, id)
	if err != nil {
		return apierrors.InternalError("检查分类使用情况失败")
	}
	if count > 0 {
		return apierrors.BadRequest("该分类下存在机器人，无法删除")
	}

	if err := s.botRepo.DeleteCategory(ctx, id); err != nil {
		return apierrors.InternalError("删除分类失败")
	}
	return nil
}

// ========== Conversation Operations ==========

// GetConversation retrieves a conversation by ID
func (s *BotService) GetConversation(ctx context.Context, id int64) (*entity.BotConversation, error) {
	return s.botRepo.GetConversationByID(ctx, id)
}

// ListConversations lists conversations
func (s *BotService) ListConversations(ctx context.Context, req *dto.BotConversationListRequest) ([]*entity.BotConversation, error) {
	return s.botRepo.ListConversations(ctx, req)
}

// PageConversations returns paginated conversations
func (s *BotService) PageConversations(ctx context.Context, req *dto.BotConversationListRequest) ([]*entity.BotConversation, int64, error) {
	return s.botRepo.PageConversations(ctx, req)
}

// SaveConversation creates or updates a conversation
func (s *BotService) SaveConversation(ctx context.Context, req *dto.BotConversationSaveRequest, userID int64) (*entity.BotConversation, error) {
	now := time.Now()

	if req.ID > 0 {
		// Update existing
		existing, err := s.botRepo.GetConversationByID(ctx, req.ID)
		if err != nil {
			return nil, apierrors.InternalError("获取会话失败")
		}
		if existing == nil {
			return nil, apierrors.NotFound("会话不存在")
		}

		existing.Title = req.Title
		existing.Modified = now
		existing.ModifiedBy = userID

		if err := s.botRepo.UpdateConversation(ctx, existing); err != nil {
			return nil, apierrors.InternalError("更新会话失败")
		}
		return existing, nil
	}

	// Create new
	accountID := req.AccountID
	if accountID == 0 {
		accountID = userID
	}

	conv := &entity.BotConversation{
		ID:         snowflake.MustGenerateID(),
		Title:      req.Title,
		BotID:      req.BotID,
		AccountID:  accountID,
		Created:    now,
		CreatedBy:  userID,
		Modified:   now,
		ModifiedBy: userID,
	}

	if err := s.botRepo.CreateConversation(ctx, conv); err != nil {
		return nil, apierrors.InternalError("创建会话失败")
	}
	return conv, nil
}

// UpdateConversation updates a conversation
func (s *BotService) UpdateConversation(ctx context.Context, req *dto.BotConversationSaveRequest, userID int64) error {
	existing, err := s.botRepo.GetConversationByID(ctx, req.ID)
	if err != nil {
		return apierrors.InternalError("获取会话失败")
	}
	if existing == nil {
		return apierrors.NotFound("会话不存在")
	}

	existing.Title = req.Title
	existing.Modified = time.Now()
	existing.ModifiedBy = userID

	if err := s.botRepo.UpdateConversation(ctx, existing); err != nil {
		return apierrors.InternalError("更新会话失败")
	}
	return nil
}

// DeleteConversation deletes a conversation
func (s *BotService) DeleteConversation(ctx context.Context, id int64) error {
	// Delete related messages first
	if err := s.botRepo.DeleteMessagesByConversationID(ctx, id); err != nil {
		return apierrors.InternalError("删除消息失败")
	}
	if err := s.botRepo.DeleteConversation(ctx, id); err != nil {
		return apierrors.InternalError("删除会话失败")
	}
	return nil
}

// GenerateConversationID generates a new conversation ID
func (s *BotService) GenerateConversationID() int64 {
	return snowflake.MustGenerateID()
}

// ========== Message Operations ==========

// GetMessage retrieves a message by ID
func (s *BotService) GetMessage(ctx context.Context, id int64) (*entity.BotMessage, error) {
	return s.botRepo.GetMessageByID(ctx, id)
}

// ListMessages lists messages
func (s *BotService) ListMessages(ctx context.Context, req *dto.BotMessageListRequest) ([]*entity.BotMessage, error) {
	return s.botRepo.ListMessages(ctx, req)
}

// PageMessages returns paginated messages
func (s *BotService) PageMessages(ctx context.Context, req *dto.BotMessageListRequest) ([]*entity.BotMessage, int64, error) {
	return s.botRepo.PageMessages(ctx, req)
}

// SaveMessage creates or updates a message
func (s *BotService) SaveMessage(ctx context.Context, req *dto.BotMessageSaveRequest, userID int64) (*entity.BotMessage, error) {
	now := time.Now()

	if req.ID > 0 {
		// Update existing
		existing, err := s.botRepo.GetMessageByID(ctx, req.ID)
		if err != nil {
			return nil, apierrors.InternalError("获取消息失败")
		}
		if existing == nil {
			return nil, apierrors.NotFound("消息不存在")
		}

		existing.Content = req.Content
		existing.Image = req.Image
		existing.Options = req.Options
		existing.Modified = now

		if err := s.botRepo.UpdateMessage(ctx, existing); err != nil {
			return nil, apierrors.InternalError("更新消息失败")
		}
		return existing, nil
	}

	// Create new
	accountID := req.AccountID
	if accountID == 0 {
		accountID = userID
	}

	msg := &entity.BotMessage{
		ID:             snowflake.MustGenerateID(),
		BotID:          req.BotID,
		AccountID:      accountID,
		ConversationID: req.ConversationID,
		Role:           req.Role,
		Content:        req.Content,
		Image:          req.Image,
		Options:        req.Options,
		Created:        now,
		Modified:       now,
	}

	if err := s.botRepo.CreateMessage(ctx, msg); err != nil {
		return nil, apierrors.InternalError("创建消息失败")
	}
	return msg, nil
}

// UpdateMessage updates a message
func (s *BotService) UpdateMessage(ctx context.Context, req *dto.BotMessageSaveRequest, userID int64) error {
	existing, err := s.botRepo.GetMessageByID(ctx, req.ID)
	if err != nil {
		return apierrors.InternalError("获取消息失败")
	}
	if existing == nil {
		return apierrors.NotFound("消息不存在")
	}

	existing.Content = req.Content
	existing.Image = req.Image
	existing.Options = req.Options
	existing.Modified = time.Now()

	if err := s.botRepo.UpdateMessage(ctx, existing); err != nil {
		return apierrors.InternalError("更新消息失败")
	}
	return nil
}

// DeleteMessage deletes a message
func (s *BotService) DeleteMessage(ctx context.Context, id int64) error {
	if err := s.botRepo.DeleteMessage(ctx, id); err != nil {
		return apierrors.InternalError("删除消息失败")
	}
	return nil
}

// GetRecentMessages gets recent messages for context
func (s *BotService) GetRecentMessages(ctx context.Context, conversationID int64, limit int) ([]*entity.BotMessage, error) {
	return s.botRepo.GetRecentMessages(ctx, conversationID, limit)
}
