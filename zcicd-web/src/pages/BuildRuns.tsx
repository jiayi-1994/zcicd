import { useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Table, Button, Tag, Typography, Space, App, Breadcrumb, Skeleton,
} from 'antd'
import type { ColumnsType } from 'antd/es/table'
import {
  ArrowLeftOutlined, HomeOutlined, StopOutlined,
  FileTextOutlined, PlayCircleOutlined,
} from '@ant-design/icons'
import { workflowApi, BuildRun } from '@/api/workflow'
import { projectApi } from '@/api/project'
import dayjs from 'dayjs'

const { Title } = Typography

const BUILD_STATUS_MAP: Record<string, { color: string; label: string }> = {
  pending: { color: 'default', label: '等待中' },
  running: { color: 'processing', label: '运行中' },
  succeeded: { color: 'success', label: '成功' },
  failed: { color: 'error', label: '失败' },
  cancelled: { color: 'warning', label: '已取消' },
}

export default function BuildRuns() {
  const { projectId, configId } = useParams<{ projectId: string; configId: string }>()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const { message: msg } = App.useApp()
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)

  const { data: projectData } = useQuery({
    queryKey: ['project', projectId],
    queryFn: () => projectApi.get(projectId!),
    enabled: !!projectId,
  })

  const { data: configData } = useQuery({
    queryKey: ['buildConfig', configId],
    queryFn: () => workflowApi.getBuildConfig(configId!),
    enabled: !!configId,
  })

  const hasActiveBuilds = (runsData: unknown): boolean => {
    const runs: BuildRun[] = (runsData as any)?.data ?? []
    return runs.some(r => r.status === 'running' || r.status === 'pending')
  }

  const { data: runsData, isLoading } = useQuery({
    queryKey: ['buildRuns', configId, page, pageSize],
    queryFn: () => workflowApi.listBuildRuns({ build_config_id: configId!, page, page_size: pageSize }),
    enabled: !!configId,
    refetchInterval: (query) => hasActiveBuilds(query.state.data) ? 5000 : false,
  })

  const projectName = (projectData as any)?.data?.name ?? '项目'
  const configName = (configData as any)?.data?.name ?? '构建配置'
  const runs: BuildRun[] = (runsData as any)?.data ?? []
  const total: number = (runsData as any)?.pagination?.total ?? 0

  const triggerMutation = useMutation({
    mutationFn: () => workflowApi.triggerBuild(configId!),
    onSuccess: () => {
      msg.success('构建已触发')
      queryClient.invalidateQueries({ queryKey: ['buildRuns', configId] })
    },
    onError: () => msg.error('触发构建失败'),
  })

  const cancelMutation = useMutation({
    mutationFn: (runId: string) => workflowApi.cancelBuildRun(runId),
    onSuccess: () => {
      msg.success('构建已取消')
      queryClient.invalidateQueries({ queryKey: ['buildRuns', configId] })
    },
    onError: () => msg.error('取消构建失败'),
  })

  const columns: ColumnsType<BuildRun> = [
    {
      title: '#',
      dataIndex: 'run_number',
      key: 'run_number',
      width: 70,
      render: (val: number) => <span style={{ fontWeight: 600 }}>#{val}</span>,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (val: string) => {
        const cfg = BUILD_STATUS_MAP[val] ?? { color: 'default', label: val }
        return <Tag color={cfg.color}>{cfg.label}</Tag>
      },
    },
    {
      title: '分支',
      dataIndex: 'branch',
      key: 'branch',
      width: 120,
      render: (val: string) => val || '-',
    },
    {
      title: 'Commit',
      dataIndex: 'commit_sha',
      key: 'commit_sha',
      width: 100,
      render: (val: string) => val ? <code>{val.substring(0, 8)}</code> : '-',
    },
    {
      title: '镜像 Tag',
      dataIndex: 'image_tag',
      key: 'image_tag',
      width: 150,
      render: (val: string) => val ? <Tag style={{ fontFamily: 'monospace' }}>{val}</Tag> : '-',
    },
    {
      title: '触发者',
      dataIndex: 'triggered_by',
      key: 'triggered_by',
      width: 100,
      render: (val: string) => val || '-',
    },
    {
      title: '耗时',
      dataIndex: 'duration_sec',
      key: 'duration_sec',
      width: 80,
      render: (val: number) => val ? `${val}s` : '-',
    },
    {
      title: '开始时间',
      dataIndex: 'started_at',
      key: 'started_at',
      width: 170,
      render: (val: string) => val ? dayjs(val).format('YYYY-MM-DD HH:mm:ss') : '-',
    },
    {
      title: '操作',
      key: 'actions',
      width: 140,
      render: (_: unknown, record: BuildRun) => (
        <Space size="small">
          {(record.status === 'running' || record.status === 'pending') && (
            <Button type="link" size="small" danger icon={<StopOutlined />} onClick={() => cancelMutation.mutate(record.id)}>
              取消
            </Button>
          )}
          <Button type="link" size="small" icon={<FileTextOutlined />}
            onClick={() => navigate(`/projects/${projectId}/builds/${configId}/runs/${record.id}/logs`)}>
            日志
          </Button>
        </Space>
      ),
    },
  ]

  if (isLoading) {
    return (
      <div style={{ padding: 24 }}>
        <Skeleton active paragraph={{ rows: 8 }} />
      </div>
    )
  }

  return (
    <div style={{ padding: 24 }}>
      <Breadcrumb
        style={{ marginBottom: 16 }}
        items={[
          { title: <><HomeOutlined /> 项目</>, href: '/projects', onClick: (e) => { e.preventDefault(); navigate('/projects') } },
          { title: projectName, href: `/projects/${projectId}`, onClick: (e) => { e.preventDefault(); navigate(`/projects/${projectId}`) } },
          { title: '构建管理', href: `/projects/${projectId}/builds`, onClick: (e) => { e.preventDefault(); navigate(`/projects/${projectId}/builds`) } },
          { title: configName },
        ]}
      />

      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <Space>
          <Button icon={<ArrowLeftOutlined />} onClick={() => navigate(`/projects/${projectId}/builds`)}>
            返回
          </Button>
          <Title level={4} style={{ margin: 0 }}>构建记录 - {configName}</Title>
        </Space>
        <Button type="primary" icon={<PlayCircleOutlined />} onClick={() => triggerMutation.mutate()} loading={triggerMutation.isPending}>
          触发构建
        </Button>
      </div>

      <Table<BuildRun>
        rowKey="id"
        columns={columns}
        dataSource={runs}
        pagination={{
          current: page,
          pageSize,
          total,
          showSizeChanger: true,
          showTotal: (t) => `共 ${t} 条构建记录`,
          onChange: (p, ps) => { setPage(p); setPageSize(ps) },
        }}
        locale={{ emptyText: '暂无构建记录' }}
      />
    </div>
  )
}
