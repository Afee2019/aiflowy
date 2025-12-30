package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/aiflowy/aiflowy-go/internal/entity"
	"github.com/aiflowy/aiflowy-go/pkg/snowflake"
)

// BotApiKeyRepository Bot API 密钥数据访问层
type BotApiKeyRepository struct {
	db *sql.DB
}

// NewBotApiKeyRepository 创建 BotApiKeyRepository
func NewBotApiKeyRepository() *BotApiKeyRepository {
	return &BotApiKeyRepository{db: GetDB()}
}

// Create 创建 Bot API 密钥
func (r *BotApiKeyRepository) Create(ctx context.Context, apiKey *entity.BotApiKey) error {
	if apiKey.ID == 0 {
		id, _ := snowflake.GenerateID()
		apiKey.ID = id
	}
	now := time.Now()
	apiKey.Created = &now

	query := `
		INSERT INTO tb_bot_api_key (id, api_key, bot_id, salt, options, created, created_by)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		apiKey.ID, apiKey.ApiKey, apiKey.BotID, apiKey.Salt,
		apiKey.Options, apiKey.Created, apiKey.CreatedBy,
	)
	return err
}

// Delete 删除 Bot API 密钥
func (r *BotApiKeyRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM tb_bot_api_key WHERE id = ?", id)
	return err
}

// GetByID 根据 ID 获取
func (r *BotApiKeyRepository) GetByID(ctx context.Context, id int64) (*entity.BotApiKey, error) {
	query := `
		SELECT id, api_key, bot_id, salt, options, created, created_by
		FROM tb_bot_api_key WHERE id = ?
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var apiKey entity.BotApiKey
	var options sql.NullString
	var created sql.NullTime
	var createdBy sql.NullInt64

	err := row.Scan(
		&apiKey.ID, &apiKey.ApiKey, &apiKey.BotID, &apiKey.Salt,
		&options, &created, &createdBy,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if options.Valid {
		apiKey.Options = &options.String
	}
	if created.Valid {
		apiKey.Created = &created.Time
	}
	if createdBy.Valid {
		apiKey.CreatedBy = &createdBy.Int64
	}

	return &apiKey, nil
}

// GetByApiKey 根据 API Key 获取
func (r *BotApiKeyRepository) GetByApiKey(ctx context.Context, apiKey string) (*entity.BotApiKey, error) {
	query := `
		SELECT id, api_key, bot_id, salt, options, created, created_by
		FROM tb_bot_api_key WHERE api_key = ?
	`
	row := r.db.QueryRowContext(ctx, query, apiKey)

	var key entity.BotApiKey
	var options sql.NullString
	var created sql.NullTime
	var createdBy sql.NullInt64

	err := row.Scan(
		&key.ID, &key.ApiKey, &key.BotID, &key.Salt,
		&options, &created, &createdBy,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if options.Valid {
		key.Options = &options.String
	}
	if created.Valid {
		key.Created = &created.Time
	}
	if createdBy.Valid {
		key.CreatedBy = &createdBy.Int64
	}

	return &key, nil
}

// ListByBotID 根据 BotID 获取列表
func (r *BotApiKeyRepository) ListByBotID(ctx context.Context, botID int64) ([]*entity.BotApiKey, error) {
	query := `
		SELECT id, api_key, bot_id, options, created, created_by
		FROM tb_bot_api_key WHERE bot_id = ?
		ORDER BY created DESC
	`
	rows, err := r.db.QueryContext(ctx, query, botID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*entity.BotApiKey
	for rows.Next() {
		var apiKey entity.BotApiKey
		var options sql.NullString
		var created sql.NullTime
		var createdBy sql.NullInt64

		err := rows.Scan(
			&apiKey.ID, &apiKey.ApiKey, &apiKey.BotID,
			&options, &created, &createdBy,
		)
		if err != nil {
			continue
		}

		if options.Valid {
			apiKey.Options = &options.String
		}
		if created.Valid {
			apiKey.Created = &created.Time
		}
		if createdBy.Valid {
			apiKey.CreatedBy = &createdBy.Int64
		}

		// 不返回 salt
		list = append(list, &apiKey)
	}

	return list, nil
}
