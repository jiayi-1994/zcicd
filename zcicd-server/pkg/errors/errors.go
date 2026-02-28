package errors

import "fmt"

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func New(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func NewAppError(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func Wrap(code int, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

// Common error codes
var (
	ErrBadRequest       = New(40000, "请求参数错误")
	ErrUnauthorized     = New(40100, "未认证")
	ErrTokenExpired     = New(40101, "Token 已过期")
	ErrTokenInvalid     = New(40102, "Token 无效")
	ErrForbidden        = New(40300, "无权限")
	ErrNotFound         = New(40400, "资源不存在")
	ErrConflict         = New(40900, "资源冲突")
	ErrInternal         = New(50000, "服务器内部错误")
	ErrDatabaseError    = New(50001, "数据库错误")
	ErrExternalService  = New(50002, "外部服务调用失败")

	// Auth errors: 401xx
	ErrLoginFailed      = New(40103, "用户名或密码错误")
	ErrUserDisabled     = New(40104, "用户已禁用")

	// Project errors: 402xx
	ErrProjectExists    = New(40201, "项目名称已存在")
	ErrProjectNotFound  = New(40202, "项目不存在")

	// Service errors: 403xx
	ErrServiceExists    = New(40301, "服务名称已存在")
	ErrServiceNotFound  = New(40302, "服务不存在")

	// Environment errors: 404xx
	ErrEnvExists        = New(40401, "环境名称已存在")
	ErrEnvNotFound      = New(40402, "环境不存在")

	// Workflow errors: 405xx
	ErrWorkflowNotFound   = New(40501, "工作流不存在")
	ErrWorkflowRunNotFound = New(40502, "工作流运行不存在")

	// Build errors: 406xx
	ErrBuildConfigNotFound = New(40601, "构建配置不存在")
	ErrBuildRunNotFound    = New(40602, "构建运行不存在")

	// Deploy errors: 407xx
	ErrDeployConfigNotFound  = New(40701, "部署配置不存在")
	ErrDeployConfigExists    = New(40702, "部署配置已存在")
	ErrDeployHistoryNotFound = New(40703, "部署记录不存在")
	ErrDeployInProgress      = New(40704, "部署正在进行中")
	ErrDeploySyncFailed      = New(40705, "部署同步失败")
	ErrRollbackFailed        = New(40706, "回滚失败")

	// Approval errors: 408xx
	ErrApprovalNotFound       = New(40801, "审批记录不存在")
	ErrApprovalRequired       = New(40802, "生产环境部署需要审批")
	ErrApprovalAlreadyDecided = New(40803, "审批已处理")
	ErrNotApprover            = New(40804, "无审批权限")

	// Environment errors (additional): 404xx
	ErrEnvVariableNotFound = New(40403, "环境变量不存在")
	ErrEnvQuotaExceeded    = New(40404, "环境资源配额超限")
)
