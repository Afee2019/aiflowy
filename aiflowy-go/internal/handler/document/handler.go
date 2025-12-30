package document

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/aiflowy/aiflowy-go/internal/config"
	"github.com/aiflowy/aiflowy-go/internal/dto"
	apierrors "github.com/aiflowy/aiflowy-go/internal/errors"
	"github.com/aiflowy/aiflowy-go/internal/handler/auth"
	"github.com/aiflowy/aiflowy-go/internal/service"
	"github.com/aiflowy/aiflowy-go/pkg/response"
)

// Handler 知识库和文档 API Handler
type Handler struct {
	collectionService *service.DocumentCollectionService
	documentService   *service.DocumentService
}

// NewHandler 创建 Handler
func NewHandler() *Handler {
	return &Handler{
		collectionService: service.NewDocumentCollectionService(),
		documentService:   service.NewDocumentService(),
	}
}

// helper to get user context from JWT claims
func getUserContext(c echo.Context) (userID, tenantID, deptID int64) {
	claims := auth.GetClaims(c)
	if claims == nil {
		return 0, 0, 0
	}
	userID, _ = strconv.ParseInt(claims.UserID, 10, 64)
	tenantID, _ = strconv.ParseInt(claims.TenantID, 10, 64)
	deptID, _ = strconv.ParseInt(claims.DeptID, 10, 64)
	return
}

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(g *echo.Group) {
	// 知识库 CRUD
	documentCollection := g.Group("/documentCollection")
	documentCollection.GET("/list", h.ListDocumentCollections)
	documentCollection.POST("/list", h.ListDocumentCollections)
	documentCollection.GET("/detail", h.GetDocumentCollection)
	documentCollection.POST("/save", h.SaveDocumentCollection)
	documentCollection.POST("/remove", h.DeleteDocumentCollection)

	// 文档 CRUD
	document := g.Group("/document")
	document.GET("/list", h.ListDocuments)
	document.POST("/list", h.ListDocuments)
	document.GET("/documentList", h.ListDocumentsPaged)
	document.GET("/detail", h.GetDocument)
	document.POST("/save", h.SaveDocument)
	document.POST("/update", h.UpdateDocument)
	document.POST("/removeDoc", h.DeleteDocument)
	document.POST("/remove", h.DeleteDocument)
	document.GET("/download", h.DownloadDocument)
	document.POST("/textSplit", h.TextSplit)
	document.POST("/saveText", h.TextSplit)

	// Bot-知识库关联
	botKnowledge := g.Group("/botKnowledge")
	botKnowledge.GET("/list", h.ListBotKnowledges)
	botKnowledge.POST("/list", h.ListBotKnowledges)
	botKnowledge.POST("/updateBotKnowledgeIds", h.UpdateBotKnowledgeIDs)
	botKnowledge.POST("/getBotKnowledgeIds", h.GetBotKnowledgeIDs)
	botKnowledge.POST("/remove", h.DeleteBotKnowledge)

	// 文件上传
	commons := g.Group("/commons")
	commons.POST("/upload", h.Upload)
	commons.POST("/uploadAntd", h.Upload)
	commons.POST("/uploadPrePath", h.UploadPrePath)
}

// ========================== 知识库 CRUD ==========================

// ListDocumentCollections 获取知识库列表
func (h *Handler) ListDocumentCollections(c echo.Context) error {
	ctx := c.Request().Context()
	_, tenantID, _ := getUserContext(c)

	collections, err := h.collectionService.List(ctx, tenantID)
	if err != nil {
		return err
	}
	return response.Success(c, collections)
}

// GetDocumentCollection 获取知识库详情
func (h *Handler) GetDocumentCollection(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.QueryParam("id")

	collection, err := h.collectionService.GetDetail(ctx, id)
	if err != nil {
		return err
	}
	if collection == nil {
		return apierrors.NotFound("知识库不存在")
	}
	return response.Success(c, collection)
}

// SaveDocumentCollection 保存知识库
func (h *Handler) SaveDocumentCollection(c echo.Context) error {
	ctx := c.Request().Context()
	userID, tenantID, deptID := getUserContext(c)

	var req dto.DocumentCollectionSaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	collection, err := h.collectionService.Save(ctx, &req, tenantID, userID, deptID)
	if err != nil {
		return apierrors.BadRequest(err.Error())
	}
	return response.Success(c, collection)
}

