import request from '@/utils/request'

export interface TestConfig {
  id: string
  project_id: string
  name: string
  test_type: string
  framework: string
  repo_url: string
  branch: string
  test_command: string
  timeout: number
  enabled: boolean
  created_at: string
  updated_at: string
}

export interface TestRun {
  id: string
  test_config_id: string
  status: string
  total: number
  passed: number
  failed: number
  skipped: number
  coverage: number
  duration_sec: number
  report_url: string
  triggered_by: string
  started_at: string
  finished_at: string
  created_at: string
}

export interface ScanConfig {
  id: string
  project_id: string
  name: string
  scan_type: string
  sonar_project_key: string
  repo_url: string
  branch: string
  enabled: boolean
  created_at: string
  updated_at: string
}

export interface ScanRun {
  id: string
  scan_config_id: string
  status: string
  bugs: number
  vulnerabilities: number
  code_smells: number
  coverage: number
  duplications: number
  report_url: string
  triggered_by: string
  started_at: string
  finished_at: string
  created_at: string
}

export interface QualityGate {
  id: string
  project_id: string
  min_coverage: number
  max_bugs: number
  max_vulnerabilities: number
  max_code_smells: number
  max_duplications: number
  block_deploy: boolean
}

export const qualityApi = {
  // Tests
  listTests: (params?: { project_id?: string; page?: number; page_size?: number }) =>
    request.get(`/projects/${params?.project_id}/tests`, { params: { page: params?.page, page_size: params?.page_size } }),
  getTest: (projectId: string, id: string) => request.get(`/projects/${projectId}/tests/${id}`),
  createTest: (projectId: string, data: Partial<TestConfig>) =>
    request.post(`/projects/${projectId}/tests`, data),
  updateTest: (projectId: string, id: string, data: Partial<TestConfig>) =>
    request.put(`/projects/${projectId}/tests/${id}`, data),
  deleteTest: (projectId: string, id: string) => request.delete(`/projects/${projectId}/tests/${id}`),
  triggerTest: (projectId: string, id: string) =>
    request.post(`/projects/${projectId}/tests/${id}/run`),
  listTestRuns: (projectId: string, id: string, params?: { page?: number; page_size?: number }) =>
    request.get(`/projects/${projectId}/tests/${id}/runs`, { params }),
  getTestRun: (projectId: string, id: string, runId: string) =>
    request.get(`/projects/${projectId}/tests/${id}/runs/${runId}`),

  // Scans
  listScans: (params?: { project_id?: string; page?: number; page_size?: number }) =>
    request.get(`/projects/${params?.project_id}/scans`, { params: { page: params?.page, page_size: params?.page_size } }),
  getScan: (projectId: string, id: string) => request.get(`/projects/${projectId}/scans/${id}`),
  createScan: (projectId: string, data: Partial<ScanConfig>) =>
    request.post(`/projects/${projectId}/scans`, data),
  updateScan: (projectId: string, id: string, data: Partial<ScanConfig>) =>
    request.put(`/projects/${projectId}/scans/${id}`, data),
  deleteScan: (projectId: string, id: string) => request.delete(`/projects/${projectId}/scans/${id}`),
  triggerScan: (projectId: string, id: string) =>
    request.post(`/projects/${projectId}/scans/${id}/run`),
  listScanRuns: (projectId: string, id: string, params?: { page?: number; page_size?: number }) =>
    request.get(`/projects/${projectId}/scans/${id}/runs`, { params }),
  getScanRun: (runId: string) => request.get(`/scans/runs/${runId}`),

  // Quality Gate
  getQualityGate: (projectId: string) => request.get(`/projects/${projectId}/quality-gate`),
  upsertQualityGate: (projectId: string, data: Partial<QualityGate>) =>
    request.put(`/projects/${projectId}/quality-gate`, data),
}
