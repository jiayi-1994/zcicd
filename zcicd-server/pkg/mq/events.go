package mq

// NATS JetStream event subjects
const (
	SubjectBuildStarted      = "zcicd.build.started"
	SubjectBuildCompleted    = "zcicd.build.completed"
	SubjectBuildFailed       = "zcicd.build.failed"
	SubjectDeploySyncing     = "zcicd.deploy.syncing"
	SubjectDeploySucceeded   = "zcicd.deploy.succeeded"
	SubjectDeployFailed      = "zcicd.deploy.failed"
	SubjectDeployRollback    = "zcicd.deploy.rollback"
	SubjectWorkflowStarted   = "zcicd.workflow.started"
	SubjectWorkflowApproval  = "zcicd.workflow.approval"
	SubjectWorkflowCompleted = "zcicd.workflow.completed"
	SubjectTestCompleted     = "zcicd.test.completed"
	SubjectScanCompleted     = "zcicd.scan.completed"
	SubjectGitOpsUpdate      = "zcicd.gitops.update"
	SubjectAuditLog          = "zcicd.audit.log"
)

// Event is the standard event envelope
type Event struct {
	EventID   string      `json:"event_id"`
	EventType string      `json:"event_type"`
	Timestamp string      `json:"timestamp"`
	ProjectID string      `json:"project_id"`
	TriggeredBy string    `json:"triggered_by"`
	Payload   interface{} `json:"payload"`
}