// DeleteDocumentCollection 删除知识库
func (h *Handler) DeleteDocumentCollection(c echo.Context) error {
	ctx := c.Request().Context()

	var req struct {
		IDs []string `json:"ids"`
	}
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	for _, id := range req.IDs {
		if err := h.collectionService.Delete(ctx, id); err != nil {
			return apierrors.BadRequest(err.Error())
		}
	}
	return response.Success(c, true)
}

// ========================== 文档 CRUD ==========================

// ListDocuments 获取文档列表
func (h *Handler) ListDocuments(c echo.Context) error {
	ctx := c.Request().Context()

	// 从 query 或 body 获取知识库 ID
	id := c.QueryParam("id")
	if id == "" {
		var req struct {
			ID string `json:"id"`
		}
		c.Bind(&req)
		id = req.ID
	}

	if id == "" {
		return apierrors.BadRequest("知识库 ID 不能为空")
	}

	documents, err := h.documentService.List(ctx, id)
	if err != nil {
		return err
	}
	return response.Success(c, documents)
}

// ListDocumentsPaged 分页获取文档列表
func (h *Handler) ListDocumentsPaged(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.QueryParam("id")
	if id == "" {
		return apierrors.BadRequest("知识库 ID 不能为空")
	}

	title := c.QueryParam("title")
	pageNumber, _ := strconv.Atoi(c.QueryParam("pageNumber"))
	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))

	result, err := h.documentService.ListPaged(ctx, id, title, pageNumber, pageSize)
	if err != nil {
		return err
	}
	return response.Success(c, result)
}

// GetDocument 获取文档详情
func (h *Handler) GetDocument(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.QueryParam("id")

	document, err := h.documentService.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if document == nil {
		return apierrors.NotFound("文档不存在")
	}
	return response.Success(c, document)
}

// SaveDocument 保存文档
func (h *Handler) SaveDocument(c echo.Context) error {
	ctx := c.Request().Context()
	userID, _, _ := getUserContext(c)

	var req dto.DocumentSaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	if req.CollectionID == "" {
		return apierrors.BadRequest("知识库 ID 不能为空")
	}

	document, err := h.documentService.Save(ctx, &req, userID)
	if err != nil {
		return apierrors.BadRequest(err.Error())
	}
	return response.Success(c, document)
}

// UpdateDocument 更新文档
func (h *Handler) UpdateDocument(c echo.Context) error {
	ctx := c.Request().Context()
	userID, _, _ := getUserContext(c)

	var req dto.DocumentSaveRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	if req.ID == "" {
		return apierrors.BadRequest("文档 ID 不能为空")
	}

	// 如果有排序更新
	if req.OrderNo != nil {
		if err := h.documentService.UpdatePosition(ctx, req.ID, *req.OrderNo, userID); err != nil {
			return apierrors.BadRequest(err.Error())
		}
	}

	_, err := h.documentService.Save(ctx, &req, userID)
	if err != nil {
		return apierrors.BadRequest(err.Error())
	}
	return response.Success(c, true)
}

// DeleteDocument 删除文档
func (h *Handler) DeleteDocument(c echo.Context) error {
	ctx := c.Request().Context()

	var req struct {
		ID  string   `json:"id"`
		IDs []string `json:"ids"`
	}
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	// 支持单个 ID 或多个 ID
	ids := req.IDs
	if req.ID != "" {
		ids = []string{req.ID}
	}

	for _, id := range ids {
		if err := h.documentService.Delete(ctx, id); err != nil {
			return apierrors.BadRequest(err.Error())
		}
	}
	return response.Success(c, true)
}

// DownloadDocument 下载文档
func (h *Handler) DownloadDocument(c echo.Context) error {
	ctx := c.Request().Context()
	documentID := c.QueryParam("documentId")

	document, err := h.documentService.GetByID(ctx, documentID)
	if err != nil {
		return err
	}
	if document == nil {
		return apierrors.NotFound("文档不存在")
	}

	// 获取文件路径
	rootPath := getUploadRootPath()
	filePath := filepath.Join(rootPath, document.DocumentPath)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return apierrors.NotFound("文件不存在")
	}

	// 构造文件名
	filename := document.Title
	if document.DocumentType != "" {
		filename = fmt.Sprintf("%s.%s", document.Title, document.DocumentType)
	}

	return c.Attachment(filePath, filename)
}

// TextSplit 文本拆分
func (h *Handler) TextSplit(c echo.Context) error {
	ctx := c.Request().Context()
	userID, _, _ := getUserContext(c)

	var req dto.TextSplitRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	result, err := h.documentService.TextSplit(ctx, &req, userID)
	if err != nil {
		return apierrors.BadRequest(err.Error())
	}
	return response.Success(c, result)
}

// ========================== Bot-知识库关联 ==========================

