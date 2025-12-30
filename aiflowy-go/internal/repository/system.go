package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/aiflowy/aiflowy-go/internal/dto"
	"github.com/aiflowy/aiflowy-go/internal/entity"
	"github.com/aiflowy/aiflowy-go/pkg/snowflake"
)

// SystemRepository handles database operations for system management
type SystemRepository struct {
	db *sql.DB
}

// NewSystemRepository creates a new SystemRepository
func NewSystemRepository() *SystemRepository {
	return &SystemRepository{
		db: GetDB(),
	}
}

// ========== Account Operations ==========

// AccountPage returns paginated accounts
func (r *SystemRepository) AccountPage(ctx context.Context, req *dto.AccountQueryRequest) ([]*entity.SysAccount, int64, error) {
	// Build WHERE clause
	var conditions []string
	var args []interface{}

	if req.LoginName != "" {
		conditions = append(conditions, "login_name LIKE ?")
		args = append(args, "%"+req.LoginName+"%")
	}
	if req.Nickname != "" {
		conditions = append(conditions, "nickname LIKE ?")
		args = append(args, "%"+req.Nickname+"%")
	}
	if req.Status != nil {
		conditions = append(conditions, "status = ?")
		args = append(args, *req.Status)
	}
	if req.DeptID != nil {
		conditions = append(conditions, "dept_id = ?")
		args = append(args, *req.DeptID)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM tb_sys_account %s", whereClause)
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count accounts: %w", err)
	}

	// Query page
	orderBy := "ORDER BY id DESC"
	if req.SortKey != "" {
		sortType := "ASC"
		if strings.ToUpper(req.SortType) == "DESC" {
			sortType = "DESC"
		}
		orderBy = fmt.Sprintf("ORDER BY %s %s", toSnakeCase(req.SortKey), sortType)
	}

	query := fmt.Sprintf(`
		SELECT id, dept_id, tenant_id, login_name, password, account_type,
		       nickname, mobile, email, avatar, status, remark,
		       created, created_by, modified, modified_by
		FROM tb_sys_account %s %s LIMIT ? OFFSET ?
	`, whereClause, orderBy)

	args = append(args, req.GetPageSize(), req.GetOffset())
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query accounts: %w", err)
	}
	defer rows.Close()

	var accounts []*entity.SysAccount
	for rows.Next() {
		a := &entity.SysAccount{}
		if err := rows.Scan(&a.ID, &a.DeptID, &a.TenantID, &a.LoginName, &a.Password,
			&a.AccountType, &a.Nickname, &a.Mobile, &a.Email, &a.Avatar,
			&a.Status, &a.Remark, &a.Created, &a.CreatedBy, &a.Modified, &a.ModifiedBy); err != nil {
			return nil, 0, fmt.Errorf("scan account: %w", err)
		}
		accounts = append(accounts, a)
	}

	return accounts, total, nil
}

