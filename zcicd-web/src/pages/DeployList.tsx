import { useState } from 'react'
import { useParams } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Table, Button, Modal, Form, Input, Select, Switch, Tag, Space, App, Popconfirm, Typography, Tooltip,
} from 'antd'
import {
  PlusOutlined, SyncOutlined, RollbackOutlined, HistoryOutlined, EditOutlined, DeleteOutlined,
} from '@ant-design/icons'
import { deployApi, DeployConfig } from '@/api/deploy'

const { Title } = Typography

const SYNC_STATUS_MAP: Record<string, { color: string; label: string }> = {
  Synced: { color: 'green', label: '已同步' },
  OutOfSync: { color: 'orange', label: '未同步' },
  Unknown: { color: 'default', label: '未知' },
}

const HEALTH_STATUS_MAP: Record<string, { color: string; label: string }> = {
  Healthy: { color: 'green', label: '健康' },
  Progressing: { color: 'blue', label: '进行中' },
  Degraded: { color: 'red', label: '降级' },
  Suspended: { color: 'default', label: '暂停' },
  Missing: { color: 'orange', label: '缺失' },
  Unknown: { color: 'default', label: '未知' },
}

const STRATEGY_MAP: Record<string, string> = {
  rolling: '滚动更新',
  'blue-green': '蓝绿部署',
  canary: '金丝雀发布',
}

