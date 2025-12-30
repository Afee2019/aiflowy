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
)

// SystemService handles system management business logic
type SystemService struct {
	repo *repository.SystemRepository
}

// NewSystemService creates a new SystemService
func NewSystemService() *SystemService {
	return &SystemService{
		repo: repository.NewSystemRepository(),
	}
}

// ========== Account Operations ==========

// AccountPage returns paginated accounts
func (s *SystemService) AccountPage(ctx context.Context, req *dto.AccountQueryRequest) (*dto.PageResponse, error) {
	accounts, total, err := s.repo.AccountPage(ctx, req)
	if err != nil {
		return nil, apierrors.InternalError("查询用户列表失败")
	}

	// Note: Role IDs are loaded on demand in AccountDetail
	// to avoid N+1 query issue on list pages

	return dto.NewPageResponse(req.GetPageNumber(), req.GetPageSize(), total, accounts), nil
}

// AccountSave creates or updates an account
func (s *SystemService) AccountSave(ctx context.Context, req *dto.AccountSaveRequest, operatorID int64) (*entity.SysAccount, error) {
	// Check if login name exists
	exists, err := s.repo.AccountExistsByLoginName(ctx, req.LoginName, req.ID)
	if err != nil {
		return nil, apierrors.InternalError("检查用户名失败")
	}
	if exists {
		return nil, apierrors.New(1, "用户名已存在")
	}

	if req.ID == 0 {
		// Create new account
		if req.Password == "" {
			return nil, apierrors.BadRequest("密码不能为空")
		}

		// Hash password with BCrypt
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, apierrors.InternalError("密码加密失败")
		}

		account := &entity.SysAccount{
			DeptID:      req.DeptID,
			TenantID:    1000000, // Default tenant ID
			LoginName:   req.LoginName,
			Password:    string(hashedPassword),
			AccountType: req.AccountType,
			Nickname:    req.Nickname,
			Mobile:      req.Mobile,
			Email:       req.Email,
			Avatar:      req.Avatar,
			Status:      req.Status,
			Remark:      req.Remark,
			CreatedBy:   operatorID,
		}

		if err := s.repo.AccountSave(ctx, account); err != nil {
			return nil, apierrors.InternalError("创建用户失败")
		}

		// Save role associations
		if len(req.RoleIds) > 0 {
			s.repo.SaveAccountRoles(ctx, account.ID, req.RoleIds)
		}

		return account, nil
	}

	// Update existing account
	accountRepo := repository.NewAccountRepository()
	account, err := accountRepo.FindByID(ctx, req.ID)
	if err != nil || account == nil {
		return nil, apierrors.NotFound("用户不存在")
	}

	account.DeptID = req.DeptID
	account.Nickname = req.Nickname
	account.Mobile = req.Mobile
	account.Email = req.Email
	account.Avatar = req.Avatar
	account.Status = req.Status
	account.Remark = req.Remark
	account.ModifiedBy = operatorID

	if err := s.repo.AccountUpdate(ctx, account); err != nil {
		return nil, apierrors.InternalError("更新用户失败")
	}

	// Update role associations
	s.repo.SaveAccountRoles(ctx, account.ID, req.RoleIds)

	return account, nil
}

// AccountDelete deletes an account
func (s *SystemService) AccountDelete(ctx context.Context, id int64) error {
	// Check if it's a super admin (ID = 1)
	if id == 1 {
		return apierrors.New(1, "不能删除超级管理员")
	}

	accountRepo := repository.NewAccountRepository()
	account, err := accountRepo.FindByID(ctx, id)
	if err != nil || account == nil {
		return apierrors.NotFound("用户不存在")
	}

	// Check account type
	if account.AccountType == 99 { // Super admin
		return apierrors.New(1, "不能删除超级管理员")
	}

	return s.repo.AccountDelete(ctx, id)
}

// AccountDetail returns account detail
func (s *SystemService) AccountDetail(ctx context.Context, id int64) (*entity.SysAccount, error) {
	accountRepo := repository.NewAccountRepository()
	account, err := accountRepo.FindByID(ctx, id)
	if err != nil {
		return nil, apierrors.InternalError("查询用户失败")
	}
	if account == nil {
		return nil, apierrors.NotFound("用户不存在")
	}
	return account, nil
}

