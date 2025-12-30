package system

import (
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/aiflowy/aiflowy-go/internal/dto"
	apierrors "github.com/aiflowy/aiflowy-go/internal/errors"
	"github.com/aiflowy/aiflowy-go/internal/handler/auth"
	"github.com/aiflowy/aiflowy-go/internal/service"
	"github.com/aiflowy/aiflowy-go/pkg/response"
)

// Handler handles system management endpoints
type Handler struct {
	systemService *service.SystemService
}

// NewHandler creates a new system handler
func NewHandler() *Handler {
	return &Handler{
		systemService: service.NewSystemService(),
	}
}

// ========== Account Handlers ==========

// AccountPage returns paginated accounts
// GET /api/v1/sysAccount/page
func (h *Handler) AccountPage(c echo.Context) error {
	var req dto.AccountQueryRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	result, err := h.systemService.AccountPage(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return response.Success(c, result)
}

// AccountList returns all accounts
// GET /api/v1/sysAccount/list
func (h *Handler) AccountList(c echo.Context) error {
	req := &dto.AccountQueryRequest{
		PageRequest: dto.PageRequest{PageSize: 1000},
	}
	result, err := h.systemService.AccountPage(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return response.Success(c, result.Rows)
}

// AccountSave creates a new account
// POST /api/v1/sysAccount/save
func (h *Handler) AccountSave(c echo.Context) error {
	var req dto.AccountSaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	operatorID := auth.GetCurrentUserID(c)
	account, err := h.systemService.AccountSave(c.Request().Context(), &req, operatorID)
	if err != nil {
		return err
	}

	return response.Success(c, map[string]interface{}{"id": account.ID})
}

// AccountUpdate updates an account
// POST /api/v1/sysAccount/update
func (h *Handler) AccountUpdate(c echo.Context) error {
	var req dto.AccountSaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	if req.ID == 0 {
		return apierrors.BadRequest("用户ID不能为空")
	}

	operatorID := auth.GetCurrentUserID(c)
	_, err := h.systemService.AccountSave(c.Request().Context(), &req, operatorID)
	if err != nil {
		return err
	}

	return response.Success(c, nil)
}

// AccountRemove deletes an account
// POST /api/v1/sysAccount/remove
func (h *Handler) AccountRemove(c echo.Context) error {
	var req dto.IDRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	if err := h.systemService.AccountDelete(c.Request().Context(), req.ID); err != nil {
		return err
	}

	return response.Success(c, true)
}

// AccountDetail returns account detail
// GET /api/v1/sysAccount/detail
func (h *Handler) AccountDetail(c echo.Context) error {
	idStr := c.QueryParam("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return apierrors.BadRequest("无效的用户ID")
	}

	account, err := h.systemService.AccountDetail(c.Request().Context(), id)
	if err != nil {
		return err
	}

	// Load role IDs
	roleIds, _ := h.systemService.GetAccountRoleIds(c.Request().Context(), id)

	return response.Success(c, map[string]interface{}{
		"id":          account.ID,
		"deptId":      account.DeptID,
		"loginName":   account.LoginName,
		"accountType": account.AccountType,
		"nickname":    account.Nickname,
		"mobile":      account.Mobile,
		"email":       account.Email,
		"avatar":      account.Avatar,
		"status":      account.Status,
		"remark":      account.Remark,
		"created":     account.Created,
		"roleIds":     roleIds,
	})
}

// UpdatePassword updates current user's password
// POST /api/v1/sysAccount/updatePassword
func (h *Handler) UpdatePassword(c echo.Context) error {
	var req dto.UpdatePasswordRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	userID := auth.GetCurrentUserID(c)
	if err := h.systemService.UpdatePassword(c.Request().Context(), userID, &req); err != nil {
		return err
	}

	return response.Success(c, nil)
}

// MyProfile returns current user's profile
// GET /api/v1/sysAccount/myProfile
func (h *Handler) MyProfile(c echo.Context) error {
	userID := auth.GetCurrentUserID(c)
	account, err := h.systemService.AccountDetail(c.Request().Context(), userID)
	if err != nil {
		return err
	}

	return response.Success(c, map[string]interface{}{
		"id":        account.ID,
		"loginName": account.LoginName,
		"nickname":  account.Nickname,
		"mobile":    account.Mobile,
		"email":     account.Email,
		"avatar":    account.Avatar,
	})
}

// ========== Role Handlers ==========

// RolePage returns paginated roles
// GET /api/v1/sysRole/page
func (h *Handler) RolePage(c echo.Context) error {
	var req dto.RoleQueryRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	result, err := h.systemService.RolePage(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return response.Success(c, result)
}

// RoleList returns all roles
// GET /api/v1/sysRole/list
func (h *Handler) RoleList(c echo.Context) error {
	roles, err := h.systemService.RoleList(c.Request().Context())
	if err != nil {
		return err
	}
	return response.Success(c, roles)
}

// RoleSave creates a new role
// POST /api/v1/sysRole/save
func (h *Handler) RoleSave(c echo.Context) error {
	var req dto.RoleSaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	operatorID := auth.GetCurrentUserID(c)
	role, err := h.systemService.RoleSave(c.Request().Context(), &req, operatorID)
	if err != nil {
		return err
	}

	return response.Success(c, map[string]interface{}{"id": role.ID})
}

// RoleUpdate updates a role
// POST /api/v1/sysRole/update
func (h *Handler) RoleUpdate(c echo.Context) error {
	var req dto.RoleSaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	if req.ID == 0 {
		return apierrors.BadRequest("角色ID不能为空")
	}

	operatorID := auth.GetCurrentUserID(c)
	_, err := h.systemService.RoleSave(c.Request().Context(), &req, operatorID)
	if err != nil {
		return err
	}

	return response.Success(c, nil)
}

// RoleRemove deletes a role
// POST /api/v1/sysRole/remove
func (h *Handler) RoleRemove(c echo.Context) error {
	var req dto.IDRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	if err := h.systemService.RoleDelete(c.Request().Context(), req.ID); err != nil {
		return err
	}

	return response.Success(c, true)
}

// RoleDetail returns role detail
// GET /api/v1/sysRole/detail
func (h *Handler) RoleDetail(c echo.Context) error {
	idStr := c.QueryParam("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return apierrors.BadRequest("无效的角色ID")
	}

	role, err := h.systemService.RoleDetail(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return response.Success(c, role)
}

// GetRoleMenuIds returns menu IDs for a role
// GET /api/v1/sysRole/getRoleMenuIds
func (h *Handler) GetRoleMenuIds(c echo.Context) error {
	roleIDStr := c.QueryParam("roleId")
	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		return apierrors.BadRequest("无效的角色ID")
	}

	menuIds, err := h.systemService.GetRoleMenuIds(c.Request().Context(), roleID)
	if err != nil {
		return err
	}

	return response.Success(c, menuIds)
}

// ========== Menu Handlers ==========

// MenuList returns all menus
// GET /api/v1/sysMenu/list
func (h *Handler) MenuList(c echo.Context) error {
	asTree := c.QueryParam("asTree") == "true"

	menus, err := h.systemService.MenuList(c.Request().Context(), asTree)
	if err != nil {
		return err
	}

	return response.Success(c, menus)
}

// MenuSave creates a new menu
// POST /api/v1/sysMenu/save
func (h *Handler) MenuSave(c echo.Context) error {
	var req dto.MenuSaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	operatorID := auth.GetCurrentUserID(c)
	menu, err := h.systemService.MenuSave(c.Request().Context(), &req, operatorID)
	if err != nil {
		return err
	}

	return response.Success(c, map[string]interface{}{"id": menu.ID})
}

// MenuUpdate updates a menu
// POST /api/v1/sysMenu/update
func (h *Handler) MenuUpdate(c echo.Context) error {
	var req dto.MenuSaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	if req.ID == 0 {
		return apierrors.BadRequest("菜单ID不能为空")
	}

	operatorID := auth.GetCurrentUserID(c)
	_, err := h.systemService.MenuSave(c.Request().Context(), &req, operatorID)
	if err != nil {
		return err
	}

	return response.Success(c, nil)
}

// MenuRemove deletes a menu
// POST /api/v1/sysMenu/remove
func (h *Handler) MenuRemove(c echo.Context) error {
	var req dto.IDRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	if err := h.systemService.MenuDelete(c.Request().Context(), req.ID); err != nil {
		return err
	}

	return response.Success(c, true)
}

// MenuDetail returns menu detail
// GET /api/v1/sysMenu/detail
func (h *Handler) MenuDetail(c echo.Context) error {
	idStr := c.QueryParam("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return apierrors.BadRequest("无效的菜单ID")
	}

	menu, err := h.systemService.MenuDetail(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return response.Success(c, menu)
}

// GetMenuCheckedByRoleId returns checked menu IDs for a role
// GET /api/v1/sysMenu/getCheckedByRoleId/:roleId
func (h *Handler) GetMenuCheckedByRoleId(c echo.Context) error {
	roleIDStr := c.Param("roleId")
	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		return apierrors.BadRequest("无效的角色ID")
	}

	menuIds, err := h.systemService.GetMenuCheckedByRoleId(c.Request().Context(), roleID)
	if err != nil {
		return err
	}

	return response.Success(c, menuIds)
}

// ========== Dept Handlers ==========

// DeptList returns all departments
// GET /api/v1/sysDept/list
func (h *Handler) DeptList(c echo.Context) error {
	asTree := c.QueryParam("asTree") != "false" // Default to true

	depts, err := h.systemService.DeptList(c.Request().Context(), asTree)
	if err != nil {
		return err
	}

	return response.Success(c, depts)
}

// DeptSave creates a new department
// POST /api/v1/sysDept/save
func (h *Handler) DeptSave(c echo.Context) error {
	var req dto.DeptSaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	operatorID := auth.GetCurrentUserID(c)
	dept, err := h.systemService.DeptSave(c.Request().Context(), &req, operatorID)
	if err != nil {
		return err
	}

	return response.Success(c, map[string]interface{}{"id": dept.ID})
}

// DeptUpdate updates a department
// POST /api/v1/sysDept/update
func (h *Handler) DeptUpdate(c echo.Context) error {
	var req dto.DeptSaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	if req.ID == 0 {
		return apierrors.BadRequest("部门ID不能为空")
	}

	operatorID := auth.GetCurrentUserID(c)
	_, err := h.systemService.DeptSave(c.Request().Context(), &req, operatorID)
	if err != nil {
		return err
	}

	return response.Success(c, nil)
}

// DeptRemove deletes a department
// POST /api/v1/sysDept/remove
func (h *Handler) DeptRemove(c echo.Context) error {
	var req dto.IDRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	if err := h.systemService.DeptDelete(c.Request().Context(), req.ID); err != nil {
		return err
	}

	return response.Success(c, true)
}

// DeptDetail returns department detail
// GET /api/v1/sysDept/detail
func (h *Handler) DeptDetail(c echo.Context) error {
	idStr := c.QueryParam("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return apierrors.BadRequest("无效的部门ID")
	}

	dept, err := h.systemService.DeptDetail(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return response.Success(c, dept)
}

// ========== Dict Handlers ==========

// DictPage returns paginated dictionaries
// GET /api/v1/sysDict/page
func (h *Handler) DictPage(c echo.Context) error {
	var req dto.DictQueryRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	result, err := h.systemService.DictPage(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return response.Success(c, result)
}

// DictList returns all dictionaries
// GET /api/v1/sysDict/list
func (h *Handler) DictList(c echo.Context) error {
	req := &dto.DictQueryRequest{
		PageRequest: dto.PageRequest{PageSize: 1000},
	}
	result, err := h.systemService.DictPage(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return response.Success(c, result.Rows)
}

// DictSave creates a new dictionary
// POST /api/v1/sysDict/save
func (h *Handler) DictSave(c echo.Context) error {
	var req dto.DictSaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	dict, err := h.systemService.DictSave(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return response.Success(c, map[string]interface{}{"id": dict.ID})
}

// DictUpdate updates a dictionary
// POST /api/v1/sysDict/update
func (h *Handler) DictUpdate(c echo.Context) error {
	var req dto.DictSaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	if req.ID == 0 {
		return apierrors.BadRequest("字典ID不能为空")
	}

	_, err := h.systemService.DictSave(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return response.Success(c, nil)
}

// DictRemove deletes a dictionary
// POST /api/v1/sysDict/remove
func (h *Handler) DictRemove(c echo.Context) error {
	var req dto.IDRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	if err := h.systemService.DictDelete(c.Request().Context(), req.ID); err != nil {
		return err
	}

	return response.Success(c, true)
}

// DictDetail returns dictionary detail
// GET /api/v1/sysDict/detail
func (h *Handler) DictDetail(c echo.Context) error {
	idStr := c.QueryParam("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return apierrors.BadRequest("无效的字典ID")
	}

	dict, err := h.systemService.DictDetail(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return response.Success(c, dict)
}

// ========== DictItem Handlers ==========

// DictItemPage returns paginated dictionary items
// GET /api/v1/sysDictItem/page
func (h *Handler) DictItemPage(c echo.Context) error {
	var req dto.DictItemQueryRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	result, err := h.systemService.DictItemPage(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return response.Success(c, result)
}

// DictItemList returns dictionary items by dict ID
// GET /api/v1/sysDictItem/list
func (h *Handler) DictItemList(c echo.Context) error {
	dictIDStr := c.QueryParam("dictId")
	dictID, _ := strconv.ParseInt(dictIDStr, 10, 64)

	req := &dto.DictItemQueryRequest{
		PageRequest: dto.PageRequest{PageSize: 1000},
		DictID:      dictID,
	}
	result, err := h.systemService.DictItemPage(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return response.Success(c, result.Rows)
}

// DictItemSave creates a new dictionary item
// POST /api/v1/sysDictItem/save
func (h *Handler) DictItemSave(c echo.Context) error {
	var req dto.DictItemSaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	item, err := h.systemService.DictItemSave(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return response.Success(c, map[string]interface{}{"id": item.ID})
}

// DictItemUpdate updates a dictionary item
// POST /api/v1/sysDictItem/update
func (h *Handler) DictItemUpdate(c echo.Context) error {
	var req dto.DictItemSaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	if req.ID == 0 {
		return apierrors.BadRequest("字典项ID不能为空")
	}

	_, err := h.systemService.DictItemSave(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return response.Success(c, nil)
}

// DictItemRemove deletes a dictionary item
// POST /api/v1/sysDictItem/remove
func (h *Handler) DictItemRemove(c echo.Context) error {
	var req dto.IDRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	if err := h.systemService.DictItemDelete(c.Request().Context(), req.ID); err != nil {
		return err
	}

	return response.Success(c, true)
}

// DictItemDetail returns dictionary item detail
// GET /api/v1/sysDictItem/detail
func (h *Handler) DictItemDetail(c echo.Context) error {
	idStr := c.QueryParam("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return apierrors.BadRequest("无效的字典项ID")
	}

	item, err := h.systemService.DictItemDetail(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return response.Success(c, item)
}
