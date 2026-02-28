package middleware

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/zcicd/zcicd-server/pkg/response"
)

// NewCasbinEnforcer creates a Casbin enforcer with the given model and policy adapter.
func NewCasbinEnforcer(modelPath string, policyAdapter interface{}) (*casbin.Enforcer, error) {
	e, err := casbin.NewEnforcer(modelPath, policyAdapter)
	if err != nil {
		return nil, err
	}
	if err := e.LoadPolicy(); err != nil {
		return nil, err
	}
	return e, nil
}

// CasbinRBAC enforces RBAC based on request path and method.
func CasbinRBAC(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			response.Unauthorized(c, "user not authenticated")
			c.Abort()
			return
		}

		ok, err := enforcer.Enforce(userID.(string), c.Request.URL.Path, c.Request.Method)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, 50000, "permission check failed")
			c.Abort()
			return
		}
		if !ok {
			response.Forbidden(c, "permission denied")
			c.Abort()
			return
		}

		c.Next()
	}
}

// httpMethodToAction maps an HTTP method to a project-level action.
func httpMethodToAction(method string) string {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return "read"
	case http.MethodPost:
		return "write"
	case http.MethodPut, http.MethodPatch:
		return "write"
	case http.MethodDelete:
		return "admin"
	default:
		return "read"
	}
}

// ProjectRBAC enforces RBAC scoped to a project extracted from the URL param ":project_id".
func ProjectRBAC(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			response.Unauthorized(c, "user not authenticated")
			c.Abort()
			return
		}

		projectID := c.Param("project_id")
		if projectID == "" {
			response.BadRequest(c, "missing project_id")
			c.Abort()
			return
		}

		action := httpMethodToAction(c.Request.Method)

		ok, err := enforcer.Enforce(userID.(string), projectID, action)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, 50000, "permission check failed")
			c.Abort()
			return
		}
		if !ok {
			response.Forbidden(c, "permission denied")
			c.Abort()
			return
		}

		c.Next()
	}
}
