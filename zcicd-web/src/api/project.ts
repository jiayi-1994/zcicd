import request from '@/utils/request'

export interface Project {
  id: string
  name: string
  identifier: string
  description: string
  owner_id: string
  repo_url: string
  default_branch: string
  visibility: string
  status: string
  created_at: string
  updated_at: string
}

export interface Service {
  id: string
  project_id: string
  name: string
  service_type: string
  language: string
  repo_url: string
  branch: string
  status: string
  created_at: string
}

export interface Environment {
  id: string
  project_id: string
  name: string
  env_type: string
  namespace: string
  cluster: string
  auto_deploy: boolean
  created_at: string
}

export const projectApi = {
  list: (params?: { page?: number; page_size?: number; keyword?: string }) =>
    request.get('/projects', { params }),
  get: (id: string) => request.get(`/projects/${id}`),
  create: (data: Partial<Project>) => request.post('/projects', data),
  update: (id: string, data: Partial<Project>) => request.put(`/projects/${id}`, data),
  delete: (id: string) => request.delete(`/projects/${id}`),
  listServices: (projectId: string, params?: { page?: number; page_size?: number }) =>
    request.get(`/projects/${projectId}/services`, { params }),
  createService: (projectId: string, data: Partial<Service>) =>
    request.post(`/projects/${projectId}/services`, data),
  getService: (id: string) => request.get(`/services/${id}`),
  updateService: (id: string, data: Partial<Service>) => request.put(`/services/${id}`, data),
  deleteService: (id: string) => request.delete(`/services/${id}`),
  listEnvironments: (projectId: string) =>
    request.get(`/projects/${projectId}/environments`),
  createEnvironment: (projectId: string, data: Partial<Environment>) =>
    request.post(`/projects/${projectId}/environments`, data),
  updateEnvironment: (id: string, data: Partial<Environment>) =>
    request.put(`/environments/${id}`, data),
  deleteEnvironment: (id: string) => request.delete(`/environments/${id}`),
}
