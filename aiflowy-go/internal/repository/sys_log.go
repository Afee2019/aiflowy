package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/aiflowy/aiflowy-go/internal/entity"
	"github.com/aiflowy/aiflowy-go/pkg/snowflake"
)

// SysLogRepository 操作日志数据访问层
type SysLogRepository struct {
	db *sql.DB
}

// NewSysLogRepository 创建 SysLogRepository
func NewSysLogRepository() *SysLogRepository {
	return &SysLogRepository{db: GetDB()}
}

// Create 创建日志
func (r *SysLogRepository) Create(ctx context.Context, log *entity.SysLog) error {
	if log.ID == 0 {
		id, _ := snowflake.GenerateID()
		log.ID = id
	}
	now := time.Now()
	log.Created = &now

	query := `
		INSERT INTO tb_sys_log (id, account_id, action_name, action_type, action_class, action_method,
			action_url, action_ip, action_params, action_body, status, created)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		log.ID, log.AccountID, log.ActionName, log.ActionType,
		log.ActionClass, log.ActionMethod, log.ActionURL, log.ActionIP,
		log.ActionParams, log.ActionBody, log.Status, log.Created,
	)
	return err
}

// Page 分页查询
func (r *SysLogRepository) Page(ctx context.Context, pageNum, pageSize int, actionName, actionType string) ([]*entity.SysLog, int64, error) {
	// 构建查询条件
	countQuery := `SELECT COUNT(*) FROM tb_sys_log WHERE 1=1`
	query := `
		SELECT l.id, l.account_id, l.action_name, l.action_type, l.action_class, l.action_method,
			l.action_url, l.action_ip, l.action_params, l.action_body, l.status, l.created,
			a.id, a.login_name, a.nickname
		FROM tb_sys_log l
		LEFT JOIN tb_sys_account a ON l.account_id = a.id
		WHERE 1=1
	`
	args := []interface{}{}
	countArgs := []interface{}{}

	if actionName != "" {
		countQuery += " AND action_name LIKE ?"
		query += " AND l.action_name LIKE ?"
		args = append(args, "%"+actionName+"%")
		countArgs = append(countArgs, "%"+actionName+"%")
	}
	if actionType != "" {
		countQuery += " AND action_type = ?"
		query += " AND l.action_type = ?"
		args = append(args, actionType)
		countArgs = append(countArgs, actionType)
	}

	// 查询总数
	var total int64
	r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)

	// 添加排序和分页
	query += " ORDER BY l.created DESC LIMIT ? OFFSET ?"
	offset := (pageNum - 1) * pageSize
	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*entity.SysLog
	for rows.Next() {
		var log entity.SysLog
		var accountID sql.NullInt64
		var actionName, actionType, actionClass, actionMethod sql.NullString
		var actionURL, actionIP, actionParams, actionBody sql.NullString
		var created sql.NullTime
		var accID sql.NullInt64
		var accLoginName, accNickname sql.NullString

		err := rows.Scan(
			&log.ID, &accountID, &actionName, &actionType, &actionClass, &actionMethod,
			&actionURL, &actionIP, &actionParams, &actionBody, &log.Status, &created,
			&accID, &accLoginName, &accNickname,
		)
		if err != nil {
			continue
		}

		if accountID.Valid {
			log.AccountID = &accountID.Int64
		}
		if actionName.Valid {
			log.ActionName = actionName.String
		}
		if actionType.Valid {
			log.ActionType = actionType.String
		}
		if actionClass.Valid {
			log.ActionClass = actionClass.String
		}
		if actionMethod.Valid {
			log.ActionMethod = actionMethod.String
		}
		if actionURL.Valid {
			log.ActionURL = actionURL.String
		}
		if actionIP.Valid {
			log.ActionIP = actionIP.String
		}
		if actionParams.Valid {
			log.ActionParams = actionParams.String
		}
		if actionBody.Valid {
			log.ActionBody = actionBody.String
		}
		if created.Valid {
			log.Created = &created.Time
		}

		// 关联用户信息
		if accID.Valid {
			log.Account = &entity.SysAccount{
				ID: accID.Int64,
			}
			if accLoginName.Valid {
				log.Account.LoginName = accLoginName.String
			}
			if accNickname.Valid {
				log.Account.Nickname = accNickname.String
			}
		}

		list = append(list, &log)
	}

	return list, total, nil
}

// Delete 删除日志
func (r *SysLogRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM tb_sys_log WHERE id = ?", id)
	return err
}

// DeleteBefore 删除指定日期之前的日志
func (r *SysLogRepository) DeleteBefore(ctx context.Context, before time.Time) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM tb_sys_log WHERE created < ?", before)
	return err
}