// UpdatePassword updates user password
func (s *SystemService) UpdatePassword(ctx context.Context, userID int64, req *dto.UpdatePasswordRequest) error {
	if req.NewPassword != req.ConfirmPassword {
		return apierrors.New(2, "两次密码不一致")
	}

	accountRepo := repository.NewAccountRepository()
	account, err := accountRepo.FindByID(ctx, userID)
	if err != nil || account == nil {
		return apierrors.NotFound("用户不存在")
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(req.Password)); err != nil {
		return apierrors.New(1, "密码不正确")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return apierrors.InternalError("密码加密失败")
	}

	// Update password directly via SQL
	db := repository.GetDB()
	_, err = db.ExecContext(ctx, "UPDATE tb_sys_account SET password=? WHERE id=?", string(hashedPassword), userID)
	if err != nil {
		return apierrors.InternalError("更新密码失败")
	}

	return nil
}

// ========== Role Operations ==========

// RolePage returns paginated roles
func (s *SystemService) RolePage(ctx context.Context, req *dto.RoleQueryRequest) (*dto.PageResponse, error) {
	roles, total, err := s.repo.RolePage(ctx, req)
	if err != nil {
		return nil, apierrors.InternalError("查询角色列表失败")
	}
	return dto.NewPageResponse(req.GetPageNumber(), req.GetPageSize(), total, roles), nil
}

// RoleList returns all roles
func (s *SystemService) RoleList(ctx context.Context) ([]*entity.SysRole, error) {
	return s.repo.RoleList(ctx)
}

// RoleSave creates or updates a role
func (s *SystemService) RoleSave(ctx context.Context, req *dto.RoleSaveRequest, operatorID int64) (*entity.SysRole, error) {
	if req.ID == 0 {
		role := &entity.SysRole{
			TenantID:          1000000,
			RoleName:          req.RoleName,
			RoleKey:           req.RoleKey,
			Status:            req.Status,
			DataScope:         req.DataScope,
			MenuCheckStrictly: req.MenuCheckStrictly,
			DeptCheckStrictly: req.DeptCheckStrictly,
			Remark:            req.Remark,
			CreatedBy:         operatorID,
		}

		if err := s.repo.RoleSave(ctx, role); err != nil {
			return nil, apierrors.InternalError("创建角色失败")
		}

		// Save menu associations
		if len(req.MenuIds) > 0 {
			s.repo.SaveRoleMenus(ctx, role.ID, req.MenuIds)
		}

		return role, nil
	}

	// Update existing role
	role, err := s.repo.RoleFindByID(ctx, req.ID)
	if err != nil || role == nil {
		return nil, apierrors.NotFound("角色不存在")
	}

	role.RoleName = req.RoleName
	role.RoleKey = req.RoleKey
	role.Status = req.Status
	role.DataScope = req.DataScope
	role.MenuCheckStrictly = req.MenuCheckStrictly
	role.DeptCheckStrictly = req.DeptCheckStrictly
	role.Remark = req.Remark
	role.ModifiedBy = operatorID

	if err := s.repo.RoleUpdate(ctx, role); err != nil {
		return nil, apierrors.InternalError("更新角色失败")
	}

	// Update menu associations
	s.repo.SaveRoleMenus(ctx, role.ID, req.MenuIds)

	return role, nil
}

// RoleDelete deletes a role
func (s *SystemService) RoleDelete(ctx context.Context, id int64) error {
	role, err := s.repo.RoleFindByID(ctx, id)
	if err != nil || role == nil {
		return apierrors.NotFound("角色不存在")
	}

	if role.RoleKey == entity.RoleKeySuperAdmin {
		return apierrors.New(1, "超级管理员角色不能删除")
	}
	if role.RoleKey == entity.RoleKeyTenantAdmin {
		return apierrors.New(1, "租户管理员角色不能删除")
	}

	return s.repo.RoleDelete(ctx, id)
}

// RoleDetail returns role detail
func (s *SystemService) RoleDetail(ctx context.Context, id int64) (*entity.SysRole, error) {
	role, err := s.repo.RoleFindByID(ctx, id)
	if err != nil {
		return nil, apierrors.InternalError("查询角色失败")
	}
	if role == nil {
		return nil, apierrors.NotFound("角色不存在")
	}

	// Load menu IDs
	menuIds, _ := s.repo.GetRoleMenuIds(ctx, id)
	role.MenuIds = menuIds

	return role, nil
}

