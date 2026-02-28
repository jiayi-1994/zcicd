package engine

// RolloutStrategy for Argo Rollouts.
type RolloutStrategy struct {
	Type   string      `json:"type"`   // canary/bluegreen
	Config interface{} `json:"config"` // strategy-specific config
}

// CanaryConfig for canary deployments.
type CanaryConfig struct {
	Steps         []CanaryStep `json:"steps"`
	MaxSurge      string       `json:"max_surge"`
	MaxUnavailable string      `json:"max_unavailable"`
}

// CanaryStep defines a single canary step.
type CanaryStep struct {
	SetWeight *int   `json:"setWeight,omitempty"`
	Pause     *Pause `json:"pause,omitempty"`
}

// Pause duration for canary steps.
type Pause struct {
	Duration string `json:"duration,omitempty"` // e.g. "5m"
}

// BlueGreenConfig for blue-green deployments.
type BlueGreenConfig struct {
	AutoPromotionEnabled bool `json:"auto_promotion_enabled"`
	PreviewReplicaCount  int  `json:"preview_replica_count"`
}

// RolloutStatus from Argo Rollouts.
type RolloutStatus struct {
	Phase          string `json:"phase"`           // Healthy/Degraded/Paused/Progressing
	CurrentStep    int    `json:"current_step"`
	TotalSteps     int    `json:"total_steps"`
	StableRevision string `json:"stable_revision"`
	CanaryRevision string `json:"canary_revision"`
	Message        string `json:"message"`
}
