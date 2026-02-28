import request from '@/utils/request'

export interface Cluster {
  id: string
  name: string
  display_name: string
  description: string
  provider: string
  api_server_url: string
  status: string
  node_count: number
  version: string
  created_at: string
  updated_at: string
}

export interface Integration {
  id: string
  name: string
  type: string
  provider: string
  status: string
  last_check_at: string
  created_at: string
}

export interface AuditLog {
  id: string
  user_id: string
  username: string
  action: string
  resource_type: string
  resource_id: string
  resource_name: string
  project_id: string
  detail: Record<string, unknown>
  ip_address: string
  created_at: string
}

export interface NotifyChannel {
  id: string
  name: string
  type: string
  config: Record<string, unknown>
  enabled: boolean
  created_at: string
}

export interface NotifyRule {
  id: string
  name: string
  event_type: string
  severity: string
  project_id: string
  channel_id: string
  enabled: boolean
  created_at: string
}

export interface OverviewStats {
  projects: number
  services: number
  clusters: number
  environments: number
  registries: number
  integrations: number
}

export interface DailyStat {
  id: string
  stat_date: string
  total_builds: number
  successful_builds: number
  failed_builds: number
  total_deploys: number
  successful_deploys: number
  failed_deploys: number
}

export const systemApi = {
  // Dashboard
  getDashboardOverview: () => request.get('/system/dashboard/overview'),
  getDashboardTrends: (days?: number) =>
    request.get('/system/dashboard/trends', { params: days ? { days } : undefined }),

  // Clusters
  listClusters: (params?: { page?: number; page_size?: number }) =>
    request.get('/system/clusters', { params }),
  getCluster: (id: string) => request.get(`/system/clusters/${id}`),
  createCluster: (data: Partial<Cluster>) => request.post('/system/clusters', data),
  updateCluster: (id: string, data: Partial<Cluster>) =>
    request.put(`/system/clusters/${id}`, data),
  deleteCluster: (id: string) => request.delete(`/system/clusters/${id}`),

  // Integrations
  listIntegrations: (params?: { page?: number; page_size?: number }) =>
    request.get('/system/integrations', { params }),
  createIntegration: (data: Partial<Integration>) =>
    request.post('/system/integrations', data),
  updateIntegration: (id: string, data: Partial<Integration>) =>
    request.put(`/system/integrations/${id}`, data),
  deleteIntegration: (id: string) => request.delete(`/system/integrations/${id}`),

  // Audit logs
  listAuditLogs: (params?: {
    user_id?: string; action?: string; resource_type?: string;
    project_id?: string; page?: number; page_size?: number
  }) => request.get('/system/audit-logs', { params }),

  // Notifications
  listChannels: (params?: { page?: number; page_size?: number }) =>
    request.get('/notifications/channels', { params }),
  createChannel: (data: Partial<NotifyChannel>) =>
    request.post('/notifications/channels', data),
  updateChannel: (id: string, data: Partial<NotifyChannel>) =>
    request.put(`/notifications/channels/${id}`, data),
  deleteChannel: (id: string) => request.delete(`/notifications/channels/${id}`),
  listRules: (params?: { page?: number; page_size?: number }) =>
    request.get('/notifications/rules', { params }),
  createRule: (data: Partial<NotifyRule>) => request.post('/notifications/rules', data),
  updateRule: (id: string, data: Partial<NotifyRule>) =>
    request.put(`/notifications/rules/${id}`, data),
}