// GetRoleMenuIds returns menu IDs for a role
func (s *SystemService) GetRoleMenuIds(ctx context.Context, roleID int64) ([]int64, error) {
	return s.repo.GetRoleMenuIds(ctx, roleID)
}

// ========== Menu Operations ==========

// MenuList returns all menus
func (s *SystemService) MenuList(ctx context.Context, asTree bool) ([]*entity.SysMenu, error) {
	menus, err := s.repo.MenuList(ctx)
	if err != nil {
		return nil, apierrors.InternalError("查询菜单列表失败")
	}

	if asTree {
		return buildMenuTree(menus, 0), nil
	}
	return menus, nil
}

// MenuSave creates or updates a menu
func (s *SystemService) MenuSave(ctx context.Context, req *dto.MenuSaveRequest, operatorID int64) (*entity.SysMenu, error) {
	if req.ID == 0 {
		menu := &entity.SysMenu{
			ParentID:      req.ParentID,
			MenuType:      req.MenuType,
			MenuTitle:     req.MenuTitle,
			MenuUrl:       req.MenuUrl,
			Component:     req.Component,
			MenuIcon:      req.MenuIcon,
			IsShow:        req.IsShow,
			PermissionTag: req.PermissionTag,
			SortNo:        req.SortNo,
			Status:        req.Status,
			Remark:        req.Remark,
			CreatedBy:     operatorID,
		}

		if err := s.repo.MenuSave(ctx, menu); err != nil {
			return nil, apierrors.InternalError("创建菜单失败")
		}

		// Auto-add to super admin role
		superAdminRoleID := int64(1) // Assuming super admin role ID is 1
		menuIds, _ := s.repo.GetRoleMenuIds(ctx, superAdminRoleID)
		menuIds = append(menuIds, menu.ID)
		s.repo.SaveRoleMenus(ctx, superAdminRoleID, menuIds)

		return menu, nil
	}

	// Update existing menu
	menu, err := s.repo.MenuFindByID(ctx, req.ID)
	if err != nil || menu == nil {
		return nil, apierrors.NotFound("菜单不存在")
	}

	menu.ParentID = req.ParentID
	menu.MenuType = req.MenuType
	menu.MenuTitle = req.MenuTitle
	menu.MenuUrl = req.MenuUrl
	menu.Component = req.Component
	menu.MenuIcon = req.MenuIcon
	menu.IsShow = req.IsShow
	menu.PermissionTag = req.PermissionTag
	menu.SortNo = req.SortNo
	menu.Status = req.Status
	menu.Remark = req.Remark
	menu.ModifiedBy = operatorID

	if err := s.repo.MenuUpdate(ctx, menu); err != nil {
		return nil, apierrors.InternalError("更新菜单失败")
	}

	return menu, nil
}

// MenuDelete deletes a menu
func (s *SystemService) MenuDelete(ctx context.Context, id int64) error {
	return s.repo.MenuDelete(ctx, id)
}

// MenuDetail returns menu detail
func (s *SystemService) MenuDetail(ctx context.Context, id int64) (*entity.SysMenu, error) {
	menu, err := s.repo.MenuFindByID(ctx, id)
	if err != nil {
		return nil, apierrors.InternalError("查询菜单失败")
	}
	if menu == nil {
		return nil, apierrors.NotFound("菜单不存在")
	}
	return menu, nil
}

// GetMenuCheckedByRoleId returns menu IDs for a role
func (s *SystemService) GetMenuCheckedByRoleId(ctx context.Context, roleID int64) ([]string, error) {
	ids, err := s.repo.GetRoleMenuIds(ctx, roleID)
	if err != nil {
		return nil, err
	}
	result := make([]string, len(ids))
	for i, id := range ids {
		result[i] = strconv.FormatInt(id, 10)
	}
	return result, nil
}

// ========== Dept Operations ==========

// DeptList returns all departments
func (s *SystemService) DeptList(ctx context.Context, asTree bool) ([]*entity.SysDept, error) {
	depts, err := s.repo.DeptList(ctx)
	if err != nil {
		return nil, apierrors.InternalError("查询部门列表失败")
	}

	if asTree {
		return buildDeptTree(depts, 0), nil
	}
	return depts, nil
}