// AccountSave creates a new account
func (r *SystemRepository) AccountSave(ctx context.Context, a *entity.SysAccount) error {
	id, _ := snowflake.GenerateID()
	now := time.Now()

	query := `
		INSERT INTO tb_sys_account (id, dept_id, tenant_id, login_name, password, account_type,
		       nickname, mobile, email, avatar, status, remark, created, created_by, modified, modified_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query, id, a.DeptID, a.TenantID, a.LoginName, a.Password,
		a.AccountType, a.Nickname, a.Mobile, a.Email, a.Avatar, a.Status, a.Remark,
		now, a.CreatedBy, now, a.CreatedBy)
	if err != nil {
		return fmt.Errorf("insert account: %w", err)
	}
	a.ID = id
	return nil
}

// AccountUpdate updates an account
func (r *SystemRepository) AccountUpdate(ctx context.Context, a *entity.SysAccount) error {
	query := `
		UPDATE tb_sys_account SET dept_id=?, nickname=?, mobile=?, email=?, avatar=?,
		       status=?, remark=?, modified=?, modified_by=?
		WHERE id=?
	`
	_, err := r.db.ExecContext(ctx, query, a.DeptID, a.Nickname, a.Mobile, a.Email, a.Avatar,
		a.Status, a.Remark, time.Now(), a.ModifiedBy, a.ID)
	return err
}

// AccountDelete deletes an account
func (r *SystemRepository) AccountDelete(ctx context.Context, id int64) error {
	// Delete role associations first
	r.db.ExecContext(ctx, "DELETE FROM tb_sys_account_role WHERE account_id=?", id)
	// Delete account
	_, err := r.db.ExecContext(ctx, "DELETE FROM tb_sys_account WHERE id=?", id)
	return err
}

// AccountExistsByLoginName checks if login name exists
func (r *SystemRepository) AccountExistsByLoginName(ctx context.Context, loginName string, excludeID int64) (bool, error) {
	query := "SELECT COUNT(*) FROM tb_sys_account WHERE login_name=? AND id!=?"
	var count int
	if err := r.db.QueryRowContext(ctx, query, loginName, excludeID).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetAccountRoleIds returns role IDs for an account
func (r *SystemRepository) GetAccountRoleIds(ctx context.Context, accountID int64) ([]int64, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT role_id FROM tb_sys_account_role WHERE account_id=?", accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// SaveAccountRoles saves account-role associations
func (r *SystemRepository) SaveAccountRoles(ctx context.Context, accountID int64, roleIds []int64) error {
	// Delete existing
	r.db.ExecContext(ctx, "DELETE FROM tb_sys_account_role WHERE account_id=?", accountID)
	// Insert new
	for _, roleID := range roleIds {
		id, _ := snowflake.GenerateID()
		r.db.ExecContext(ctx, "INSERT INTO tb_sys_account_role (id, account_id, role_id) VALUES (?, ?, ?)",
			id, accountID, roleID)
	}
	return nil
}

// ========== Role Operations ==========

// RolePage returns paginated roles
func (r *SystemRepository) RolePage(ctx context.Context, req *dto.RoleQueryRequest) ([]*entity.SysRole, int64, error) {
	var conditions []string
	var args []interface{}

	if req.RoleName != "" {
		conditions = append(conditions, "role_name LIKE ?")
		args = append(args, "%"+req.RoleName+"%")
	}
	if req.RoleKey != "" {
		conditions = append(conditions, "role_key LIKE ?")
		args = append(args, "%"+req.RoleKey+"%")
	}
	if req.Status != nil {
		conditions = append(conditions, "status = ?")
		args = append(args, *req.Status)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	var total int64
	r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(*) FROM tb_sys_role %s", whereClause), args...).Scan(&total)

	// Query
	query := fmt.Sprintf(`
		SELECT id, tenant_id, role_name, role_key, status, data_scope,
		       menu_check_strictly, dept_check_strictly, created, created_by, modified, modified_by, remark
		FROM tb_sys_role %s ORDER BY id DESC LIMIT ? OFFSET ?
	`, whereClause)
	args = append(args, req.GetPageSize(), req.GetOffset())

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var roles []*entity.SysRole
	for rows.Next() {
		ro := &entity.SysRole{}
		rows.Scan(&ro.ID, &ro.TenantID, &ro.RoleName, &ro.RoleKey, &ro.Status, &ro.DataScope,
			&ro.MenuCheckStrictly, &ro.DeptCheckStrictly, &ro.Created, &ro.CreatedBy, &ro.Modified, &ro.ModifiedBy, &ro.Remark)
		roles = append(roles, ro)
	}
	return roles, total, nil
}

// RoleList returns all roles
func (r *SystemRepository) RoleList(ctx context.Context) ([]*entity.SysRole, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, role_name, role_key, status, data_scope,
		       menu_check_strictly, dept_check_strictly, created, created_by, modified, modified_by, remark
		FROM tb_sys_role ORDER BY id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*entity.SysRole
	for rows.Next() {
		ro := &entity.SysRole{}
		rows.Scan(&ro.ID, &ro.TenantID, &ro.RoleName, &ro.RoleKey, &ro.Status, &ro.DataScope,
			&ro.MenuCheckStrictly, &ro.DeptCheckStrictly, &ro.Created, &ro.CreatedBy, &ro.Modified, &ro.ModifiedBy, &ro.Remark)
		roles = append(roles, ro)
	}
	return roles, nil
}

// RoleFindByID finds a role by ID
func (r *SystemRepository) RoleFindByID(ctx context.Context, id int64) (*entity.SysRole, error) {
	ro := &entity.SysRole{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, role_name, role_key, status, data_scope,
		       menu_check_strictly, dept_check_strictly, created, created_by, modified, modified_by, remark
		FROM tb_sys_role WHERE id=?
	`, id).Scan(&ro.ID, &ro.TenantID, &ro.RoleName, &ro.RoleKey, &ro.Status, &ro.DataScope,
		&ro.MenuCheckStrictly, &ro.DeptCheckStrictly, &ro.Created, &ro.CreatedBy, &ro.Modified, &ro.ModifiedBy, &ro.Remark)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return ro, err
}

