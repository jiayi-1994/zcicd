import { useState, useCallback } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Table, Button, Tag, Modal, Form, Input, Select, Typography,
  Space, Popconfirm, App, Breadcrumb, Skeleton, Switch,
} from 'antd'
import type { ColumnsType } from 'antd/es/table'
import {
  PlusOutlined, EditOutlined, DeleteOutlined, ArrowLeftOutlined,
  PlayCircleOutlined, HomeOutlined,
} from '@ant-design/icons'
import { workflowApi, Workflow } from '@/api/workflow'
import { projectApi } from '@/api/project'
import dayjs from 'dayjs'

const { Title } = Typography

const TRIGGER_TYPE_OPTIONS = [
  { value: 'manual', label: '手动触发' },
  { value: 'webhook', label: 'Webhook' },
  { value: 'cron', label: '定时触发' },
]

const TRIGGER_TYPE_MAP: Record<string, { color: string; label: string }> = {
  manual: { color: 'blue', label: '手动' },
  webhook: { color: 'green', label: 'Webhook' },
  cron: { color: 'orange', label: '定时' },
}

export default function WorkflowList() {
  const { projectId } = useParams<{ projectId: string }>()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const { message: msg } = App.useApp()
  const [form] = Form.useForm()

  const [modalOpen, setModalOpen] = useState(false)
  const [editingWorkflow, setEditingWorkflow] = useState<Workflow | null>(null)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)

  const { data: projectData } = useQuery({
    queryKey: ['project', projectId],
    queryFn: () => projectApi.get(projectId!),
    enabled: !!projectId,
  })

  const { data: workflowsData, isLoading } = useQuery({
    queryKey: ['workflows', projectId, page, pageSize],
    queryFn: () => workflowApi.listWorkflows({ project_id: projectId!, page, page_size: pageSize }),
    enabled: !!projectId,
  })

  const projectName = (projectData as any)?.data?.name ?? '项目'
  const workflows: Workflow[] = (workflowsData as any)?.data ?? []
  const total: number = (workflowsData as any)?.pagination?.total ?? 0

  const createMutation = useMutation({
    mutationFn: (data: Partial<Workflow>) => workflowApi.createWorkflow({ ...data, project_id: projectId! }),
    onSuccess: () => {
      msg.success('工作流创建成功')
      queryClient.invalidateQueries({ queryKey: ['workflows', projectId] })
      closeModal()
    },
    onError: () => msg.error('工作流创建失败'),
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Workflow> }) =>
      workflowApi.updateWorkflow(id, data),
    onSuccess: () => {
      msg.success('工作流更新成功')
      queryClient.invalidateQueries({ queryKey: ['workflows', projectId] })
      closeModal()
    },
    onError: () => msg.error('工作流更新失败'),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => workflowApi.deleteWorkflow(id),
    onSuccess: () => {
      msg.success('工作流已删除')
      queryClient.invalidateQueries({ queryKey: ['workflows', projectId] })
    },
    onError: () => msg.error('删除失败'),
  })

  const triggerMutation = useMutation({
    mutationFn: (id: string) => workflowApi.triggerWorkflow(id),
    onSuccess: () => {
      msg.success('工作流已触发')
      queryClient.invalidateQueries({ queryKey: ['workflows', projectId] })
    },
    onError: () => msg.error('触发失败'),
  })

  const toggleMutation = useMutation({
    mutationFn: ({ id, enabled }: { id: string; enabled: boolean }) =>
      workflowApi.updateWorkflow(id, { enabled }),
    onSuccess: () => {
      msg.success('状态已更新')
      queryClient.invalidateQueries({ queryKey: ['workflows', projectId] })
    },
    onError: () => msg.error('更新失败'),
  })

  const openCreate = useCallback(() => {
    setEditingWorkflow(null)
    form.resetFields()
    form.setFieldsValue({ trigger_type: 'manual' })
    setModalOpen(true)
  }, [form])

  const openEdit = useCallback((record: Workflow) => {
    setEditingWorkflow(record)
    form.setFieldsValue(record)
    setModalOpen(true)
  }, [form])

  const closeModal = useCallback(() => {
    setModalOpen(false)
    setEditingWorkflow(null)
    form.resetFields()
  }, [form])

  const handleSubmit = useCallback(async () => {
    const values = await form.validateFields()
    if (editingWorkflow) {
      updateMutation.mutate({ id: editingWorkflow.id, data: values })
    } else {
      createMutation.mutate(values)
    }
  }, [form, editingWorkflow, updateMutation, createMutation])

  const columns: ColumnsType<Workflow> = [
    {
      title: '工作流名称',
      dataIndex: 'name',
      key: 'name',
      render: (text: string, record: Workflow) => (
        <a onClick={() => navigate(`/projects/${projectId}/workflows/${record.id}`)} style={{ fontWeight: 600 }}>
          {text}
        </a>
      ),
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
      width: 200,
      render: (val: string) => val || '-',
    },
    {
      title: '触发方式',
      dataIndex: 'trigger_type',
      key: 'trigger_type',
      width: 110,
      render: (val: string) => {
        const cfg = TRIGGER_TYPE_MAP[val] ?? { color: 'default', label: val }
        return <Tag color={cfg.color}>{cfg.label}</Tag>
      },
    },
    {
      title: '启用',
      dataIndex: 'enabled',
      key: 'enabled',
      width: 80,
      render: (val: boolean, record: Workflow) => (
        <Switch size="small" checked={val} onChange={(checked) => toggleMutation.mutate({ id: record.id, enabled: checked })} />
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 170,
      render: (val: string) => val ? dayjs(val).format('YYYY-MM-DD HH:mm:ss') : '-',
    },
    {
      title: '操作',
      key: 'actions',
      width: 200,
      render: (_: unknown, record: Workflow) => (
        <Space size="small">
          <Button type="link" size="small" icon={<PlayCircleOutlined />} onClick={() => triggerMutation.mutate(record.id)}>
            触发
          </Button>
          <Button type="link" size="small" icon={<EditOutlined />} onClick={() => openEdit(record)}>
            编辑
          </Button>
          <Popconfirm
            title="确认删除"
            description={`确定要删除工作流「${record.name}」吗？`}
            onConfirm={() => deleteMutation.mutate(record.id)}
            okText="删除"
            cancelText="取消"
            okButtonProps={{ danger: true }}
          >
            <Button type="link" size="small" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
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
          { title: '工作流' },
        ]}
      />

      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <Space>
          <Button icon={<ArrowLeftOutlined />} onClick={() => navigate(`/projects/${projectId}`)}>
            返回
          </Button>
          <Title level={4} style={{ margin: 0 }}>工作流</Title>
        </Space>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
          新建工作流
        </Button>
      </div>

      <Table<Workflow>
        rowKey="id"
        columns={columns}
        dataSource={workflows}
        pagination={{
          current: page,
          pageSize,
          total,
          showSizeChanger: true,
          showTotal: (t) => `共 ${t} 个工作流`,
          onChange: (p, ps) => { setPage(p); setPageSize(ps) },
        }}
        locale={{ emptyText: '暂无工作流，点击「新建工作流」创建' }}
      />

      <Modal
        title={editingWorkflow ? '编辑工作流' : '新建工作流'}
        open={modalOpen}
        onOk={handleSubmit}
        onCancel={closeModal}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
        okText={editingWorkflow ? '保存' : '创建'}
        cancelText="取消"
        width={600}
        destroyOnClose
      >
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item name="name" label="工作流名称" rules={[{ required: true, message: '请输入工作流名称' }]}>
            <Input placeholder="例如：deploy-pipeline" />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={3} placeholder="工作流描述（可选）" />
          </Form.Item>
          <Form.Item name="trigger_type" label="触发方式" rules={[{ required: true, message: '请选择触发方式' }]}>
            <Select options={TRIGGER_TYPE_OPTIONS} placeholder="选择触发方式" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
