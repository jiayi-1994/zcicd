package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	RequestID string      `json:"request_id"`
}

type PaginatedResponse struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
	RequestID  string      `json:"request_id"`
}

type Pagination struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

func getRequestID(c *gin.Context) string {
	if id := c.GetString("request_id"); id != "" {
		return id
	}
	return uuid.New().String()
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:      0,
		Message:   "success",
		Data:      data,
		RequestID: getRequestID(c),
	})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:      0,
		Message:   "created",
		Data:      data,
		RequestID: getRequestID(c),
	})
}

func OKWithPage(c *gin.Context, data interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, PaginatedResponse{
		Code:    0,
		Message: "success",
		Data:    data,
		Pagination: &Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
		RequestID: getRequestID(c),
	})
}

func Error(c *gin.Context, httpCode int, errCode int, message string) {
	c.JSON(httpCode, Response{
		Code:      errCode,
		Message:   message,
		RequestID: getRequestID(c),
	})
}

func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, 40000, message)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, 40100, message)
}

func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, 40300, message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, 40400, message)
}

func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, 50000, message)
}