export default function DeployList() {
  const { projectId } = useParams<{ projectId: string }>()
  const queryClient = useQueryClient()
  const { message } = App.useApp()
  const [form] = Form.useForm()

  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [modalOpen, setModalOpen] = useState(false)
  const [editingDeploy, setEditingDeploy] = useState<DeployConfig | null>(null)

  // --- Data fetching ---
  const { data, isLoading } = useQuery({
    queryKey: ['deploys', projectId, page, pageSize],
    queryFn: async () => {
      const res: any = await deployApi.list({ project_id: projectId, page, page_size: pageSize })
      return res
    },
    enabled: !!projectId,
  })

  const deploys: DeployConfig[] = data?.data ?? []
  const total: number = data?.pagination?.total ?? 0

  // --- Mutations ---
  const createMutation = useMutation({
    mutationFn: (values: Partial<DeployConfig>) => deployApi.create({ ...values, project_id: projectId }),
    onSuccess: () => {
      message.success('部署配置创建成功')
      queryClient.invalidateQueries({ queryKey: ['deploys'] })
      closeModal()
    },
    onError: (err: any) => message.error(err?.message || '创建失败'),
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, values }: { id: string; values: Partial<DeployConfig> }) =>
      deployApi.update(id, values),
    onSuccess: () => {
      message.success('部署配置更新成功')
      queryClient.invalidateQueries({ queryKey: ['deploys'] })
      closeModal()
    },
    onError: (err: any) => message.error(err?.message || '更新失败'),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => deployApi.delete(id),
    onSuccess: () => {
      message.success('部署配置已删除')
      queryClient.invalidateQueries({ queryKey: ['deploys'] })
    },
    onError: (err: any) => message.error(err?.message || '删除失败'),
  })

  const syncMutation = useMutation({
    mutationFn: (id: string) => deployApi.triggerSync(id),
    onSuccess: () => {
      message.success('同步已触发')
      queryClient.invalidateQueries({ queryKey: ['deploys'] })
    },
    onError: (err: any) => message.error(err?.message || '同步失败'),
  })

  const rollbackMutation = useMutation({
    mutationFn: (id: string) => deployApi.rollback(id),
    onSuccess: () => {
      message.success('回滚已触发')
      queryClient.invalidateQueries({ queryKey: ['deploys'] })
    },
    onError: (err: any) => message.error(err?.message || '回滚失败'),
  })

  // --- Modal helpers ---
  const openCreate = () => {
    setEditingDeploy(null)
    form.resetFields()
    form.setFieldsValue({ deploy_strategy: 'rolling', auto_sync: false, require_approval: true })
    setModalOpen(true)
  }

  const openEdit = (deploy: DeployConfig) => {
    setEditingDeploy(deploy)
    form.setFieldsValue({
      service_id: deploy.service_id,
      env_id: deploy.env_id,
      git_repo_url: deploy.git_repo_url,
      git_path: deploy.git_path,
      git_branch: deploy.git_branch,
      deploy_strategy: deploy.deploy_strategy,
      auto_sync: deploy.auto_sync,
      require_approval: deploy.require_approval,
    })
    setModalOpen(true)
  }

  const closeModal = () => {
    setModalOpen(false)
    setEditingDeploy(null)
    form.resetFields()
  }

  const handleSubmit = async () => {
    const values = await form.validateFields()
    if (editingDeploy) {
      updateMutation.mutate({ id: editingDeploy.id, values })
    } else {
      createMutation.mutate(values)
    }
  }

  // --- Table columns ---
  const columns = [
    {
      title: '服务名称',
      dataIndex: 'service_id',
      key: 'service_id',
      width: 150,
    },
    {
      title: '环境',
      dataIndex: 'env_id',
      key: 'env_id',
      width: 120,
    },
    {
      title: '部署策略',
      dataIndex: 'deploy_strategy',
      key: 'deploy_strategy',
      width: 120,
      render: (strategy: string) => STRATEGY_MAP[strategy] || strategy,
    },
    {
      title: '同步状态',
      key: 'sync_status',
      width: 100,
      render: (_: any, record: DeployConfig) => {
        const status = SYNC_STATUS_MAP[record.argocd_app_name] || SYNC_STATUS_MAP.Unknown
        return <Tag color={status.color}>{status.label}</Tag>
      },
    },
    {
      title: '健康状态',
      key: 'health_status',
      width: 100,
      render: (_: any, record: DeployConfig) => {
        const status = HEALTH_STATUS_MAP[record.argocd_app_name] || HEALTH_STATUS_MAP.Unknown
        return <Tag color={status.color}>{status.label}</Tag>
      },
    },
    {
      title: '自动同步',
      dataIndex: 'auto_sync',
      key: 'auto_sync',
      width: 100,
      render: (autoSync: boolean) => (
        <Tag color={autoSync ? 'green' : 'default'}>{autoSync ? '开启' : '关闭'}</Tag>
      ),
    },
    {
      title: '操作',
      key: 'actions',
      width: 240,
      fixed: 'right' as const,
      render: (_: any, record: DeployConfig) => (
        <Space size="small">
          <Tooltip title="同步">
            <Button
              type="link"
              size="small"
              icon={<SyncOutlined />}
              onClick={() => syncMutation.mutate(record.id)}
              loading={syncMutation.isPending}
            />
          </Tooltip>
          <Tooltip title="回滚">
            <Button
              type="link"
              size="small"
              icon={<RollbackOutlined />}
              onClick={() => rollbackMutation.mutate(record.id)}
              loading={rollbackMutation.isPending}
            />
          </Tooltip>
          <Tooltip title="历史记录">
            <Button type="link" size="small" icon={<HistoryOutlined />} />
          </Tooltip>
          <Tooltip title="编辑">
            <Button
              type="link"
              size="small"
              icon={<EditOutlined />}
              onClick={() => openEdit(record)}
            />
          </Tooltip>
          <Popconfirm
            title="确认删除"
            description="确定要删除此部署配置吗？"
            onConfirm={() => deleteMutation.mutate(record.id)}
            okText="删除"
            cancelText="取消"
            okButtonProps={{ danger: true }}
          >
            <Tooltip title="删除">
              <Button
                type="link"
                size="small"
                danger
                icon={<DeleteOutlined />}
              />
            </Tooltip>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div style={{ padding: 24 }}>
      {/* Page header */}
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <Title level={3} style={{ margin: 0 }}>部署管理</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
          新建部署配置
        </Button>
      </div>

      {/* Table */}
      <Table
        columns={columns}
        dataSource={deploys}
        rowKey="id"
        loading={isLoading}
        pagination={{
          current: page,
          pageSize: pageSize,
          total: total,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (t) => `共 ${t} 条记录`,
          onChange: (p, ps) => { setPage(p); setPageSize(ps) },
        }}
        scroll={{ x: 1200 }}
      />

      {/* Create / Edit Modal */}
      <Modal
        title={editingDeploy ? '编辑部署配置' : '新建部署配置'}
        open={modalOpen}
        onCancel={closeModal}
        onOk={handleSubmit}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
        okText={editingDeploy ? '保存' : '创建'}
        cancelText="取消"
        destroyOnClose
        width={600}
      >
        <Form
          form={form}
          layout="vertical"
          style={{ marginTop: 16 }}
          initialValues={{ deploy_strategy: 'rolling', auto_sync: false, require_approval: true }}
        >
          <Form.Item
            name="service_id"
            label="服务ID"
            rules={[{ required: true, message: '请输入服务ID' }]}
          >
            <Input placeholder="请输入服务ID" />
          </Form.Item>

          <Form.Item
            name="env_id"
            label="环境ID"
            rules={[{ required: true, message: '请输入环境ID' }]}
          >
            <Input placeholder="请输入环境ID" />
          </Form.Item>

          <Form.Item
            name="git_repo_url"
            label="Git仓库地址"
            rules={[{ required: true, message: '请输入Git仓库地址' }]}
          >
            <Input placeholder="https://github.com/org/repo" />
          </Form.Item>

          <Form.Item
            name="git_path"
            label="Git路径"
            rules={[{ required: true, message: '请输入Git路径' }]}
          >
            <Input placeholder="例如：manifests/app" />
          </Form.Item>

          <Form.Item
            name="git_branch"
            label="Git分支"
            rules={[{ required: true, message: '请输入Git分支' }]}
          >
            <Input placeholder="例如：main" />
          </Form.Item>

          <Form.Item
            name="deploy_strategy"
            label="部署策略"
            rules={[{ required: true, message: '请选择部署策略' }]}
          >
            <Select
              options={[
                { value: 'rolling', label: '滚动更新' },
                { value: 'blue-green', label: '蓝绿部署' },
                { value: 'canary', label: '金丝雀发布' },
              ]}
            />
          </Form.Item>

          <Form.Item name="auto_sync" label="自动同步" valuePropName="checked">
            <Switch />
          </Form.Item>

          <Form.Item name="require_approval" label="需要审批" valuePropName="checked">
            <Switch />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
