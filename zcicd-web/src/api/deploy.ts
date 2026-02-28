import request from '@/utils/request'

export interface DeployConfig {
  id: string
  project_id: string
  service_id: string
  env_id: string
  argocd_app_name: string
  git_repo_url: string
  git_path: string
  git_branch: string
  deploy_strategy: string
  auto_sync: boolean
  require_approval: boolean
  values_override: Record<string, unknown>
  created_at: string
  updated_at: string
}

export interface DeployHistory {
  id: string
  deploy_config_id: string
  version: number
  status: string
  image_tag: string
  commit_sha: string
  sync_status: string
  health_status: string
  deployed_by: string
  started_at: string
  finished_at: string
  created_at: string
}

export interface ApprovalRecord {
  id: string
  deploy_config_id: string
  status: string
  requested_by: string
  approved_by: string
  comment: string
  created_at: string
}

export interface EnvVariable {
  id: string
  env_id: string
  key: string
  value: string
  is_secret: boolean
  created_at: string
}

export interface EnvQuota {
  id: string
  env_id: string
  cpu_limit: string
  memory_limit: string
  pod_limit: number
}

export const deployApi = {
  list: (params?: { project_id?: string; env_id?: string; page?: number; page_size?: number }) =>
    request.get('/deploys', { params }),
  listByEnv: (params?: { env_id?: string }) =>
    request.get('/deploys/by-env', { params }),
  get: (id: string) => request.get(`/deploys/${id}`),
  create: (data: Partial<DeployConfig>) => request.post('/deploys', data),
  update: (id: string, data: Partial<DeployConfig>) => request.put(`/deploys/${id}`, data),
  delete: (id: string) => request.delete(`/deploys/${id}`),
  triggerSync: (id: string) => request.post(`/deploys/${id}/sync`),
  rollback: (id: string, data?: { history_id?: string }) => request.post(`/deploys/${id}/rollback`, data),
  getStatus: (id: string) => request.get(`/deploys/${id}/status`),
  getResources: (id: string) => request.get(`/deploys/${id}/resources`),
  listHistory: (id: string, params?: { page?: number; page_size?: number }) =>
    request.get(`/deploys/${id}/history`, { params }),
  getHistory: (id: string, historyId: string) => request.get(`/deploys/${id}/history/${historyId}`),
  // Rollout (Argo Rollouts)
  getRolloutStatus: (id: string) => request.get(`/deploys/${id}/rollout`),
  promoteRollout: (id: string) => request.post(`/deploys/${id}/rollout/promote`),
  abortRollout: (id: string) => request.post(`/deploys/${id}/rollout/abort`),
  // Approvals
  listPendingApprovals: () => request.get('/approvals/pending'),
  getApproval: (id: string) => request.get(`/approvals/${id}`),
  approve: (id: string, data?: { comment?: string }) => request.post(`/approvals/${id}/approve`, data),
  reject: (id: string, data?: { comment?: string }) => request.post(`/approvals/${id}/reject`, data),
  // Environment variables
  listEnvVars: (envId: string) => request.get(`/environments/${envId}/variables`),
  createEnvVar: (envId: string, data: Partial<EnvVariable>) =>
    request.post(`/environments/${envId}/variables`, data),
  updateEnvVar: (envId: string, varId: string, data: Partial<EnvVariable>) =>
    request.put(`/environments/${envId}/variables/${varId}`, data),
  deleteEnvVar: (envId: string, varId: string) =>
    request.delete(`/environments/${envId}/variables/${varId}`),
  batchUpsertVars: (envId: string, data: { variables: Partial<EnvVariable>[] }) =>
    request.put(`/environments/${envId}/variables/batch`, data),
  getEnvQuota: (envId: string) => request.get(`/environments/${envId}/quota`),
  upsertEnvQuota: (envId: string, data: Partial<EnvQuota>) =>
    request.put(`/environments/${envId}/quota`, data),
}
