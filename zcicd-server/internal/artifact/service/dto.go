package service

type CreateRegistryReq struct {
	Name         string `json:"name" binding:"required"`
	RegistryType string `json:"registry_type"`
	Endpoint     string `json:"endpoint" binding:"required"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	IsDefault    bool   `json:"is_default"`
}

type UpdateRegistryReq struct {
	Name         string `json:"name"`
	RegistryType string `json:"registry_type"`
	Endpoint     string `json:"endpoint"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	IsDefault    *bool  `json:"is_default"`
}

type TriggerScanReq struct {
	RegistryID string `json:"registry_id" binding:"required"`
	Tag        string `json:"tag" binding:"required"`
}

type CreateChartReq struct {
	Name    string `json:"name" binding:"required"`
	RepoURL string `json:"repo_url" binding:"required"`
}