// DeptSave creates or updates a department
func (s *SystemService) DeptSave(ctx context.Context, req *dto.DeptSaveRequest, operatorID int64) (*entity.SysDept, error) {
	// Calculate ancestors
	ancestors := "0"
	if req.ParentID > 0 {
		parent, err := s.repo.DeptFindByID(ctx, req.ParentID)
		if err != nil || parent == nil {
			return nil, apierrors.NotFound("父部门不存在")
		}
		if parent.Ancestors != "" {
			ancestors = parent.Ancestors + "," + strconv.FormatInt(req.ParentID, 10)
		} else {
			ancestors = strconv.FormatInt(req.ParentID, 10)
		}
	}

	if req.ID == 0 {
		dept := &entity.SysDept{
			TenantID:  1000000,
			ParentID:  req.ParentID,
			Ancestors: ancestors,
			DeptName:  req.DeptName,
			DeptCode:  req.DeptCode,
			SortNo:    req.SortNo,
			Status:    req.Status,
			Remark:    req.Remark,
			CreatedBy: operatorID,
		}

		if err := s.repo.DeptSave(ctx, dept); err != nil {
			return nil, apierrors.InternalError("创建部门失败")
		}

		return dept, nil
	}

	// Update existing department
	dept, err := s.repo.DeptFindByID(ctx, req.ID)
	if err != nil || dept == nil {
		return nil, apierrors.NotFound("部门不存在")
	}

	dept.ParentID = req.ParentID
	dept.Ancestors = ancestors
	dept.DeptName = req.DeptName
	dept.DeptCode = req.DeptCode
	dept.SortNo = req.SortNo
	dept.Status = req.Status
	dept.Remark = req.Remark
	dept.ModifiedBy = operatorID

	if err := s.repo.DeptUpdate(ctx, dept); err != nil {
		return nil, apierrors.InternalError("更新部门失败")
	}

	return dept, nil
}

// DeptDelete deletes a department
func (s *SystemService) DeptDelete(ctx context.Context, id int64) error {
	dept, err := s.repo.DeptFindByID(ctx, id)
	if err != nil || dept == nil {
		return apierrors.NotFound("部门不存在")
	}

	if dept.IsRoot() {
		return apierrors.New(1, "无法删除根部门")
	}

	hasEmployees, _ := s.repo.DeptHasEmployees(ctx, id)
	if hasEmployees {
		return apierrors.New(1, "该部门下有员工，不能删除")
	}

	return s.repo.DeptDelete(ctx, id)
}

// DeptDetail returns department detail
func (s *SystemService) DeptDetail(ctx context.Context, id int64) (*entity.SysDept, error) {
	dept, err := s.repo.DeptFindByID(ctx, id)
	if err != nil {
		return nil, apierrors.InternalError("查询部门失败")
	}
	if dept == nil {
		return nil, apierrors.NotFound("部门不存在")
	}
	return dept, nil
}

// ========== Dict Operations ==========

// DictPage returns paginated dictionaries
func (s *SystemService) DictPage(ctx context.Context, req *dto.DictQueryRequest) (*dto.PageResponse, error) {
	dicts, total, err := s.repo.DictPage(ctx, req)
	if err != nil {
		return nil, apierrors.InternalError("查询字典列表失败")
	}
	return dto.NewPageResponse(req.GetPageNumber(), req.GetPageSize(), total, dicts), nil
}

// DictSave creates or updates a dictionary
func (s *SystemService) DictSave(ctx context.Context, req *dto.DictSaveRequest) (*entity.SysDict, error) {
	if req.ID == 0 {
		dict := &entity.SysDict{
			Name:        req.Name,
			Code:        req.Code,
			Description: req.Description,
			DictType:    req.DictType,
			SortNo:      req.SortNo,
			Status:      req.Status,
			Options:     req.Options,
		}

		if err := s.repo.DictSave(ctx, dict); err != nil {
			return nil, apierrors.InternalError("创建字典失败")
		}

		return dict, nil
	}

	// Update existing dictionary
	dict, err := s.repo.DictFindByID(ctx, req.ID)
	if err != nil || dict == nil {
		return nil, apierrors.NotFound("字典不存在")
	}

	dict.Name = req.Name
	dict.Code = req.Code
	dict.Description = req.Description
	dict.DictType = req.DictType
	dict.SortNo = req.SortNo
	dict.Status = req.Status
	dict.Options = req.Options

	if err := s.repo.DictUpdate(ctx, dict); err != nil {
		return nil, apierrors.InternalError("更新字典失败")
	}

	return dict, nil
}