// RoleSave creates a new role
func (r *SystemRepository) RoleSave(ctx context.Context, ro *entity.SysRole) error {
	id, _ := snowflake.GenerateID()
	now := time.Now()
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO tb_sys_role (id, tenant_id, role_name, role_key, status, data_scope,
		       menu_check_strictly, dept_check_strictly, created, created_by, modified, modified_by, remark)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, id, ro.TenantID, ro.RoleName, ro.RoleKey, ro.Status, ro.DataScope,
		ro.MenuCheckStrictly, ro.DeptCheckStrictly, now, ro.CreatedBy, now, ro.CreatedBy, ro.Remark)
	if err != nil {
		return err
	}
	ro.ID = id
	return nil
}

// RoleUpdate updates a role
func (r *SystemRepository) RoleUpdate(ctx context.Context, ro *entity.SysRole) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE tb_sys_role SET role_name=?, role_key=?, status=?, data_scope=?,
		       menu_check_strictly=?, dept_check_strictly=?, modified=?, modified_by=?, remark=?
		WHERE id=?
	`, ro.RoleName, ro.RoleKey, ro.Status, ro.DataScope,
		ro.MenuCheckStrictly, ro.DeptCheckStrictly, time.Now(), ro.ModifiedBy, ro.Remark, ro.ID)
	return err
}

// RoleDelete deletes a role
func (r *SystemRepository) RoleDelete(ctx context.Context, id int64) error {
	r.db.ExecContext(ctx, "DELETE FROM tb_sys_role_menu WHERE role_id=?", id)
	r.db.ExecContext(ctx, "DELETE FROM tb_sys_account_role WHERE role_id=?", id)
	_, err := r.db.ExecContext(ctx, "DELETE FROM tb_sys_role WHERE id=?", id)
	return err
}

// GetRoleMenuIds returns menu IDs for a role
func (r *SystemRepository) GetRoleMenuIds(ctx context.Context, roleID int64) ([]int64, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT menu_id FROM tb_sys_role_menu WHERE role_id=?", roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		rows.Scan(&id)
		ids = append(ids, id)
	}
	return ids, nil
}

// SaveRoleMenus saves role-menu associations
func (r *SystemRepository) SaveRoleMenus(ctx context.Context, roleID int64, menuIds []int64) error {
	r.db.ExecContext(ctx, "DELETE FROM tb_sys_role_menu WHERE role_id=?", roleID)
	for _, menuID := range menuIds {
		id, _ := snowflake.GenerateID()
		r.db.ExecContext(ctx, "INSERT INTO tb_sys_role_menu (id, role_id, menu_id) VALUES (?, ?, ?)",
			id, roleID, menuID)
	}
	return nil
}

// ========== Menu Operations ==========

// MenuList returns all menus
func (r *SystemRepository) MenuList(ctx context.Context) ([]*entity.SysMenu, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, parent_id, menu_type, menu_title, menu_url, component, menu_icon,
		       is_show, permission_tag, sort_no, status, created, created_by, modified, modified_by, remark
		FROM tb_sys_menu ORDER BY sort_no ASC, id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var menus []*entity.SysMenu
	for rows.Next() {
		m := &entity.SysMenu{}
		var menuUrl, component, menuIcon, permissionTag, remark sql.NullString
		rows.Scan(&m.ID, &m.ParentID, &m.MenuType, &m.MenuTitle, &menuUrl, &component, &menuIcon,
			&m.IsShow, &permissionTag, &m.SortNo, &m.Status, &m.Created, &m.CreatedBy, &m.Modified, &m.ModifiedBy, &remark)
		m.MenuUrl = menuUrl.String
		m.Component = component.String
		m.MenuIcon = menuIcon.String
		m.PermissionTag = permissionTag.String
		m.Remark = remark.String
		menus = append(menus, m)
	}
	return menus, nil
}

// MenuFindByID finds a menu by ID
func (r *SystemRepository) MenuFindByID(ctx context.Context, id int64) (*entity.SysMenu, error) {
	m := &entity.SysMenu{}
	var menuUrl, component, menuIcon, permissionTag, remark sql.NullString
	err := r.db.QueryRowContext(ctx, `
		SELECT id, parent_id, menu_type, menu_title, menu_url, component, menu_icon,
		       is_show, permission_tag, sort_no, status, created, created_by, modified, modified_by, remark
		FROM tb_sys_menu WHERE id=?
	`, id).Scan(&m.ID, &m.ParentID, &m.MenuType, &m.MenuTitle, &menuUrl, &component, &menuIcon,
		&m.IsShow, &permissionTag, &m.SortNo, &m.Status, &m.Created, &m.CreatedBy, &m.Modified, &m.ModifiedBy, &remark)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	m.MenuUrl = menuUrl.String
	m.Component = component.String
	m.MenuIcon = menuIcon.String
	m.PermissionTag = permissionTag.String
	m.Remark = remark.String
	return m, err
}

// MenuSave creates a new menu
func (r *SystemRepository) MenuSave(ctx context.Context, m *entity.SysMenu) error {
	id, _ := snowflake.GenerateID()
	now := time.Now()
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO tb_sys_menu (id, parent_id, menu_type, menu_title, menu_url, component, menu_icon,
		       is_show, permission_tag, sort_no, status, created, created_by, modified, modified_by, remark)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, id, m.ParentID, m.MenuType, m.MenuTitle, m.MenuUrl, m.Component, m.MenuIcon,
		m.IsShow, m.PermissionTag, m.SortNo, m.Status, now, m.CreatedBy, now, m.CreatedBy, m.Remark)
	if err != nil {
		return err
	}
	m.ID = id
	return nil
}

// MenuUpdate updates a menu
func (r *SystemRepository) MenuUpdate(ctx context.Context, m *entity.SysMenu) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE tb_sys_menu SET parent_id=?, menu_type=?, menu_title=?, menu_url=?, component=?, menu_icon=?,
		       is_show=?, permission_tag=?, sort_no=?, status=?, modified=?, modified_by=?, remark=?
		WHERE id=?
	`, m.ParentID, m.MenuType, m.MenuTitle, m.MenuUrl, m.Component, m.MenuIcon,
		m.IsShow, m.PermissionTag, m.SortNo, m.Status, time.Now(), m.ModifiedBy, m.Remark, m.ID)
	return err
}

