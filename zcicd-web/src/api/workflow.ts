import request from '@/utils/request'

export interface WorkflowStage {
  name: string
  order: number
  steps: { name: string; type: string; config: Record<string, unknown> }[]
}

export interface Workflow {
  id: string
  project_id: string
  name: string
  description: string
  trigger_type: string
  enabled: boolean
  stages: WorkflowStage[]
  created_at: string
  updated_at: string
}

export interface WorkflowRun {
  id: string
  workflow_id: string
  run_number: number
  status: string
  trigger_type: string
  triggered_by: string
  started_at: string
  finished_at: string
  duration_sec: number
}

export interface BuildConfig {
  id: string
  project_id: string
  service_id: string
  name: string
  template_id: string
  repo_url: string
  branch: string
  dockerfile_path: string
  docker_context: string
  image_repo: string
  tag_strategy: string
  cache_enabled: boolean
  build_script: string
  created_at: string
  updated_at: string
}

export interface BuildRun {
  id: string
  build_config_id: string
  run_number: number
  status: string
  branch: string
  commit_sha: string
  image_tag: string
  tekton_ref: string
  triggered_by: string
  started_at: string
  finished_at: string
  duration_sec: number
}

export interface BuildTemplate {
  id: string
  name: string
  language: string
  framework: string
  description: string
  is_system: boolean
  task_yaml: string
}

export const workflowApi = {
  // Workflows
  listWorkflows: (params?: { project_id?: string; page?: number; page_size?: number }) =>
    request.get('/workflows', { params }),
  getWorkflow: (id: string) => request.get(`/workflows/${id}`),
  createWorkflow: (data: Partial<Workflow>) => request.post('/workflows', data),
  updateWorkflow: (id: string, data: Partial<Workflow>) => request.put(`/workflows/${id}`, data),
  deleteWorkflow: (id: string) => request.delete(`/workflows/${id}`),
  triggerWorkflow: (id: string) => request.post(`/workflows/${id}/trigger`),
  listWorkflowRuns: (id: string, params?: { page?: number; page_size?: number }) =>
    request.get(`/workflows/${id}/runs`, { params }),
  getWorkflowRun: (workflowId: string, runId: string) =>
    request.get(`/workflows/${workflowId}/runs/${runId}`),

  // Build configs
  listBuildConfigs: (params?: { project_id?: string; page?: number; page_size?: number }) =>
    request.get('/build-configs', { params }),
  getBuildConfig: (id: string) => request.get(`/build-configs/${id}`),
  createBuildConfig: (data: Partial<BuildConfig>) => request.post('/build-configs', data),
  updateBuildConfig: (id: string, data: Partial<BuildConfig>) => request.put(`/build-configs/${id}`, data),
  deleteBuildConfig: (id: string) => request.delete(`/build-configs/${id}`),
  triggerBuild: (id: string) => request.post(`/build-configs/${id}/trigger`),

  // Build runs
  listBuildRuns: (params?: { build_config_id?: string; page?: number; page_size?: number }) =>
    request.get('/build-runs', { params }),
  getBuildRun: (runId: string) => request.get(`/build-runs/${runId}`),
  cancelBuildRun: (runId: string) => request.post(`/build-runs/${runId}/cancel`),

  // Build templates
  listBuildTemplates: (params?: { page?: number; page_size?: number }) =>
    request.get('/build-templates', { params }),
  getBuildTemplate: (id: string) => request.get(`/build-templates/${id}`),
  createBuildTemplate: (data: Partial<BuildTemplate>) => request.post('/build-templates', data),
  updateBuildTemplate: (id: string, data: Partial<BuildTemplate>) => request.put(`/build-templates/${id}`, data),
  deleteBuildTemplate: (id: string) => request.delete(`/build-templates/${id}`),
}
