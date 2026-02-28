package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/zcicd/zcicd-server/internal/workflow/service"
	"github.com/zcicd/zcicd-server/pkg/response"
)

type WebhookHandler struct {
	workflowSvc *service.WorkflowService
	buildSvc    *service.BuildService
}

func NewWebhookHandler(workflowSvc *service.WorkflowService, buildSvc *service.BuildService) *WebhookHandler {
	return &WebhookHandler{workflowSvc: workflowSvc, buildSvc: buildSvc}
}

type githubPushPayload struct {
	Ref        string `json:"ref"`
	After      string `json:"after"`
	Repository struct {
		CloneURL string `json:"clone_url"`
	} `json:"repository"`
	HeadCommit struct {
		Message string `json:"message"`
	} `json:"head_commit"`
}

type gitlabPushPayload struct {
	Ref        string `json:"ref"`
	After      string `json:"after"`
	Project    struct {
		GitHTTPURL string `json:"git_http_url"`
	} `json:"project"`
}

func (h *WebhookHandler) HandleGitHub(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.BadRequest(c, "failed to read body")
		return
	}

	// Verify signature if secret configured
	secret := c.Query("secret")
	if secret != "" {
		sig := c.GetHeader("X-Hub-Signature-256")
		if !verifyGitHubSignature(body, sig, secret) {
			response.Error(c, 403, 40301, "invalid signature")
			return
		}
	}

	event := c.GetHeader("X-GitHub-Event")
	if event != "push" {
		response.OK(c, gin.H{"message": "event ignored", "event": event})
		return
	}

	var payload githubPushPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		response.BadRequest(c, "invalid payload")
		return
	}

	branch := extractBranch(payload.Ref)
	h.triggerMatchingWorkflows(c, payload.Repository.CloneURL, branch, payload.After)
}

func (h *WebhookHandler) HandleGitLab(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.BadRequest(c, "failed to read body")
		return
	}

	// Verify token
	secret := c.Query("secret")
	if secret != "" {
		token := c.GetHeader("X-Gitlab-Token")
		if token != secret {
			response.Error(c, 403, 40301, "invalid token")
			return
		}
	}

	event := c.GetHeader("X-Gitlab-Event")
	if event != "Push Hook" {
		response.OK(c, gin.H{"message": "event ignored", "event": event})
		return
	}

	var payload gitlabPushPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		response.BadRequest(c, "invalid payload")
		return
	}

	branch := extractBranch(payload.Ref)
	h.triggerMatchingWorkflows(c, payload.Project.GitHTTPURL, branch, payload.After)
}

func (h *WebhookHandler) triggerMatchingWorkflows(c *gin.Context, repoURL, branch, commitSHA string) {
	triggered, err := h.workflowSvc.TriggerByWebhook(c.Request.Context(), repoURL, branch, commitSHA)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"triggered": triggered})
}

func extractBranch(ref string) string {
	// refs/heads/main -> main
	if len(ref) > 11 {
		return ref[11:]
	}
	return ref
}

func verifyGitHubSignature(body []byte, signature, secret string) bool {
	if len(signature) < 8 {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}