// MenuDelete deletes a menu
func (r *SystemRepository) MenuDelete(ctx context.Context, id int64) error {
	r.db.ExecContext(ctx, "DELETE FROM tb_sys_role_menu WHERE menu_id=?", id)
	_, err := r.db.ExecContext(ctx, "DELETE FROM tb_sys_menu WHERE id=?", id)
	return err
}

// ========== Dept Operations ==========

// DeptList returns all departments
func (r *SystemRepository) DeptList(ctx context.Context) ([]*entity.SysDept, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, parent_id, ancestors, dept_name, dept_code,
		       sort_no, status, created, created_by, modified, modified_by, remark
		FROM tb_sys_dept ORDER BY sort_no ASC, id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var depts []*entity.SysDept
	for rows.Next() {
		d := &entity.SysDept{}
		var ancestors, deptCode, remark sql.NullString
		rows.Scan(&d.ID, &d.TenantID, &d.ParentID, &ancestors, &d.DeptName, &deptCode,
			&d.SortNo, &d.Status, &d.Created, &d.CreatedBy, &d.Modified, &d.ModifiedBy, &remark)
		d.Ancestors = ancestors.String
		d.DeptCode = deptCode.String
		d.Remark = remark.String
		depts = append(depts, d)
	}
	return depts, nil
}

// DeptFindByID finds a department by ID
func (r *SystemRepository) DeptFindByID(ctx context.Context, id int64) (*entity.SysDept, error) {
	d := &entity.SysDept{}
	var ancestors, deptCode, remark sql.NullString
	err := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, parent_id, ancestors, dept_name, dept_code,
		       sort_no, status, created, created_by, modified, modified_by, remark
		FROM tb_sys_dept WHERE id=?
	`, id).Scan(&d.ID, &d.TenantID, &d.ParentID, &ancestors, &d.DeptName, &deptCode,
		&d.SortNo, &d.Status, &d.Created, &d.CreatedBy, &d.Modified, &d.ModifiedBy, &remark)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	d.Ancestors = ancestors.String
	d.DeptCode = deptCode.String
	d.Remark = remark.String
	return d, err
}

// DeptSave creates a new department
func (r *SystemRepository) DeptSave(ctx context.Context, d *entity.SysDept) error {
	id, _ := snowflake.GenerateID()
	now := time.Now()
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO tb_sys_dept (id, tenant_id, parent_id, ancestors, dept_name, dept_code,
		       sort_no, status, created, created_by, modified, modified_by, remark)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, id, d.TenantID, d.ParentID, d.Ancestors, d.DeptName, d.DeptCode,
		d.SortNo, d.Status, now, d.CreatedBy, now, d.CreatedBy, d.Remark)
	if err != nil {
		return err
	}
	d.ID = id
	return nil
}

// DeptUpdate updates a department
func (r *SystemRepository) DeptUpdate(ctx context.Context, d *entity.SysDept) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE tb_sys_dept SET parent_id=?, ancestors=?, dept_name=?, dept_code=?,
		       sort_no=?, status=?, modified=?, modified_by=?, remark=?
		WHERE id=?
	`, d.ParentID, d.Ancestors, d.DeptName, d.DeptCode,
		d.SortNo, d.Status, time.Now(), d.ModifiedBy, d.Remark, d.ID)
	return err
}

