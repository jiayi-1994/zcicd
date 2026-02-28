package service

import "encoding/json"

type CreateProjectRequest struct {
	Name          string `json:"name" binding:"required,min=2,max=100"`
	Identifier    string `json:"identifier" binding:"required,min=2,max=50,alphanum"`
	Description   string `json:"description"`
	RepoURL       string `json:"repo_url"`
	DefaultBranch string `json:"default_branch"`
	Visibility    string `json:"visibility" binding:"omitempty,oneof=public private"`
}

type UpdateProjectRequest struct {
	Name          string `json:"name" binding:"omitempty,min=2,max=100"`
	Description   string `json:"description"`
	RepoURL       string `json:"repo_url"`
	DefaultBranch string `json:"default_branch"`
	Visibility    string `json:"visibility" binding:"omitempty,oneof=public private"`
	Status        string `json:"status" binding:"omitempty,oneof=active archived"`
}

type CreateServiceRequest struct {
	Name            string          `json:"name" binding:"required,min=2,max=100"`
	ServiceType     string          `json:"service_type" binding:"required,oneof=backend frontend middleware"`
	Language        string          `json:"language"`
	RepoURL         string          `json:"repo_url"`
	Branch          string          `json:"branch"`
	DockerfilePath  string          `json:"dockerfile_path"`
	BuildContext    string          `json:"build_context"`
	DeployType      string          `json:"deploy_type" binding:"omitempty,oneof=helm k8s_yaml kustomize"`
	HelmChartPath   string          `json:"helm_chart_path"`
	HelmValues      json.RawMessage `json:"helm_values"`
	HealthCheckPath string          `json:"health_check_path"`
	Ports           json.RawMessage `json:"ports"`
	EnvVars         json.RawMessage `json:"env_vars"`
	Resources       json.RawMessage `json:"resources"`
}

type UpdateServiceRequest struct {
	Name            string          `json:"name" binding:"omitempty,min=2,max=100"`
	ServiceType     string          `json:"service_type" binding:"omitempty,oneof=backend frontend middleware"`
	Language        string          `json:"language"`
	RepoURL         string          `json:"repo_url"`
	Branch          string          `json:"branch"`
	DockerfilePath  string          `json:"dockerfile_path"`
	BuildContext    string          `json:"build_context"`
	DeployType      string          `json:"deploy_type" binding:"omitempty,oneof=helm k8s_yaml kustomize"`
	HelmChartPath   string          `json:"helm_chart_path"`
	HelmValues      json.RawMessage `json:"helm_values"`
	HealthCheckPath string          `json:"health_check_path"`
	Ports           json.RawMessage `json:"ports"`
	EnvVars         json.RawMessage `json:"env_vars"`
	Resources       json.RawMessage `json:"resources"`
	Status          string          `json:"status" binding:"omitempty,oneof=active inactive"`
}

type CreateEnvRequest struct {
	Name           string          `json:"name" binding:"required,min=2,max=50"`
	EnvType        string          `json:"env_type" binding:"required,oneof=dev testing staging production"`
	Namespace      string          `json:"namespace"`
	ClusterID      string          `json:"cluster_id"`
	IsProduction   bool            `json:"is_production"`
	AutoDeploy     bool            `json:"auto_deploy"`
	DeployStrategy json.RawMessage `json:"deploy_strategy"`
	GlobalEnvVars  json.RawMessage `json:"global_env_vars"`
}

type UpdateEnvRequest struct {
	Name           string          `json:"name" binding:"omitempty,min=2,max=50"`
	EnvType        string          `json:"env_type" binding:"omitempty,oneof=dev testing staging production"`
	Namespace      string          `json:"namespace"`
	ClusterID      string          `json:"cluster_id"`
	IsProduction   *bool           `json:"is_production"`
	AutoDeploy     *bool           `json:"auto_deploy"`
	DeployStrategy json.RawMessage `json:"deploy_strategy"`
	GlobalEnvVars  json.RawMessage `json:"global_env_vars"`
	Status         string          `json:"status" binding:"omitempty,oneof=active inactive"`
}
