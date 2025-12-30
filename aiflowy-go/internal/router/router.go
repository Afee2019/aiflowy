package router

import (
	"github.com/labstack/echo/v4"

	"github.com/aiflowy/aiflowy-go/internal/handler"
	"github.com/aiflowy/aiflowy-go/internal/handler/ai"
	"github.com/aiflowy/aiflowy-go/internal/handler/auth"
	"github.com/aiflowy/aiflowy-go/internal/handler/bot"
	"github.com/aiflowy/aiflowy-go/internal/handler/document"
	"github.com/aiflowy/aiflowy-go/internal/handler/model"
	"github.com/aiflowy/aiflowy-go/internal/handler/plugin"
	"github.com/aiflowy/aiflowy-go/internal/handler/system"
	"github.com/aiflowy/aiflowy-go/internal/handler/workflow"
	"github.com/aiflowy/aiflowy-go/internal/middleware"
	"github.com/aiflowy/aiflowy-go/pkg/metrics"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(e *echo.Echo) {
	// Health check endpoints
	healthHandler := handler.NewHealthHandler()
	e.GET("/health", healthHandler.Health)
	e.GET("/health/detail", healthHandler.HealthDetail)

	// Prometheus metrics endpoint
	e.GET("/metrics", metrics.MetricsHandler())

	// Test endpoints (for development/verification)
	testHandler := handler.NewTestHandler()
	test := e.Group("/test")
	test.GET("/error", testHandler.TestError)
	test.GET("/panic", testHandler.TestPanic)
	test.GET("/snowflake", testHandler.TestSnowflake)
	test.GET("/config", testHandler.TestConfig)

	// API v1 group
	apiV1 := e.Group("/api/v1")

	// Public routes (no auth required)
	publicHandler := handler.NewPublicHandler()
	public := apiV1.Group("/public")
	public.GET("/getCaptcha", publicHandler.GetCaptcha)
	public.POST("/verifyCaptcha", publicHandler.VerifyCaptcha)

	// Auth routes (no auth required for login)
	authHandler := auth.NewHandler()
	authGroup := apiV1.Group("/auth")
	authGroup.POST("/login", authHandler.Login)
	authGroup.POST("/logout", authHandler.Logout)

	// Auth routes that require authentication
	authProtected := apiV1.Group("/auth")
	authProtected.Use(middleware.JWTAuth())
	authProtected.GET("/getPermissions", authHandler.GetPermissions)
	authProtected.GET("/getRoles", authHandler.GetRoles)
	authProtected.GET("/getMenus", authHandler.GetMenus)
	authProtected.GET("/getUserInfo", authHandler.GetUserInfo)

	// System management routes (all require authentication)
	systemHandler := system.NewHandler()

	// Account management
	sysAccount := apiV1.Group("/sysAccount")
	sysAccount.Use(middleware.JWTAuth())
	sysAccount.GET("/page", systemHandler.AccountPage)
	sysAccount.GET("/list", systemHandler.AccountList)
	sysAccount.GET("/detail", systemHandler.AccountDetail)
	sysAccount.POST("/save", systemHandler.AccountSave)
	sysAccount.POST("/update", systemHandler.AccountUpdate)
	sysAccount.POST("/remove", systemHandler.AccountRemove)
	sysAccount.POST("/updatePassword", systemHandler.UpdatePassword)
	sysAccount.GET("/myProfile", systemHandler.MyProfile)

	// Role management
	sysRole := apiV1.Group("/sysRole")
	sysRole.Use(middleware.JWTAuth())
	sysRole.GET("/page", systemHandler.RolePage)
	sysRole.GET("/list", systemHandler.RoleList)
	sysRole.GET("/detail", systemHandler.RoleDetail)
	sysRole.POST("/save", systemHandler.RoleSave)
	sysRole.POST("/update", systemHandler.RoleUpdate)
	sysRole.POST("/remove", systemHandler.RoleRemove)
	sysRole.GET("/getRoleMenuIds", systemHandler.GetRoleMenuIds)

	// Menu management
	sysMenu := apiV1.Group("/sysMenu")
	sysMenu.Use(middleware.JWTAuth())
	sysMenu.GET("/list", systemHandler.MenuList)
	sysMenu.GET("/detail", systemHandler.MenuDetail)
	sysMenu.POST("/save", systemHandler.MenuSave)
	sysMenu.POST("/update", systemHandler.MenuUpdate)
	sysMenu.POST("/remove", systemHandler.MenuRemove)
	sysMenu.GET("/getCheckedByRoleId/:roleId", systemHandler.GetMenuCheckedByRoleId)

	// Department management
	sysDept := apiV1.Group("/sysDept")
	sysDept.Use(middleware.JWTAuth())
	sysDept.GET("/list", systemHandler.DeptList)
	sysDept.GET("/detail", systemHandler.DeptDetail)
	sysDept.POST("/save", systemHandler.DeptSave)
	sysDept.POST("/update", systemHandler.DeptUpdate)
	sysDept.POST("/remove", systemHandler.DeptRemove)

	// Dictionary management
	sysDict := apiV1.Group("/sysDict")
	sysDict.Use(middleware.JWTAuth())
	sysDict.GET("/page", systemHandler.DictPage)
	sysDict.GET("/list", systemHandler.DictList)
	sysDict.GET("/detail", systemHandler.DictDetail)
	sysDict.POST("/save", systemHandler.DictSave)
	sysDict.POST("/update", systemHandler.DictUpdate)
	sysDict.POST("/remove", systemHandler.DictRemove)

	// Dictionary Item management
	sysDictItem := apiV1.Group("/sysDictItem")
	sysDictItem.Use(middleware.JWTAuth())
	sysDictItem.GET("/page", systemHandler.DictItemPage)
	sysDictItem.GET("/list", systemHandler.DictItemList)
	sysDictItem.GET("/detail", systemHandler.DictItemDetail)
	sysDictItem.POST("/save", systemHandler.DictItemSave)
	sysDictItem.POST("/update", systemHandler.DictItemUpdate)
	sysDictItem.POST("/remove", systemHandler.DictItemRemove)

	// Model management routes (all require authentication)
	modelHandler := model.NewHandler()

	// Model Provider management
	modelProvider := apiV1.Group("/modelProvider")
	modelProvider.Use(middleware.JWTAuth())
	modelProvider.GET("/page", modelHandler.ProviderPage)
	modelProvider.GET("/list", modelHandler.ProviderList)
	modelProvider.GET("/detail", modelHandler.ProviderDetail)
	modelProvider.POST("/save", modelHandler.ProviderSave)
	modelProvider.POST("/update", modelHandler.ProviderUpdate)
	modelProvider.POST("/remove", modelHandler.ProviderRemove)

	// Model management
	modelGroup := apiV1.Group("/model")
	modelGroup.Use(middleware.JWTAuth())
	modelGroup.GET("/page", modelHandler.ModelPage)
	modelGroup.GET("/list", modelHandler.ModelList)
	modelGroup.GET("/detail", modelHandler.ModelDetail)
	modelGroup.POST("/save", modelHandler.ModelSave)
	modelGroup.POST("/update", modelHandler.ModelUpdate)
	modelGroup.POST("/remove", modelHandler.ModelRemove)
	modelGroup.GET("/getList", modelHandler.GetList)
	modelGroup.GET("/selectLlmByProviderAndModelType", modelHandler.SelectLlmByProviderAndModelType)
	modelGroup.GET("/selectLlmList", modelHandler.SelectLlmList)
	modelGroup.POST("/addAiLlm", modelHandler.AddAiLlm)
	modelGroup.POST("/addAllLlm", modelHandler.AddAllLlm)
	modelGroup.POST("/updateByEntity", modelHandler.UpdateByEntity)
	modelGroup.POST("/removeByEntity", modelHandler.RemoveByEntity)
	modelGroup.POST("/removeLlmByIds", modelHandler.RemoveLlmByIds)
	modelGroup.GET("/verifyLlmConfig", modelHandler.VerifyLlmConfig)

	// AI routes (all require authentication)
	aiHandler := ai.NewHandler()
	aiGroup := apiV1.Group("/ai")
	aiGroup.Use(middleware.JWTAuth())
	aiGroup.GET("/test", aiHandler.Test)
	aiGroup.POST("/test", aiHandler.Test)
	aiGroup.POST("/chat", aiHandler.Chat)

	// Bot routes (all require authentication)
	botHandler := bot.NewHandler()

	// Bot management
	botGroup := apiV1.Group("/bot")
	botGroup.Use(middleware.JWTAuth())
	botGroup.GET("/page", botHandler.BotPage)
	botGroup.GET("/list", botHandler.BotList)
	botGroup.GET("/detail", botHandler.BotDetail)
	botGroup.GET("/getDetail", botHandler.GetDetail)
	botGroup.POST("/save", botHandler.BotSave)
	botGroup.POST("/update", botHandler.BotUpdate)
	botGroup.POST("/updateLlmOptions", botHandler.UpdateLlmOptions)
	botGroup.POST("/remove", botHandler.BotRemove)
	botGroup.GET("/generateConversationId", botHandler.GenerateConversationId)
	botGroup.POST("/chat", botHandler.Chat) // Bot streaming chat API
	botGroup.POST("/voiceInput", botHandler.VoiceInput)
	botGroup.POST("/prompt/chore/chat", botHandler.PromptChoreChat)

	// Bot API Key management
	botApiKeyHandler := bot.NewBotApiKeyHandler()
	botApiKey := apiV1.Group("/botApiKey")
	botApiKey.Use(middleware.JWTAuth())
	botApiKeyHandler.Register(botApiKey)

	// Bot Category management
	botCategory := apiV1.Group("/botCategory")
	botCategory.Use(middleware.JWTAuth())
	botCategory.GET("/list", botHandler.CategoryList)
	botCategory.GET("/detail", botHandler.CategoryDetail)
	botCategory.POST("/save", botHandler.CategorySave)
	botCategory.POST("/update", botHandler.CategoryUpdate)
	botCategory.POST("/remove", botHandler.CategoryRemove)

	// Bot Conversation management
	botConversation := apiV1.Group("/botConversation")
	botConversation.Use(middleware.JWTAuth())
	botConversation.GET("/page", botHandler.ConversationPage)
	botConversation.GET("/list", botHandler.ConversationList)
	botConversation.GET("/detail", botHandler.ConversationDetail)
	botConversation.POST("/save", botHandler.ConversationSave)
	botConversation.POST("/update", botHandler.ConversationUpdate)
	botConversation.POST("/remove", botHandler.ConversationRemove)

	// Bot Message management
	botMessage := apiV1.Group("/botMessage")
	botMessage.Use(middleware.JWTAuth())
	botMessage.GET("/page", botHandler.MessagePage)
	botMessage.GET("/list", botHandler.MessageList)
	botMessage.GET("/detail", botHandler.MessageDetail)
	botMessage.POST("/save", botHandler.MessageSave)
	botMessage.POST("/update", botHandler.MessageUpdate)
	botMessage.POST("/remove", botHandler.MessageRemove)

	// Plugin routes (all require authentication)
	pluginHandler := plugin.NewHandler()
	protectedGroup := apiV1.Group("")
	protectedGroup.Use(middleware.JWTAuth())
	pluginHandler.RegisterRoutes(protectedGroup)

	// Workflow routes (all require authentication)
	workflowHandler := workflow.NewHandler()
	workflowHandler.RegisterRoutes(protectedGroup)

	// Document/Knowledge routes (all require authentication)
	documentHandler := document.NewHandler()
	documentHandler.RegisterRoutes(protectedGroup)

	// System auxiliary routes (all require authentication)
	// API Key management
	apiKeyHandler := system.NewSysApiKeyHandler()
	sysApiKey := apiV1.Group("/sysApiKey")
	sysApiKey.Use(middleware.JWTAuth())
	apiKeyHandler.Register(sysApiKey)

	// Log management
	logHandler := system.NewSysLogHandler()
	sysLog := apiV1.Group("/sysLog")
	sysLog.Use(middleware.JWTAuth())
	logHandler.Register(sysLog)

	// Option management
	optionHandler := system.NewSysOptionHandler()
	sysOption := apiV1.Group("/sysOption")
	sysOption.Use(middleware.JWTAuth())
	optionHandler.Register(sysOption)

	// Job management
	jobHandler := system.NewSysJobHandler()
	sysJob := apiV1.Group("/sysJob")
	sysJob.Use(middleware.JWTAuth())
	jobHandler.Register(sysJob)
}