// DeptDelete deletes a department
func (r *SystemRepository) DeptDelete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM tb_sys_dept WHERE id=?", id)
	return err
}

// DeptHasEmployees checks if department has employees
func (r *SystemRepository) DeptHasEmployees(ctx context.Context, deptID int64) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM tb_sys_account WHERE dept_id=?", deptID).Scan(&count)
	return count > 0, err
}

// ========== Dict Operations ==========

// DictPage returns paginated dictionaries
func (r *SystemRepository) DictPage(ctx context.Context, req *dto.DictQueryRequest) ([]*entity.SysDict, int64, error) {
	var conditions []string
	var args []interface{}

	if req.Name != "" {
		conditions = append(conditions, "name LIKE ?")
		args = append(args, "%"+req.Name+"%")
	}
	if req.Code != "" {
		conditions = append(conditions, "code LIKE ?")
		args = append(args, "%"+req.Code+"%")
	}
	if req.Status != nil {
		conditions = append(conditions, "status = ?")
		args = append(args, *req.Status)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	var total int64
	r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(*) FROM tb_sys_dict %s", whereClause), args...).Scan(&total)

	query := fmt.Sprintf(`
		SELECT id, name, code, description, dict_type, sort_no, status, options, created, modified
		FROM tb_sys_dict %s ORDER BY sort_no ASC, id DESC LIMIT ? OFFSET ?
	`, whereClause)
	args = append(args, req.GetPageSize(), req.GetOffset())

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var dicts []*entity.SysDict
	for rows.Next() {
		d := &entity.SysDict{}
		var name, description, options sql.NullString
		var created, modified sql.NullTime
		rows.Scan(&d.ID, &name, &d.Code, &description, &d.DictType, &d.SortNo, &d.Status, &options, &created, &modified)
		d.Name = name.String
		d.Description = description.String
		d.Options = options.String
		if created.Valid {
			d.Created = created.Time
		}
		if modified.Valid {
			d.Modified = modified.Time
		}
		dicts = append(dicts, d)
	}
	return dicts, total, nil
}

// DictFindByID finds a dictionary by ID
func (r *SystemRepository) DictFindByID(ctx context.Context, id int64) (*entity.SysDict, error) {
	d := &entity.SysDict{}
	var name, description, options sql.NullString
	var created, modified sql.NullTime
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, code, description, dict_type, sort_no, status, options, created, modified
		FROM tb_sys_dict WHERE id=?
	`, id).Scan(&d.ID, &name, &d.Code, &description, &d.DictType, &d.SortNo, &d.Status, &options, &created, &modified)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	d.Name = name.String
	d.Description = description.String
	d.Options = options.String
	if created.Valid {
		d.Created = created.Time
	}
	if modified.Valid {
		d.Modified = modified.Time
	}
	return d, err
}

// DictSave creates a new dictionary
func (r *SystemRepository) DictSave(ctx context.Context, d *entity.SysDict) error {
	id, _ := snowflake.GenerateID()
	now := time.Now()
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO tb_sys_dict (id, name, code, description, dict_type, sort_no, status, options, created, modified)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, id, d.Name, d.Code, d.Description, d.DictType, d.SortNo, d.Status, d.Options, now, now)
	if err != nil {
		return err
	}
	d.ID = id
	return nil
}

// DictUpdate updates a dictionary
func (r *SystemRepository) DictUpdate(ctx context.Context, d *entity.SysDict) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE tb_sys_dict SET name=?, code=?, description=?, dict_type=?, sort_no=?, status=?, options=?, modified=?
		WHERE id=?
	`, d.Name, d.Code, d.Description, d.DictType, d.SortNo, d.Status, d.Options, time.Now(), d.ID)
	return err
}