// ListBotKnowledges 获取 Bot 关联的知识库列表
func (h *Handler) ListBotKnowledges(c echo.Context) error {
	ctx := c.Request().Context()

	botID := c.QueryParam("botId")
	if botID == "" {
		var req struct {
			BotID string `json:"botId"`
		}
		c.Bind(&req)
		botID = req.BotID
	}

	if botID == "" {
		return apierrors.BadRequest("Bot ID 不能为空")
	}

	knowledges, err := h.collectionService.ListByBotID(ctx, botID)
	if err != nil {
		return err
	}
	return response.Success(c, knowledges)
}

// GetBotKnowledgeIDs 获取 Bot 关联的知识库 ID 列表
func (h *Handler) GetBotKnowledgeIDs(c echo.Context) error {
	ctx := c.Request().Context()

	var req struct {
		BotID string `json:"botId"`
	}
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	ids, err := h.collectionService.GetBotKnowledgeIDs(ctx, req.BotID)
	if err != nil {
		return err
	}
	return response.Success(c, ids)
}

// UpdateBotKnowledgeIDs 更新 Bot-知识库关联
func (h *Handler) UpdateBotKnowledgeIDs(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.BotDocumentCollectionUpdateRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	if err := h.collectionService.UpdateBotKnowledges(ctx, &req); err != nil {
		return apierrors.BadRequest(err.Error())
	}
	return response.Success(c, nil)
}

// DeleteBotKnowledge 删除 Bot-知识库关联
func (h *Handler) DeleteBotKnowledge(c echo.Context) error {
	ctx := c.Request().Context()

	var req struct {
		BotID       string `json:"botId"`
		KnowledgeID string `json:"knowledgeId"`
	}
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	if err := h.collectionService.DeleteBotKnowledge(ctx, req.BotID, req.KnowledgeID); err != nil {
		return apierrors.BadRequest(err.Error())
	}
	return response.Success(c, true)
}

// ========================== 文件上传 ==========================

// Upload 上传文件
func (h *Handler) Upload(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return apierrors.BadRequest("请上传文件")
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		return apierrors.InternalError("打开文件失败")
	}
	defer src.Close()

	// 生成保存路径
	rootPath := getUploadRootPath()
	relativePath := generateFilePath(file.Filename, "")
	fullPath := filepath.Join(rootPath, relativePath)

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return apierrors.InternalError("创建目录失败")
	}

	// 创建目标文件
	dst, err := os.Create(fullPath)
	if err != nil {
		return apierrors.InternalError("创建文件失败")
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, src); err != nil {
		return apierrors.InternalError("保存文件失败")
	}

	return response.Success(c, dto.UploadResponse{Path: relativePath})
}

// UploadPrePath 上传文件到指定路径
func (h *Handler) UploadPrePath(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return apierrors.BadRequest("请上传文件")
	}

	prePath := c.FormValue("prePath")

	// 打开文件
	src, err := file.Open()
	if err != nil {
		return apierrors.InternalError("打开文件失败")
	}
	defer src.Close()

	// 生成保存路径
	rootPath := getUploadRootPath()
	relativePath := generateFilePath(file.Filename, prePath)
	fullPath := filepath.Join(rootPath, relativePath)

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return apierrors.InternalError("创建目录失败")
	}

	// 创建目标文件
	dst, err := os.Create(fullPath)
	if err != nil {
		return apierrors.InternalError("创建文件失败")
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, src); err != nil {
		return apierrors.InternalError("保存文件失败")
	}

	return response.Success(c, dto.UploadResponse{Path: relativePath})
}

// ========================== 辅助函数 ==========================

// getUploadRootPath 获取上传根路径
func getUploadRootPath() string {
	cfg := config.GetConfig()
	if cfg.Storage.LocalRoot != "" {
		return cfg.Storage.LocalRoot
	}
	// 默认使用当前目录下的 uploads
	return "./uploads"
}

// generateFilePath 生成文件保存路径
func generateFilePath(filename, prePath string) string {
	// 获取文件扩展名
	ext := filepath.Ext(filename)

	// 生成日期目录
	now := time.Now()
	dateDir := now.Format("2006/01/02")

	// 生成唯一文件名
	uniqueName := fmt.Sprintf("%d%s", now.UnixNano(), ext)

	if prePath != "" {
		prePath = strings.TrimPrefix(prePath, "/")
		prePath = strings.TrimSuffix(prePath, "/")
		return filepath.Join(prePath, dateDir, uniqueName)
	}

	return filepath.Join(dateDir, uniqueName)
}
