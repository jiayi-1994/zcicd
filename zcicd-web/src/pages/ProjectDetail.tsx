import { useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Card, Tabs, Table, Button, Tag, Descriptions, Space, Modal, Form,
  Input, Select, Typography, Popconfirm, App, Spin, Result, Switch, Row, Col, Divider,
} from 'antd'
import {
  ArrowLeftOutlined, PlusOutlined, EditOutlined, DeleteOutlined,
  SettingOutlined, CloudServerOutlined, AppstoreOutlined,
  BuildOutlined, BranchesOutlined,
} from '@ant-design/icons'
import { projectApi, Project, Service, Environment } from '@/api/project'
import dayjs from 'dayjs'

const { Title, Text } = Typography

const statusColors: Record<string, string> = {
  active: 'green',
  inactive: 'default',
  archived: 'red',
  running: 'green',
  stopped: 'orange',
  error: 'red',
}

const envTypeColors: Record<string, string> = {
  dev: 'blue',
  testing: 'orange',
  staging: 'purple',
  production: 'red',
}

const envTypeLabels: Record<string, string> = {
  dev: '开发',
  testing: '测试',
  staging: '预发布',
  production: '生产',
}

// --- Service Modal ---
function ServiceModal({
  open, editingService, projectId, onClose,
}: {
  open: boolean
  editingService: Service | null
  projectId: string
  onClose: () => void
}) {
  const [form] = Form.useForm()
  const { message } = App.useApp()
  const queryClient = useQueryClient()
  const createMut = useMutation({
    mutationFn: (data: Partial<Service>) => projectApi.createService(projectId, data),
    onSuccess: () => { message.success('服务创建成功'); queryClient.invalidateQueries({ queryKey: ['services', projectId] }); onClose() },
    onError: (e: any) => message.error(e?.message || '操作失败'),
  })
  const updateMut = useMutation({
    mutationFn: (data: Partial<Service>) => projectApi.updateService(editingService!.id, data),
    onSuccess: () => { message.success('服务更新成功'); queryClient.invalidateQueries({ queryKey: ['services', projectId] }); onClose() },
    onError: (e: any) => message.error(e?.message || '操作失败'),
  })

  const handleOk = async () => {
    const values = await form.validateFields()
    editingService ? updateMut.mutate(values) : createMut.mutate(values)
  }

  return (
    <Modal
      title={editingService ? '编辑服务' : '添加服务'}
      open={open}
      onCancel={onClose}
      onOk={handleOk}
      confirmLoading={createMut.isPending || updateMut.isPending}
      destroyOnClose
    >
      <Form form={form} layout="vertical" initialValues={editingService || { service_type: 'backend', language: 'go', branch: 'main' }}>
        <Form.Item name="name" label="服务名称" rules={[{ required: true, message: '请输入服务名称' }]}>
          <Input placeholder="my-service" />
        </Form.Item>
        <Form.Item name="service_type" label="类型" rules={[{ required: true }]}>
          <Select options={[
            { label: '后端服务', value: 'backend' },
            { label: '前端服务', value: 'frontend' },
            { label: '中间件', value: 'middleware' },
            { label: '任务', value: 'job' },
          ]} />
        </Form.Item>
        <Form.Item name="language" label="语言" rules={[{ required: true }]}>
          <Select options={[
            { label: 'Go', value: 'go' },
            { label: 'Java', value: 'java' },
            { label: 'Python', value: 'python' },
            { label: 'Node.js', value: 'nodejs' },
            { label: 'TypeScript', value: 'typescript' },
          ]} />
        </Form.Item>
        <Form.Item name="repo_url" label="仓库地址">
          <Input placeholder="https://github.com/org/repo" />
        </Form.Item>
        <Form.Item name="branch" label="分支" rules={[{ required: true }]}>
          <Input placeholder="main" />
        </Form.Item>
      </Form>
    </Modal>
  )
}