// DictDelete deletes a dictionary
func (r *SystemRepository) DictDelete(ctx context.Context, id int64) error {
	r.db.ExecContext(ctx, "DELETE FROM tb_sys_dict_item WHERE dict_id=?", id)
	_, err := r.db.ExecContext(ctx, "DELETE FROM tb_sys_dict WHERE id=?", id)
	return err
}

// DictItemPage returns paginated dictionary items
func (r *SystemRepository) DictItemPage(ctx context.Context, req *dto.DictItemQueryRequest) ([]*entity.SysDictItem, int64, error) {
	var conditions []string
	var args []interface{}

	if req.DictID > 0 {
		conditions = append(conditions, "dict_id = ?")
		args = append(args, req.DictID)
	}
	if req.Text != "" {
		conditions = append(conditions, "text LIKE ?")
		args = append(args, "%"+req.Text+"%")
	}
	if req.Status != nil {
		conditions = append(conditions, "status = ?")
		args = append(args, *req.Status)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	var total int64
	r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(*) FROM tb_sys_dict_item %s", whereClause), args...).Scan(&total)

	query := fmt.Sprintf(`
		SELECT id, dict_id, text, value, description, sort_no, css_content, css_class, remark, status, created, modified
		FROM tb_sys_dict_item %s ORDER BY sort_no ASC, id DESC LIMIT ? OFFSET ?
	`, whereClause)
	args = append(args, req.GetPageSize(), req.GetOffset())

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []*entity.SysDictItem
	for rows.Next() {
		i := &entity.SysDictItem{}
		var description, cssContent, cssClass, remark sql.NullString
		var created, modified sql.NullTime
		rows.Scan(&i.ID, &i.DictID, &i.Text, &i.Value, &description, &i.SortNo, &cssContent, &cssClass, &remark, &i.Status, &created, &modified)
		i.Description = description.String
		i.CssContent = cssContent.String
		i.CssClass = cssClass.String
		i.Remark = remark.String
		if created.Valid {
			i.Created = created.Time
		}
		if modified.Valid {
			i.Modified = modified.Time
		}
		items = append(items, i)
	}
	return items, total, nil
}

// DictItemSave creates a new dictionary item
func (r *SystemRepository) DictItemSave(ctx context.Context, i *entity.SysDictItem) error {
	id, _ := snowflake.GenerateID()
	now := time.Now()
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO tb_sys_dict_item (id, dict_id, text, value, description, sort_no, css_content, css_class, remark, status, created, modified)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, id, i.DictID, i.Text, i.Value, i.Description, i.SortNo, i.CssContent, i.CssClass, i.Remark, i.Status, now, now)
	if err != nil {
		return err
	}
	i.ID = id
	return nil
}

// DictItemUpdate updates a dictionary item
func (r *SystemRepository) DictItemUpdate(ctx context.Context, i *entity.SysDictItem) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE tb_sys_dict_item SET dict_id=?, text=?, value=?, description=?, sort_no=?, css_content=?, css_class=?, remark=?, status=?, modified=?
		WHERE id=?
	`, i.DictID, i.Text, i.Value, i.Description, i.SortNo, i.CssContent, i.CssClass, i.Remark, i.Status, time.Now(), i.ID)
	return err
}

// DictItemDelete deletes a dictionary item
func (r *SystemRepository) DictItemDelete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM tb_sys_dict_item WHERE id=?", id)
	return err
}

// DictItemFindByID finds a dictionary item by ID
func (r *SystemRepository) DictItemFindByID(ctx context.Context, id int64) (*entity.SysDictItem, error) {
	i := &entity.SysDictItem{}
	var description, cssContent, cssClass, remark sql.NullString
	var created, modified sql.NullTime
	err := r.db.QueryRowContext(ctx, `
		SELECT id, dict_id, text, value, description, sort_no, css_content, css_class, remark, status, created, modified
		FROM tb_sys_dict_item WHERE id=?
	`, id).Scan(&i.ID, &i.DictID, &i.Text, &i.Value, &description, &i.SortNo, &cssContent, &cssClass, &remark, &i.Status, &created, &modified)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	i.Description = description.String
	i.CssContent = cssContent.String
	i.CssClass = cssClass.String
	i.Remark = remark.String
	if created.Valid {
		i.Created = created.Time
	}
	if modified.Valid {
		i.Modified = modified.Time
	}
	return i, err
}

// Helper function to convert camelCase to snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		if r >= 'A' && r <= 'Z' {
			result.WriteRune(r + 32) // toLower
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}
