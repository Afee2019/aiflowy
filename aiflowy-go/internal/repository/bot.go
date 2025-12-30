package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/aiflowy/aiflowy-go/internal/dto"
	"github.com/aiflowy/aiflowy-go/internal/entity"
)

// BotRepository handles bot database operations
type BotRepository struct {
	db *sql.DB
}

// NewBotRepository creates a new BotRepository
func NewBotRepository(db *sql.DB) *BotRepository {
	return &BotRepository{db: db}
}

// ========== Bot Operations ==========

// GetBotByID retrieves a bot by ID
func (r *BotRepository) GetBotByID(ctx context.Context, id int64) (*entity.Bot, error) {
	query := `SELECT id, COALESCE(alias,''), dept_id, tenant_id, COALESCE(category_id,0), COALESCE(title,''),
		COALESCE(description,''), COALESCE(icon,''), COALESCE(model_id,0), COALESCE(model_options,''),
		COALESCE(status,0), COALESCE(options,''), created, COALESCE(created_by,0), modified, COALESCE(modified_by,0)
		FROM tb_bot WHERE id = ?`

	var b entity.Bot
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&b.ID, &b.Alias, &b.DeptID, &b.TenantID, &b.CategoryID, &b.Title,
		&b.Description, &b.Icon, &b.ModelID, &b.ModelOptions,
		&b.Status, &b.Options, &b.Created, &b.CreatedBy, &b.Modified, &b.ModifiedBy,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// GetBotByAlias retrieves a bot by alias
func (r *BotRepository) GetBotByAlias(ctx context.Context, alias string) (*entity.Bot, error) {
	query := `SELECT id, COALESCE(alias,''), dept_id, tenant_id, COALESCE(category_id,0), COALESCE(title,''),
		COALESCE(description,''), COALESCE(icon,''), COALESCE(model_id,0), COALESCE(model_options,''),
		COALESCE(status,0), COALESCE(options,''), created, COALESCE(created_by,0), modified, COALESCE(modified_by,0)
		FROM tb_bot WHERE alias = ?`

	var b entity.Bot
	err := r.db.QueryRowContext(ctx, query, alias).Scan(
		&b.ID, &b.Alias, &b.DeptID, &b.TenantID, &b.CategoryID, &b.Title,
		&b.Description, &b.Icon, &b.ModelID, &b.ModelOptions,
		&b.Status, &b.Options, &b.Created, &b.CreatedBy, &b.Modified, &b.ModifiedBy,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// ListBots lists bots with optional filters
func (r *BotRepository) ListBots(ctx context.Context, req *dto.BotListRequest) ([]*entity.Bot, error) {
	query := `SELECT id, COALESCE(alias,''), dept_id, tenant_id, COALESCE(category_id,0), COALESCE(title,''),
		COALESCE(description,''), COALESCE(icon,''), COALESCE(model_id,0), COALESCE(model_options,''),
		COALESCE(status,0), COALESCE(options,''), created, COALESCE(created_by,0), modified, COALESCE(modified_by,0)
		FROM tb_bot WHERE 1=1`
	var args []interface{}

	if req != nil {
		if req.CategoryID > 0 {
			query += " AND category_id = ?"
			args = append(args, req.CategoryID)
		}
		if req.Title != "" {
			query += " AND title LIKE ?"
			args = append(args, "%"+req.Title+"%")
		}
		if req.Status != nil {
			query += " AND status = ?"
			args = append(args, *req.Status)
		}
	}

	query += " ORDER BY created DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bots []*entity.Bot
	for rows.Next() {
		var b entity.Bot
		err := rows.Scan(
			&b.ID, &b.Alias, &b.DeptID, &b.TenantID, &b.CategoryID, &b.Title,
			&b.Description, &b.Icon, &b.ModelID, &b.ModelOptions,
			&b.Status, &b.Options, &b.Created, &b.CreatedBy, &b.Modified, &b.ModifiedBy,
		)
		if err != nil {
			return nil, err
		}
		bots = append(bots, &b)
	}
	return bots, nil
}

// PageBots returns paginated bots
func (r *BotRepository) PageBots(ctx context.Context, req *dto.BotListRequest) ([]*entity.Bot, int64, error) {
	countQuery := "SELECT COUNT(*) FROM tb_bot WHERE 1=1"
	query := `SELECT id, COALESCE(alias,''), dept_id, tenant_id, COALESCE(category_id,0), COALESCE(title,''),
		COALESCE(description,''), COALESCE(icon,''), COALESCE(model_id,0), COALESCE(model_options,''),
		COALESCE(status,0), COALESCE(options,''), created, COALESCE(created_by,0), modified, COALESCE(modified_by,0)
		FROM tb_bot WHERE 1=1`
	var args []interface{}

	if req != nil {
		if req.CategoryID > 0 {
			countQuery += " AND category_id = ?"
			query += " AND category_id = ?"
			args = append(args, req.CategoryID)
		}
		if req.Title != "" {
			countQuery += " AND title LIKE ?"
			query += " AND title LIKE ?"
			args = append(args, "%"+req.Title+"%")
		}
		if req.Status != nil {
			countQuery += " AND status = ?"
			query += " AND status = ?"
			args = append(args, *req.Status)
		}
	}

	var total int64
	countArgs := make([]interface{}, len(args))
	copy(countArgs, args)
	err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query += " ORDER BY created DESC LIMIT ? OFFSET ?"
	args = append(args, req.GetPageSize(), req.GetOffset())

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var bots []*entity.Bot
	for rows.Next() {
		var b entity.Bot
		err := rows.Scan(
			&b.ID, &b.Alias, &b.DeptID, &b.TenantID, &b.CategoryID, &b.Title,
			&b.Description, &b.Icon, &b.ModelID, &b.ModelOptions,
			&b.Status, &b.Options, &b.Created, &b.CreatedBy, &b.Modified, &b.ModifiedBy,
		)
		if err != nil {
			return nil, 0, err
		}
		bots = append(bots, &b)
	}
	return bots, total, nil
}

// CreateBot creates a new bot
func (r *BotRepository) CreateBot(ctx context.Context, b *entity.Bot) error {
	query := `INSERT INTO tb_bot
		(id, alias, dept_id, tenant_id, category_id, title, description, icon, model_id, model_options,
		status, options, created, created_by, modified, modified_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	var categoryID, modelID interface{} = nil, nil
	if b.CategoryID > 0 {
		categoryID = b.CategoryID
	}
	if b.ModelID > 0 {
		modelID = b.ModelID
	}

	_, err := r.db.ExecContext(ctx, query,
		b.ID, b.Alias, b.DeptID, b.TenantID, categoryID, b.Title, b.Description, b.Icon, modelID, b.ModelOptions,
		b.Status, b.Options, b.Created, b.CreatedBy, b.Modified, b.ModifiedBy,
	)
	return err
}

// UpdateBot updates an existing bot
func (r *BotRepository) UpdateBot(ctx context.Context, b *entity.Bot) error {
	query := `UPDATE tb_bot SET
		alias = ?, category_id = ?, title = ?, description = ?, icon = ?, model_id = ?, model_options = ?,
		status = ?, options = ?, modified = ?, modified_by = ?
		WHERE id = ?`

	var categoryID, modelID interface{} = nil, nil
	if b.CategoryID > 0 {
		categoryID = b.CategoryID
	}
	if b.ModelID > 0 {
		modelID = b.ModelID
	}

	_, err := r.db.ExecContext(ctx, query,
		b.Alias, categoryID, b.Title, b.Description, b.Icon, modelID, b.ModelOptions,
		b.Status, b.Options, b.Modified, b.ModifiedBy, b.ID,
	)
	return err
}

// UpdateBotLlmOptions updates bot's LLM options
func (r *BotRepository) UpdateBotLlmOptions(ctx context.Context, id, modelID int64, modelOptions string, modifiedBy int64) error {
	query := `UPDATE tb_bot SET model_id = ?, model_options = ?, modified = ?, modified_by = ? WHERE id = ?`

	var mID interface{} = nil
	if modelID > 0 {
		mID = modelID
	}

	_, err := r.db.ExecContext(ctx, query, mID, modelOptions, time.Now(), modifiedBy, id)
	return err
}

// DeleteBot deletes a bot
func (r *BotRepository) DeleteBot(ctx context.Context, id int64) error {
	query := "DELETE FROM tb_bot WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// ========== Bot Category Operations ==========

// GetCategoryByID retrieves a category by ID
func (r *BotRepository) GetCategoryByID(ctx context.Context, id int64) (*entity.BotCategory, error) {
	query := `SELECT id, category_name, COALESCE(sort_no,0), status, created, created_by, modified, modified_by
		FROM tb_bot_category WHERE id = ?`

	var c entity.BotCategory
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.CategoryName, &c.SortNo, &c.Status, &c.Created, &c.CreatedBy, &c.Modified, &c.ModifiedBy,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// ListCategories lists all bot categories
func (r *BotRepository) ListCategories(ctx context.Context, req *dto.BotCategoryListRequest) ([]*entity.BotCategory, error) {
	query := `SELECT id, category_name, COALESCE(sort_no,0), status, created, created_by, modified, modified_by
		FROM tb_bot_category WHERE 1=1`
	var args []interface{}

	if req != nil {
		if req.CategoryName != "" {
			query += " AND category_name LIKE ?"
			args = append(args, "%"+req.CategoryName+"%")
		}
		if req.Status != nil {
			query += " AND status = ?"
			args = append(args, *req.Status)
		}
	}

	query += " ORDER BY sort_no ASC, created DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*entity.BotCategory
	for rows.Next() {
		var c entity.BotCategory
		err := rows.Scan(
			&c.ID, &c.CategoryName, &c.SortNo, &c.Status, &c.Created, &c.CreatedBy, &c.Modified, &c.ModifiedBy,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &c)
	}
	return categories, nil
}

// CreateCategory creates a new category
func (r *BotRepository) CreateCategory(ctx context.Context, c *entity.BotCategory) error {
	query := `INSERT INTO tb_bot_category
		(id, category_name, sort_no, status, created, created_by, modified, modified_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		c.ID, c.CategoryName, c.SortNo, c.Status, c.Created, c.CreatedBy, c.Modified, c.ModifiedBy,
	)
	return err
}

// UpdateCategory updates an existing category
func (r *BotRepository) UpdateCategory(ctx context.Context, c *entity.BotCategory) error {
	query := `UPDATE tb_bot_category SET
		category_name = ?, sort_no = ?, status = ?, modified = ?, modified_by = ?
		WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		c.CategoryName, c.SortNo, c.Status, c.Modified, c.ModifiedBy, c.ID,
	)
	return err
}

// DeleteCategory deletes a category
func (r *BotRepository) DeleteCategory(ctx context.Context, id int64) error {
	query := "DELETE FROM tb_bot_category WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// CountBotsByCategory counts bots in a category
func (r *BotRepository) CountBotsByCategory(ctx context.Context, categoryID int64) (int64, error) {
	var count int64
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM tb_bot WHERE category_id = ?", categoryID).Scan(&count)
	return count, err
}

// ========== Conversation Operations ==========

// GetConversationByID retrieves a conversation by ID
func (r *BotRepository) GetConversationByID(ctx context.Context, id int64) (*entity.BotConversation, error) {
	query := `SELECT id, title, COALESCE(bot_id,0), COALESCE(account_id,0), created, COALESCE(created_by,0), modified, COALESCE(modified_by,0)
		FROM tb_bot_conversation WHERE id = ?`

	var c entity.BotConversation
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.Title, &c.BotID, &c.AccountID, &c.Created, &c.CreatedBy, &c.Modified, &c.ModifiedBy,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// ListConversations lists conversations with optional filters
func (r *BotRepository) ListConversations(ctx context.Context, req *dto.BotConversationListRequest) ([]*entity.BotConversation, error) {
	query := `SELECT id, title, COALESCE(bot_id,0), COALESCE(account_id,0), created, COALESCE(created_by,0), modified, COALESCE(modified_by,0)
		FROM tb_bot_conversation WHERE 1=1`
	var args []interface{}

	if req != nil {
		if req.BotID > 0 {
			query += " AND bot_id = ?"
			args = append(args, req.BotID)
		}
		if req.AccountID > 0 {
			query += " AND account_id = ?"
			args = append(args, req.AccountID)
		}
	}

	query += " ORDER BY modified DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []*entity.BotConversation
	for rows.Next() {
		var c entity.BotConversation
		err := rows.Scan(
			&c.ID, &c.Title, &c.BotID, &c.AccountID, &c.Created, &c.CreatedBy, &c.Modified, &c.ModifiedBy,
		)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, &c)
	}
	return conversations, nil
}

// PageConversations returns paginated conversations
func (r *BotRepository) PageConversations(ctx context.Context, req *dto.BotConversationListRequest) ([]*entity.BotConversation, int64, error) {
	countQuery := "SELECT COUNT(*) FROM tb_bot_conversation WHERE 1=1"
	query := `SELECT id, title, COALESCE(bot_id,0), COALESCE(account_id,0), created, COALESCE(created_by,0), modified, COALESCE(modified_by,0)
		FROM tb_bot_conversation WHERE 1=1`
	var args []interface{}

	if req != nil {
		if req.BotID > 0 {
			countQuery += " AND bot_id = ?"
			query += " AND bot_id = ?"
			args = append(args, req.BotID)
		}
		if req.AccountID > 0 {
			countQuery += " AND account_id = ?"
			query += " AND account_id = ?"
			args = append(args, req.AccountID)
		}
	}

	var total int64
	countArgs := make([]interface{}, len(args))
	copy(countArgs, args)
	err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query += " ORDER BY modified DESC LIMIT ? OFFSET ?"
	args = append(args, req.GetPageSize(), req.GetOffset())

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var conversations []*entity.BotConversation
	for rows.Next() {
		var c entity.BotConversation
		err := rows.Scan(
			&c.ID, &c.Title, &c.BotID, &c.AccountID, &c.Created, &c.CreatedBy, &c.Modified, &c.ModifiedBy,
		)
		if err != nil {
			return nil, 0, err
		}
		conversations = append(conversations, &c)
	}
	return conversations, total, nil
}

// CreateConversation creates a new conversation
func (r *BotRepository) CreateConversation(ctx context.Context, c *entity.BotConversation) error {
	query := `INSERT INTO tb_bot_conversation
		(id, title, bot_id, account_id, created, created_by, modified, modified_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	var botID, accountID interface{} = nil, nil
	if c.BotID > 0 {
		botID = c.BotID
	}
	if c.AccountID > 0 {
		accountID = c.AccountID
	}

	_, err := r.db.ExecContext(ctx, query,
		c.ID, c.Title, botID, accountID, c.Created, c.CreatedBy, c.Modified, c.ModifiedBy,
	)
	return err
}

// UpdateConversation updates an existing conversation
func (r *BotRepository) UpdateConversation(ctx context.Context, c *entity.BotConversation) error {
	query := `UPDATE tb_bot_conversation SET title = ?, modified = ?, modified_by = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, c.Title, c.Modified, c.ModifiedBy, c.ID)
	return err
}

// DeleteConversation deletes a conversation
func (r *BotRepository) DeleteConversation(ctx context.Context, id int64) error {
	query := "DELETE FROM tb_bot_conversation WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// DeleteConversationsByBotID deletes all conversations for a bot
func (r *BotRepository) DeleteConversationsByBotID(ctx context.Context, botID int64) error {
	query := "DELETE FROM tb_bot_conversation WHERE bot_id = ?"
	_, err := r.db.ExecContext(ctx, query, botID)
	return err
}

// ========== Message Operations ==========

// GetMessageByID retrieves a message by ID
func (r *BotRepository) GetMessageByID(ctx context.Context, id int64) (*entity.BotMessage, error) {
	query := `SELECT id, COALESCE(bot_id,0), COALESCE(account_id,0), COALESCE(conversation_id,0), COALESCE(role,''),
		COALESCE(content,''), COALESCE(image,''), COALESCE(options,''), created, modified
		FROM tb_bot_message WHERE id = ?`

	var m entity.BotMessage
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID, &m.BotID, &m.AccountID, &m.ConversationID, &m.Role,
		&m.Content, &m.Image, &m.Options, &m.Created, &m.Modified,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// ListMessages lists messages with optional filters
func (r *BotRepository) ListMessages(ctx context.Context, req *dto.BotMessageListRequest) ([]*entity.BotMessage, error) {
	query := `SELECT id, COALESCE(bot_id,0), COALESCE(account_id,0), COALESCE(conversation_id,0), COALESCE(role,''),
		COALESCE(content,''), COALESCE(image,''), COALESCE(options,''), created, modified
		FROM tb_bot_message WHERE 1=1`
	var args []interface{}

	if req != nil {
		if req.BotID > 0 {
			query += " AND bot_id = ?"
			args = append(args, req.BotID)
		}
		if req.ConversationID > 0 {
			query += " AND conversation_id = ?"
			args = append(args, req.ConversationID)
		}
		if req.AccountID > 0 {
			query += " AND account_id = ?"
			args = append(args, req.AccountID)
		}
	}

	query += " ORDER BY created ASC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*entity.BotMessage
	for rows.Next() {
		var m entity.BotMessage
		err := rows.Scan(
			&m.ID, &m.BotID, &m.AccountID, &m.ConversationID, &m.Role,
			&m.Content, &m.Image, &m.Options, &m.Created, &m.Modified,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &m)
	}
	return messages, nil
}

// PageMessages returns paginated messages
func (r *BotRepository) PageMessages(ctx context.Context, req *dto.BotMessageListRequest) ([]*entity.BotMessage, int64, error) {
	countQuery := "SELECT COUNT(*) FROM tb_bot_message WHERE 1=1"
	query := `SELECT id, COALESCE(bot_id,0), COALESCE(account_id,0), COALESCE(conversation_id,0), COALESCE(role,''),
		COALESCE(content,''), COALESCE(image,''), COALESCE(options,''), created, modified
		FROM tb_bot_message WHERE 1=1`
	var args []interface{}

	if req != nil {
		if req.BotID > 0 {
			countQuery += " AND bot_id = ?"
			query += " AND bot_id = ?"
			args = append(args, req.BotID)
		}
		if req.ConversationID > 0 {
			countQuery += " AND conversation_id = ?"
			query += " AND conversation_id = ?"
			args = append(args, req.ConversationID)
		}
		if req.AccountID > 0 {
			countQuery += " AND account_id = ?"
			query += " AND account_id = ?"
			args = append(args, req.AccountID)
		}
	}

	var total int64
	countArgs := make([]interface{}, len(args))
	copy(countArgs, args)
	err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query += " ORDER BY created ASC LIMIT ? OFFSET ?"
	args = append(args, req.GetPageSize(), req.GetOffset())

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var messages []*entity.BotMessage
	for rows.Next() {
		var m entity.BotMessage
		err := rows.Scan(
			&m.ID, &m.BotID, &m.AccountID, &m.ConversationID, &m.Role,
			&m.Content, &m.Image, &m.Options, &m.Created, &m.Modified,
		)
		if err != nil {
			return nil, 0, err
		}
		messages = append(messages, &m)
	}
	return messages, total, nil
}

// CreateMessage creates a new message
func (r *BotRepository) CreateMessage(ctx context.Context, m *entity.BotMessage) error {
	query := `INSERT INTO tb_bot_message
		(id, bot_id, account_id, conversation_id, role, content, image, options, created, modified)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	var botID, accountID, conversationID interface{} = nil, nil, nil
	if m.BotID > 0 {
		botID = m.BotID
	}
	if m.AccountID > 0 {
		accountID = m.AccountID
	}
	if m.ConversationID > 0 {
		conversationID = m.ConversationID
	}

	_, err := r.db.ExecContext(ctx, query,
		m.ID, botID, accountID, conversationID, m.Role, m.Content, m.Image, m.Options, m.Created, m.Modified,
	)
	return err
}

// UpdateMessage updates an existing message
func (r *BotRepository) UpdateMessage(ctx context.Context, m *entity.BotMessage) error {
	query := `UPDATE tb_bot_message SET content = ?, image = ?, options = ?, modified = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, m.Content, m.Image, m.Options, m.Modified, m.ID)
	return err
}

// DeleteMessage deletes a message
func (r *BotRepository) DeleteMessage(ctx context.Context, id int64) error {
	query := "DELETE FROM tb_bot_message WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// DeleteMessagesByConversationID deletes all messages for a conversation
func (r *BotRepository) DeleteMessagesByConversationID(ctx context.Context, conversationID int64) error {
	query := "DELETE FROM tb_bot_message WHERE conversation_id = ?"
	_, err := r.db.ExecContext(ctx, query, conversationID)
	return err
}

// DeleteMessagesByBotID deletes all messages for a bot
func (r *BotRepository) DeleteMessagesByBotID(ctx context.Context, botID int64) error {
	query := "DELETE FROM tb_bot_message WHERE bot_id = ?"
	_, err := r.db.ExecContext(ctx, query, botID)
	return err
}

// GetRecentMessages gets recent messages for a conversation (for context)
func (r *BotRepository) GetRecentMessages(ctx context.Context, conversationID int64, limit int) ([]*entity.BotMessage, error) {
	query := `SELECT id, COALESCE(bot_id,0), COALESCE(account_id,0), COALESCE(conversation_id,0), COALESCE(role,''),
		COALESCE(content,''), COALESCE(image,''), COALESCE(options,''), created, modified
		FROM tb_bot_message WHERE conversation_id = ?
		ORDER BY created DESC LIMIT ?`

	rows, err := r.db.QueryContext(ctx, query, conversationID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*entity.BotMessage
	for rows.Next() {
		var m entity.BotMessage
		err := rows.Scan(
			&m.ID, &m.BotID, &m.AccountID, &m.ConversationID, &m.Role,
			&m.Content, &m.Image, &m.Options, &m.Created, &m.Modified,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &m)
	}

	// Reverse the order (oldest first)
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// ========== Singleton ==========

var botRepo *BotRepository
var botRepoInit = false

// GetBotRepository returns the singleton BotRepository
func GetBotRepository() *BotRepository {
	if !botRepoInit {
		botRepo = NewBotRepository(GetDB())
		botRepoInit = true
	}
	return botRepo
}

// helper for building IN clauses
func buildInClause(ids []int64) (string, []interface{}) {
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	return strings.Join(placeholders, ","), args
}

// unused but keep for potential future use
var _ = fmt.Sprintf
var _ = buildInClause