// DictDelete deletes a dictionary
func (s *SystemService) DictDelete(ctx context.Context, id int64) error {
	return s.repo.DictDelete(ctx, id)
}

// DictDetail returns dictionary detail
func (s *SystemService) DictDetail(ctx context.Context, id int64) (*entity.SysDict, error) {
	dict, err := s.repo.DictFindByID(ctx, id)
	if err != nil {
		return nil, apierrors.InternalError("查询字典失败")
	}
	if dict == nil {
		return nil, apierrors.NotFound("字典不存在")
	}
	return dict, nil
}

// DictItemPage returns paginated dictionary items
func (s *SystemService) DictItemPage(ctx context.Context, req *dto.DictItemQueryRequest) (*dto.PageResponse, error) {
	items, total, err := s.repo.DictItemPage(ctx, req)
	if err != nil {
		return nil, apierrors.InternalError("查询字典项列表失败")
	}
	return dto.NewPageResponse(req.GetPageNumber(), req.GetPageSize(), total, items), nil
}

// DictItemSave creates or updates a dictionary item
func (s *SystemService) DictItemSave(ctx context.Context, req *dto.DictItemSaveRequest) (*entity.SysDictItem, error) {
	if req.ID == 0 {
		item := &entity.SysDictItem{
			DictID:      req.DictID,
			Text:        req.Text,
			Value:       req.Value,
			Description: req.Description,
			SortNo:      req.SortNo,
			CssContent:  req.CssContent,
			CssClass:    req.CssClass,
			Remark:      req.Remark,
			Status:      req.Status,
		}

		if err := s.repo.DictItemSave(ctx, item); err != nil {
			return nil, apierrors.InternalError("创建字典项失败")
		}

		return item, nil
	}

	// Update existing item
	item, err := s.repo.DictItemFindByID(ctx, req.ID)
	if err != nil || item == nil {
		return nil, apierrors.NotFound("字典项不存在")
	}

	item.DictID = req.DictID
	item.Text = req.Text
	item.Value = req.Value
	item.Description = req.Description
	item.SortNo = req.SortNo
	item.CssContent = req.CssContent
	item.CssClass = req.CssClass
	item.Remark = req.Remark
	item.Status = req.Status

	if err := s.repo.DictItemUpdate(ctx, item); err != nil {
		return nil, apierrors.InternalError("更新字典项失败")
	}

	return item, nil
}

// DictItemDelete deletes a dictionary item
func (s *SystemService) DictItemDelete(ctx context.Context, id int64) error {
	return s.repo.DictItemDelete(ctx, id)
}

// DictItemDetail returns dictionary item detail
func (s *SystemService) DictItemDetail(ctx context.Context, id int64) (*entity.SysDictItem, error) {
	item, err := s.repo.DictItemFindByID(ctx, id)
	if err != nil {
		return nil, apierrors.InternalError("查询字典项失败")
	}
	if item == nil {
		return nil, apierrors.NotFound("字典项不存在")
	}
	return item, nil
}

// Helper functions

func buildMenuTree(menus []*entity.SysMenu, parentID int64) []*entity.SysMenu {
	var result []*entity.SysMenu
	for _, m := range menus {
		if m.ParentID == parentID {
			m.Children = buildMenuTree(menus, m.ID)
			result = append(result, m)
		}
	}
	return result
}

func buildDeptTree(depts []*entity.SysDept, parentID int64) []*entity.SysDept {
	var result []*entity.SysDept
	for _, d := range depts {
		if d.ParentID == parentID {
			d.Children = buildDeptTree(depts, d.ID)
			result = append(result, d)
		}
	}
	return result
}

// GetAccountRoleIds returns role IDs for an account
func (s *SystemService) GetAccountRoleIds(ctx context.Context, accountID int64) ([]int64, error) {
	return s.repo.GetAccountRoleIds(ctx, accountID)
}

// Placeholder for validation
func validateAccountRequest(req *dto.AccountSaveRequest) error {
	if req.LoginName == "" {
		return fmt.Errorf("登录名不能为空")
	}
	return nil
}
