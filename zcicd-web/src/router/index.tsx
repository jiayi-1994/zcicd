import { lazy, Suspense } from 'react'
import { Navigate, RouteObject } from 'react-router-dom'
import MainLayout from '@/layouts/MainLayout'
import AuthGuard from '@/components/AuthGuard'
import { Spin } from 'antd'

const Login = lazy(() => import('@/pages/Login'))
const ProjectList = lazy(() => import('@/pages/ProjectList'))
const ProjectDetail = lazy(() => import('@/pages/ProjectDetail'))
const ServiceList = lazy(() => import('@/pages/ServiceList'))
const BuildList = lazy(() => import('@/pages/BuildList'))
const BuildRuns = lazy(() => import('@/pages/BuildRuns'))
const BuildLog = lazy(() => import('@/pages/BuildLog'))
const WorkflowList = lazy(() => import('@/pages/WorkflowList'))
const WorkflowDetail = lazy(() => import('@/pages/WorkflowDetail'))
const DeployList = lazy(() => import('@/pages/DeployList'))
const EnvironmentList = lazy(() => import('@/pages/EnvironmentList'))
const TestList = lazy(() => import('@/pages/TestList'))
const ScanList = lazy(() => import('@/pages/ScanList'))
const ArtifactRegistries = lazy(() => import('@/pages/ArtifactRegistries'))
const ClusterList = lazy(() => import('@/pages/ClusterList'))
const IntegrationList = lazy(() => import('@/pages/IntegrationList'))
const AuditLogs = lazy(() => import('@/pages/AuditLogs'))
const NotificationChannels = lazy(() => import('@/pages/NotificationChannels'))
const Dashboard = lazy(() => import('@/pages/Dashboard'))
const NotFound = lazy(() => import('@/pages/NotFound'))

const lazyLoad = (Component: React.LazyExoticComponent<any>) => (
  <Suspense fallback={<Spin size="large" style={{ display: 'flex', justifyContent: 'center', marginTop: '30vh' }} />}>
    <Component />
  </Suspense>
)

export const routes: RouteObject[] = [
  { path: '/login', element: lazyLoad(Login) },
  {
    path: '/',
    element: <AuthGuard><MainLayout /></AuthGuard>,
    children: [
      { index: true, element: <Navigate to="/dashboard" replace /> },
      // Dashboard
      { path: 'dashboard', element: lazyLoad(Dashboard) },
      // Projects
      { path: 'projects', element: lazyLoad(ProjectList) },
      { path: 'projects/:id', element: lazyLoad(ProjectDetail) },
      { path: 'projects/:projectId/services', element: lazyLoad(ServiceList) },
      { path: 'projects/:projectId/builds', element: lazyLoad(BuildList) },
      { path: 'projects/:projectId/builds/:configId/runs', element: lazyLoad(BuildRuns) },
      { path: 'projects/:projectId/builds/:configId/runs/:runId/logs', element: lazyLoad(BuildLog) },
      { path: 'projects/:projectId/workflows', element: lazyLoad(WorkflowList) },
      { path: 'projects/:projectId/workflows/:workflowId', element: lazyLoad(WorkflowDetail) },
      // Deploy & Environments
      { path: 'projects/:projectId/deploys', element: lazyLoad(DeployList) },
      { path: 'projects/:projectId/environments', element: lazyLoad(EnvironmentList) },
      // Quality
      { path: 'projects/:projectId/tests', element: lazyLoad(TestList) },
      { path: 'projects/:projectId/scans', element: lazyLoad(ScanList) },
      // Artifacts
      { path: 'artifacts/registries', element: lazyLoad(ArtifactRegistries) },
      // System
      { path: 'system/clusters', element: lazyLoad(ClusterList) },
      { path: 'system/integrations', element: lazyLoad(IntegrationList) },
      { path: 'system/audit-logs', element: lazyLoad(AuditLogs) },
      { path: 'system/notifications', element: lazyLoad(NotificationChannels) },
    ],
  },
  { path: '*', element: lazyLoad(NotFound) },
]
