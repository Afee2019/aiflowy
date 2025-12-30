package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aiflowy/aiflowy-go/internal/entity"
)

// AccountRepository handles database operations for accounts
type AccountRepository struct {
	db *sql.DB
}

// NewAccountRepository creates a new AccountRepository
func NewAccountRepository() *AccountRepository {
	return &AccountRepository{
		db: GetDB(),
	}
}

// FindByLoginName finds an account by login name
func (r *AccountRepository) FindByLoginName(ctx context.Context, loginName string) (*entity.SysAccount, error) {
	query := `
		SELECT id, dept_id, tenant_id, login_name, password, account_type,
		       nickname, mobile, email, avatar, status, remark,
		       created, created_by, modified, modified_by
		FROM tb_sys_account
		WHERE login_name = ?
		LIMIT 1
	`

	account := &entity.SysAccount{}
	err := r.db.QueryRowContext(ctx, query, loginName).Scan(
		&account.ID,
		&account.DeptID,
		&account.TenantID,
		&account.LoginName,
		&account.Password,
		&account.AccountType,
		&account.Nickname,
		&account.Mobile,
		&account.Email,
		&account.Avatar,
		&account.Status,
		&account.Remark,
		&account.Created,
		&account.CreatedBy,
		&account.Modified,
		&account.ModifiedBy,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query account: %w", err)
	}

	return account, nil
}

// FindByID finds an account by ID
func (r *AccountRepository) FindByID(ctx context.Context, id int64) (*entity.SysAccount, error) {
	query := `
		SELECT id, dept_id, tenant_id, login_name, password, account_type,
		       nickname, mobile, email, avatar, status, remark,
		       created, created_by, modified, modified_by
		FROM tb_sys_account
		WHERE id = ?
		LIMIT 1
	`

	account := &entity.SysAccount{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&account.ID,
		&account.DeptID,
		&account.TenantID,
		&account.LoginName,
		&account.Password,
		&account.AccountType,
		&account.Nickname,
		&account.Mobile,
		&account.Email,
		&account.Avatar,
		&account.Status,
		&account.Remark,
		&account.Created,
		&account.CreatedBy,
		&account.Modified,
		&account.ModifiedBy,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query account: %w", err)
	}

	return account, nil
}

// GetPermissionsByAccountID returns permission tags for an account
func (r *AccountRepository) GetPermissionsByAccountID(ctx context.Context, accountID int64) ([]string, error) {
	query := `
		SELECT DISTINCT m.permission_tag
		FROM tb_sys_menu m
		INNER JOIN tb_sys_role_menu rm ON m.id = rm.menu_id
		INNER JOIN tb_sys_account_role ar ON rm.role_id = ar.role_id
		WHERE ar.account_id = ?
		  AND m.permission_tag IS NOT NULL
		  AND m.permission_tag != ''
		  AND m.status = 1
	`

	rows, err := r.db.QueryContext(ctx, query, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to query permissions: %w", err)
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var permission string
		if err := rows.Scan(&permission); err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		permissions = append(permissions, permission)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating permissions: %w", err)
	}

	return permissions, nil
}

// GetRolesByAccountID returns role keys for an account
func (r *AccountRepository) GetRolesByAccountID(ctx context.Context, accountID int64) ([]string, error) {
	query := `
		SELECT r.role_key
		FROM tb_sys_role r
		INNER JOIN tb_sys_account_role ar ON r.id = ar.role_id
		WHERE ar.account_id = ? AND r.status = 1
	`

	rows, err := r.db.QueryContext(ctx, query, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to query roles: %w", err)
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating roles: %w", err)
	}

	return roles, nil
}

// GetMenusByAccountID returns menus for an account
func (r *AccountRepository) GetMenusByAccountID(ctx context.Context, accountID int64) ([]map[string]interface{}, error) {
	query := `
		SELECT DISTINCT m.id, m.parent_id, m.menu_type, m.menu_title,
		       m.menu_url, m.component, m.menu_icon, m.is_show,
		       m.permission_tag, m.sort_no, m.status
		FROM tb_sys_menu m
		INNER JOIN tb_sys_role_menu rm ON m.id = rm.menu_id
		INNER JOIN tb_sys_account_role ar ON rm.role_id = ar.role_id
		WHERE ar.account_id = ? AND m.status = 1
		ORDER BY m.sort_no
	`

	rows, err := r.db.QueryContext(ctx, query, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to query menus: %w", err)
	}
	defer rows.Close()

	var menus []map[string]interface{}
	for rows.Next() {
		var (
			id            int64
			parentID      int64
			menuType      int
			menuTitle     string
			menuURL       sql.NullString
			component     sql.NullString
			menuIcon      sql.NullString
			isShow        int
			permissionTag sql.NullString
			sortNo        sql.NullInt32
			status        int
		)

		if err := rows.Scan(&id, &parentID, &menuType, &menuTitle, &menuURL, &component, &menuIcon, &isShow, &permissionTag, &sortNo, &status); err != nil {
			return nil, fmt.Errorf("failed to scan menu: %w", err)
		}

		menu := map[string]interface{}{
			"id":            id,
			"parentId":      parentID,
			"menuType":      menuType,
			"menuTitle":     menuTitle,
			"menuUrl":       menuURL.String,
			"component":     component.String,
			"menuIcon":      menuIcon.String,
			"isShow":        isShow,
			"permissionTag": permissionTag.String,
			"sortNo":        sortNo.Int32,
			"status":        status,
		}
		menus = append(menus, menu)
	}

	return menus, nil
}
