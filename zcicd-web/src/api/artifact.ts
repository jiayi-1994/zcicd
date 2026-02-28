import request from '@/utils/request'

export interface ImageRegistry {
  id: string
  name: string
  type: string
  endpoint: string
  username: string
  is_default: boolean
  status: string
  created_at: string
}

export interface ImageScan {
  id: string
  image_name: string
  image_tag: string
  status: string
  critical: number
  high: number
  medium: number
  low: number
  report_url: string
  scanned_at: string
}

export interface HelmChart {
  id: string
  name: string
  version: string
  app_version: string
  description: string
  repo_url: string
  created_at: string
}

export const artifactApi = {
  listRegistries: (params?: { page?: number; page_size?: number }) =>
    request.get('/artifacts/registries', { params }),
  getRegistry: (id: string) => request.get(`/artifacts/registries/${id}`),
  createRegistry: (data: Partial<ImageRegistry>) => request.post('/artifacts/registries', data),
  updateRegistry: (id: string, data: Partial<ImageRegistry>) =>
    request.put(`/artifacts/registries/${id}`, data),
  deleteRegistry: (id: string) => request.delete(`/artifacts/registries/${id}`),
  getScanResults: (name: string) => request.get(`/artifacts/images/${name}/scan`),
  triggerScan: (name: string) => request.post(`/artifacts/images/${name}/scan`),
  listCharts: (params?: { page?: number; page_size?: number }) =>
    request.get('/artifacts/charts', { params }),
  getChart: (name: string) => request.get(`/artifacts/charts/${name}`),
  createChart: (data: Partial<HelmChart>) => request.post('/artifacts/charts', data),
  deleteChart: (name: string) => request.delete(`/artifacts/charts/${name}`),
}
