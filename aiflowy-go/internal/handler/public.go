package handler

import (
	"github.com/labstack/echo/v4"

	apierrors "github.com/aiflowy/aiflowy-go/internal/errors"
	"github.com/aiflowy/aiflowy-go/pkg/captcha"
	"github.com/aiflowy/aiflowy-go/pkg/response"
)

// PublicHandler handles public endpoints that don't require authentication
type PublicHandler struct{}

// NewPublicHandler creates a new PublicHandler
func NewPublicHandler() *PublicHandler {
	return &PublicHandler{}
}

// GetCaptcha generates and returns a captcha image
// GET /api/v1/public/getCaptcha
func (h *PublicHandler) GetCaptcha(c echo.Context) error {
	result, err := captcha.Generate()
	if err != nil {
		return apierrors.InternalError("生成验证码失败")
	}

	return response.Success(c, result)
}

// VerifyCaptcha verifies a captcha (mainly for testing)
// POST /api/v1/public/verifyCaptcha
func (h *PublicHandler) VerifyCaptcha(c echo.Context) error {
	var req struct {
		CaptchaID   string `json:"captchaId"`
		CaptchaCode string `json:"captchaCode"`
	}

	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	if captcha.Verify(req.CaptchaID, req.CaptchaCode) {
		return response.Success(c, map[string]bool{"valid": true})
	}

	return response.Success(c, map[string]bool{"valid": false})
}
