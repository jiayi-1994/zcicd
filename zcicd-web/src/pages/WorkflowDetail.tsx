import { useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Table, Button, Tag, Typography, Space, App, Breadcrumb,
  Skeleton, Descriptions, Card, Switch,
} from 'antd'
import type { ColumnsType } from 'antd/es/table'
import {
  ArrowLeftOutlined, HomeOutlined, PlayCircleOutlined,
} from '@ant-design/icons'
import { workflowApi, Workflow, WorkflowRun } from '@/api/workflow'
import { projectApi } from '@/api/project'
import dayjs from 'dayjs'

const { Title } = Typography

const RUN_STATUS_MAP: Record<string, { color: string; label: string }> = {
  pending: { color: 'default', label: '等待中' },
  running: { color: 'processing', label: '运行中' },
  succeeded: { color: 'success', label: '成功' },
  failed: { color: 'error', label: '失败' },
  cancelled: { color: 'warning', label: '已取消' },
}

const TRIGGER_TYPE_MAP: Record<string, string> = {
  manual: '手动',
  webhook: 'Webhook',
  cron: '定时',
}

export default function WorkflowDetail() {
  const { projectId, workflowId } = useParams<{ projectId: string; workflowId: string }>()
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

  const { data: workflowData, isLoading } = useQuery({
    queryKey: ['workflow', workflowId],
    queryFn: () => workflowApi.getWorkflow(workflowId!),
    enabled: !!workflowId,
  })

  const { data: runsData } = useQuery({
    queryKey: ['workflowRuns', workflowId, page, pageSize],
    queryFn: () => workflowApi.listWorkflowRuns(workflowId!, { page, page_size: pageSize }),
    enabled: !!workflowId,
    refetchInterval: (query) => {
      const runs: WorkflowRun[] = (query.state.data as any)?.data ?? []
      return runs.some(r => r.status === 'running' || r.status === 'pending') ? 5000 : false
    },
  })

  const projectName = (projectData as any)?.data?.name ?? '项目'
  const workflow: Workflow | undefined = (workflowData as any)?.data
  const runs: WorkflowRun[] = (runsData as any)?.data ?? []
  const total: number = (runsData as any)?.pagination?.total ?? 0

  const triggerMutation = useMutation({
    mutationFn: () => workflowApi.triggerWorkflow(workflowId!),
    onSuccess: () => {
      msg.success('工作流已触发')
      queryClient.invalidateQueries({ queryKey: ['workflowRuns', workflowId] })
    },
    onError: () => msg.error('触发失败'),
  })

  const toggleMutation = useMutation({
    mutationFn: (enabled: boolean) => workflowApi.updateWorkflow(workflowId!, { enabled }),
    onSuccess: () => {
      msg.success('状态已更新')
      queryClient.invalidateQueries({ queryKey: ['workflow', workflowId] })
    },
    onError: () => msg.error('更新失败'),
  })

  const columns: ColumnsType<WorkflowRun> = [
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
        const cfg = RUN_STATUS_MAP[val] ?? { color: 'default', label: val }
        return <Tag color={cfg.color}>{cfg.label}</Tag>
      },
    },
    {
      title: '触发方式',
      dataIndex: 'trigger_type',
      key: 'trigger_type',
      width: 100,
      render: (val: string) => TRIGGER_TYPE_MAP[val] ?? val,
    },
    {
      title: '触发者',
      dataIndex: 'triggered_by',
      key: 'triggered_by',
      width: 100,
      render: (val: string) => val || '-',
    },
    {
      title: '开始时间',
      dataIndex: 'started_at',
      key: 'started_at',
      width: 170,
      render: (val: string) => val ? dayjs(val).format('YYYY-MM-DD HH:mm:ss') : '-',
    },
    {
      title: '耗时',
      dataIndex: 'duration_sec',
      key: 'duration_sec',
      width: 80,
      render: (val: number) => val ? `${val}s` : '-',
    },
  ]

  if (isLoading) {
    return (
      <div style={{ padding: 24 }}>
        <Skeleton active paragraph={{ rows: 10 }} />
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
          { title: '工作流', href: `/projects/${projectId}/workflows`, onClick: (e) => { e.preventDefault(); navigate(`/projects/${projectId}/workflows`) } },
          { title: workflow?.name ?? '' },
        ]}
      />

      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <Space>
          <Button icon={<ArrowLeftOutlined />} onClick={() => navigate(`/projects/${projectId}/workflows`)}>
            返回
          </Button>
          <Title level={4} style={{ margin: 0 }}>{workflow?.name}</Title>
        </Space>
        <Button type="primary" icon={<PlayCircleOutlined />} onClick={() => triggerMutation.mutate()} loading={triggerMutation.isPending}>
          触发工作流
        </Button>
      </div>

      <Card style={{ marginBottom: 24 }}>
        <Descriptions column={{ xs: 1, sm: 2, md: 3 }} size="small">
          <Descriptions.Item label="名称">{workflow?.name || '-'}</Descriptions.Item>
          <Descriptions.Item label="描述">{workflow?.description || '-'}</Descriptions.Item>
          <Descriptions.Item label="触发方式">
            <Tag color={workflow?.trigger_type === 'manual' ? 'blue' : workflow?.trigger_type === 'webhook' ? 'green' : 'orange'}>
              {TRIGGER_TYPE_MAP[workflow?.trigger_type ?? ''] ?? workflow?.trigger_type}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label="启用状态">
            <Switch size="small" checked={workflow?.enabled} onChange={(checked) => toggleMutation.mutate(checked)} />
          </Descriptions.Item>
          <Descriptions.Item label="创建时间">{workflow?.created_at ? dayjs(workflow.created_at).format('YYYY-MM-DD HH:mm:ss') : '-'}</Descriptions.Item>
          <Descriptions.Item label="更新时间">{workflow?.updated_at ? dayjs(workflow.updated_at).format('YYYY-MM-DD HH:mm:ss') : '-'}</Descriptions.Item>
        </Descriptions>
      </Card>

      <Title level={5} style={{ marginBottom: 16 }}>运行记录</Title>

      <Table<WorkflowRun>
        rowKey="id"
        columns={columns}
        dataSource={runs}
        pagination={{
          current: page,
          pageSize,
          total,
          showSizeChanger: true,
          showTotal: (t) => `共 ${t} 条运行记录`,
          onChange: (p, ps) => { setPage(p); setPageSize(ps) },
        }}
        locale={{ emptyText: '暂无运行记录' }}
      />
    </div>
  )
}