// --- Environment Modal ---
function EnvModal({
  open, editingEnv, projectId, onClose,
}: {
  open: boolean
  editingEnv: Environment | null
  projectId: string
  onClose: () => void
}) {
  const [form] = Form.useForm()
  const { message } = App.useApp()
  const queryClient = useQueryClient()

  const createMut = useMutation({
    mutationFn: (data: Partial<Environment>) => projectApi.createEnvironment(projectId, data),
    onSuccess: () => { message.success('环境创建成功'); queryClient.invalidateQueries({ queryKey: ['environments', projectId] }); onClose() },
    onError: (e: any) => message.error(e?.message || '操作失败'),
  })
  const updateMut = useMutation({
    mutationFn: (data: Partial<Environment>) => projectApi.updateEnvironment(editingEnv!.id, data),
    onSuccess: () => { message.success('环境更新成功'); queryClient.invalidateQueries({ queryKey: ['environments', projectId] }); onClose() },
    onError: (e: any) => message.error(e?.message || '操作失败'),
  })

  const handleOk = async () => {
    const values = await form.validateFields()
    editingEnv ? updateMut.mutate(values) : createMut.mutate(values)
  }

  return (
    <Modal
      title={editingEnv ? '编辑环境' : '添加环境'}
      open={open}
      onCancel={onClose}
      onOk={handleOk}
      confirmLoading={createMut.isPending || updateMut.isPending}
      destroyOnClose
    >
      <Form form={form} layout="vertical" initialValues={editingEnv || { env_type: 'dev', auto_deploy: false }}>
        <Form.Item name="name" label="环境名称" rules={[{ required: true, message: '请输入环境名称' }]}>
          <Input placeholder="dev-01" />
        </Form.Item>
        <Form.Item name="env_type" label="环境类型" rules={[{ required: true }]}>
          <Select options={[
            { label: '开发', value: 'dev' },
            { label: '测试', value: 'testing' },
            { label: '预发布', value: 'staging' },
            { label: '生产', value: 'production' },
          ]} />
        </Form.Item>
        <Form.Item name="namespace" label="命名空间" rules={[{ required: true, message: '请输入命名空间' }]}>
          <Input placeholder="project-dev" />
        </Form.Item>
        <Form.Item name="cluster" label="集群" rules={[{ required: true, message: '请输入集群名称' }]}>
          <Input placeholder="k8s-cluster-01" />
        </Form.Item>
        <Form.Item name="auto_deploy" label="自动部署" valuePropName="checked">
          <Switch />
        </Form.Item>
      </Form>
    </Modal>
  )
}

