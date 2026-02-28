package service

type CreateChannelReq struct {
	Name        string `json:"name" binding:"required"`
	ChannelType string `json:"channel_type" binding:"required"`
	Config      string `json:"config"`
}

type UpdateChannelReq struct {
	Name    string `json:"name"`
	Config  string `json:"config"`
	Enabled *bool  `json:"enabled"`
}

type CreateRuleReq struct {
	Name      string `json:"name" binding:"required"`
	EventType string `json:"event_type" binding:"required"`
	ChannelID string `json:"channel_id" binding:"required"`
	ProjectID string `json:"project_id"`
	Severity  string `json:"severity"`
}

type UpdateRuleReq struct {
	Name      string `json:"name"`
	EventType string `json:"event_type"`
	Severity  string `json:"severity"`
	Enabled   *bool  `json:"enabled"`
}

type CreateClusterReq struct {
	Name          string `json:"name" binding:"required"`
	DisplayName   string `json:"display_name"`
	Description   string `json:"description"`
	Provider      string `json:"provider"`
	APIServerURL  string `json:"api_server_url" binding:"required"`
	KubeConfigRef string `json:"kube_config_ref"`
}

type UpdateClusterReq struct {
	DisplayName   string `json:"display_name"`
	Description   string `json:"description"`
	KubeConfigRef string `json:"kube_config_ref"`
	Status        string `json:"status"`
}

type CreateIntegrationReq struct {
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Provider string `json:"provider" binding:"required"`
	Config   string `json:"config" binding:"required"`
}

type UpdateIntegrationReq struct {
	Name   string `json:"name"`
	Config string `json:"config"`
	Status string `json:"status"`
}
