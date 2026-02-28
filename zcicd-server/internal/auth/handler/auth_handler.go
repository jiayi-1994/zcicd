package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zcicd/zcicd-server/internal/auth/service"
	apperrors "github.com/zcicd/zcicd-server/pkg/errors"
	"github.com/zcicd/zcicd-server/pkg/response"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := h.svc.Register(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Created(c, user)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	token, err := h.svc.Login(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, token)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req service.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	token, err := h.svc.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, token)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	user, err := h.svc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, user)
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	var req service.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.UpdateProfile(c.Request.Context(), userID, &req); err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, nil)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID := c.GetString("user_id")
	var req service.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.ChangePassword(c.Request.Context(), userID, &req); err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, nil)
}

func (h *AuthHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	users, total, err := h.svc.ListUsers(c.Request.Context(), page, pageSize)
	if err != nil {
		handleError(c, err)
		return
	}
	response.OKWithPage(c, users, total, page, pageSize)
}

func handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*apperrors.AppError); ok {
		httpCode := http.StatusInternalServerError
		switch {
		case appErr.Code >= 40000 && appErr.Code < 40100:
			httpCode = http.StatusBadRequest
		case appErr.Code >= 40100 && appErr.Code < 40200:
			httpCode = http.StatusUnauthorized
		case appErr.Code >= 40300 && appErr.Code < 40400:
			httpCode = http.StatusForbidden
		case appErr.Code >= 40400 && appErr.Code < 40500:
			httpCode = http.StatusNotFound
		case appErr.Code >= 40900 && appErr.Code < 41000:
			httpCode = http.StatusConflict
		}
		response.Error(c, httpCode, appErr.Code, appErr.Message)
		return
	}
	response.InternalError(c, "服务器内部错误")
}