// --- Main Page ---
export default function ProjectDetail() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { message } = App.useApp()
  const queryClient = useQueryClient()
  const [activeTab, setActiveTab] = useState('services')

  // Service modal state
  const [svcModalOpen, setSvcModalOpen] = useState(false)
  const [editingSvc, setEditingSvc] = useState<Service | null>(null)

  // Environment modal state
  const [envModalOpen, setEnvModalOpen] = useState(false)
  const [editingEnv, setEditingEnv] = useState<Environment | null>(null)

  // Settings form
  const [settingsForm] = Form.useForm()

  // --- Queries ---
  const { data: projectRes, isLoading, isError } = useQuery({
    queryKey: ['project', id],
    queryFn: () => projectApi.get(id!),
    enabled: !!id,
  })
  const project: Project | undefined = (projectRes as any)?.data

  const { data: servicesRes } = useQuery({
    queryKey: ['services', id],
    queryFn: () => projectApi.listServices(id!),
    enabled: !!id && activeTab === 'services',
  })
  const services: Service[] = (servicesRes as any)?.data || []

  const { data: envsRes } = useQuery({
    queryKey: ['environments', id],
    queryFn: () => projectApi.listEnvironments(id!),
    enabled: !!id && activeTab === 'environments',
  })
  const environments: Environment[] = (envsRes as any)?.data || []

  // --- Mutations ---
  const deleteSvcMut = useMutation({
    mutationFn: (svcId: string) => projectApi.deleteService(svcId),
    onSuccess: () => { message.success('服务已删除'); queryClient.invalidateQueries({ queryKey: ['services', id] }) },
    onError: (e: any) => message.error(e?.message || '删除失败'),
  })

  const deleteEnvMut = useMutation({
    mutationFn: (envId: string) => projectApi.deleteEnvironment(envId),
    onSuccess: () => { message.success('环境已删除'); queryClient.invalidateQueries({ queryKey: ['environments', id] }) },
    onError: (e: any) => message.error(e?.message || '删除失败'),
  })

  const updateProjectMut = useMutation({
    mutationFn: (data: Partial<Project>) => projectApi.update(id!, data),
    onSuccess: () => { message.success('项目已更新'); queryClient.invalidateQueries({ queryKey: ['project', id] }) },
    onError: (e: any) => message.error(e?.message || '更新失败'),
  })

  const deleteProjectMut = useMutation({
    mutationFn: () => projectApi.delete(id!),
    onSuccess: () => { message.success('项目已删除'); navigate('/projects') },
    onError: (e: any) => message.error(e?.message || '删除失败'),
  })

  // --- Loading / Error ---
  if (isLoading) return <Spin size="large" style={{ display: 'block', margin: '120px auto' }} />
  if (isError || !project) return <Result status="404" title="项目不存在" subTitle="请检查项目 ID 是否正确" extra={<Button onClick={() => navigate('/projects')}>返回项目列表</Button>} />

  // --- Services Table Columns ---
  const serviceColumns = [
    { title: '服务名称', dataIndex: 'name', key: 'name' },
    {
      title: '类型', dataIndex: 'service_type', key: 'service_type',
      render: (v: string) => <Tag>{v}</Tag>,
    },
    { title: '语言', dataIndex: 'language', key: 'language' },
    { title: '分支', dataIndex: 'branch', key: 'branch' },
    {
      title: '状态', dataIndex: 'status', key: 'status',
      render: (v: string) => <Tag color={statusColors[v] || 'default'}>{v}</Tag>,
    },
    {
      title: '操作', key: 'actions', width: 120,
      render: (_: unknown, record: Service) => (
        <Space size="small">
          <Button type="link" size="small" icon={<EditOutlined />} onClick={() => { setEditingSvc(record); setSvcModalOpen(true) }} />
          <Popconfirm title="确定删除该服务？" onConfirm={() => deleteSvcMut.mutate(record.id)} okText="删除" cancelText="取消">
            <Button type="link" size="small" danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ]

  // --- Tab Content ---
  const servicesTab = (
    <>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'flex-end' }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => { setEditingSvc(null); setSvcModalOpen(true) }}>
          添加服务
        </Button>
      </div>
      <Table rowKey="id" columns={serviceColumns} dataSource={services} pagination={false} />
    </>
  )

  const environmentsTab = (
    <>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'flex-end' }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => { setEditingEnv(null); setEnvModalOpen(true) }}>
          添加环境
        </Button>
      </div>
      <Row gutter={[16, 16]}>
        {environments.map((env) => (
          <Col key={env.id} xs={24} sm={12} lg={8}>
            <Card
              size="small"
              title={env.name}
              extra={
                <Space size="small">
                  <Button type="link" size="small" icon={<EditOutlined />} onClick={() => { setEditingEnv(env); setEnvModalOpen(true) }} />
                  <Popconfirm title="确定删除该环境？" onConfirm={() => deleteEnvMut.mutate(env.id)} okText="删除" cancelText="取消">
                    <Button type="link" size="small" danger icon={<DeleteOutlined />} />
                  </Popconfirm>
                </Space>
              }
            >
              <p><Text type="secondary">类型：</Text><Tag color={envTypeColors[env.env_type]}>{envTypeLabels[env.env_type] || env.env_type}</Tag></p>
              <p><Text type="secondary">命名空间：</Text>{env.namespace}</p>
              <p><Text type="secondary">集群：</Text>{env.cluster}</p>
              <p><Text type="secondary">自动部署：</Text><Switch size="small" checked={env.auto_deploy} disabled /></p>
            </Card>
          </Col>
        ))}
        {environments.length === 0 && (
          <Col span={24}><Result status="info" title="暂无环境" subTitle="点击上方按钮添加环境" /></Col>
        )}
      </Row>
    </>
  )

  const settingsTab = (
    <div style={{ maxWidth: 600 }}>
      <Form
        form={settingsForm}
        layout="vertical"
        initialValues={project}
        onFinish={(values) => updateProjectMut.mutate(values)}
      >
        <Form.Item name="name" label="项目名称" rules={[{ required: true, message: '请输入项目名称' }]}>
          <Input />
        </Form.Item>
        <Form.Item name="description" label="项目描述">
          <Input.TextArea rows={3} />
        </Form.Item>
        <Form.Item name="repo_url" label="仓库地址">
          <Input />
        </Form.Item>
        <Form.Item name="default_branch" label="默认分支">
          <Input />
        </Form.Item>
        <Form.Item name="visibility" label="可见性">
          <Select options={[
            { label: '公开', value: 'public' },
            { label: '私有', value: 'private' },
          ]} />
        </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit" loading={updateProjectMut.isPending}>
            保存设置
          </Button>
        </Form.Item>
      </Form>

      <Divider />

      <div style={{ padding: '16px', border: '1px solid #ff4d4f', borderRadius: 8 }}>
        <Title level={5} type="danger">危险操作</Title>
        <Text type="secondary">删除项目后，所有关联的服务和环境将被永久移除，此操作不可撤销。</Text>
        <div style={{ marginTop: 12 }}>
          <Popconfirm
            title="确定要删除此项目？"
            description="此操作不可撤销，所有关联数据将被永久删除。"
            onConfirm={() => deleteProjectMut.mutate()}
            okText="确认删除"
            cancelText="取消"
            okButtonProps={{ danger: true }}
          >
            <Button danger loading={deleteProjectMut.isPending}>删除项目</Button>
          </Popconfirm>
        </div>
      </div>
    </div>
  )

  return (
    <div>
      {/* Header */}
      <div style={{ marginBottom: 24 }}>
        <Space align="center" style={{ marginBottom: 16 }}>
          <Button type="text" icon={<ArrowLeftOutlined />} onClick={() => navigate('/projects')}>
            返回
          </Button>
        </Space>
        <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
          <Title level={3} style={{ margin: 0 }}>{project.name}</Title>
          <Tag>{project.identifier}</Tag>
          <Tag color={statusColors[project.status] || 'default'}>{project.status}</Tag>
        </div>
      </div>

      {/* Metadata */}
      <Card style={{ marginBottom: 24 }}>
        <Descriptions column={{ xs: 1, sm: 2, md: 3 }} size="small">
          <Descriptions.Item label="项目标识">{project.identifier}</Descriptions.Item>
          <Descriptions.Item label="仓库地址">{project.repo_url || '-'}</Descriptions.Item>
          <Descriptions.Item label="默认分支">{project.default_branch || '-'}</Descriptions.Item>
          <Descriptions.Item label="可见性">
            <Tag color={project.visibility === 'public' ? 'green' : 'orange'}>
              {project.visibility === 'public' ? '公开' : '私有'}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label="创建时间">{dayjs(project.created_at).format('YYYY-MM-DD HH:mm')}</Descriptions.Item>
          <Descriptions.Item label="更新时间">{dayjs(project.updated_at).format('YYYY-MM-DD HH:mm')}</Descriptions.Item>
        </Descriptions>
        <div style={{ marginTop: 16 }}>
          <Space>
            <Button icon={<BuildOutlined />} onClick={() => navigate(`/projects/${id}/builds`)}>
              构建管理
            </Button>
            <Button icon={<BranchesOutlined />} onClick={() => navigate(`/projects/${id}/workflows`)}>
              工作流
            </Button>
          </Space>
        </div>
      </Card>

      {/* Tabs */}
      <Card>
        <Tabs
          activeKey={activeTab}
          onChange={setActiveTab}
          items={[
            {
              key: 'services',
              label: <span><AppstoreOutlined /> 服务列表</span>,
              children: servicesTab,
            },
            {
              key: 'environments',
              label: <span><CloudServerOutlined /> 环境管理</span>,
              children: environmentsTab,
            },
            {
              key: 'settings',
              label: <span><SettingOutlined /> 项目设置</span>,
              children: settingsTab,
            },
          ]}
        />
      </Card>

      {/* Modals */}
      <ServiceModal open={svcModalOpen} editingService={editingSvc} projectId={id!} onClose={() => { setSvcModalOpen(false); setEditingSvc(null) }} />
      <EnvModal open={envModalOpen} editingEnv={editingEnv} projectId={id!} onClose={() => { setEnvModalOpen(false); setEditingEnv(null) }} />
    </div>
  )
}
