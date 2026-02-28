import { useState, useCallback } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Table, Button, Tag, Modal, Form, Input, Select, Typography,
  Space, Popconfirm, App, Breadcrumb, Tooltip, Skeleton,
} from 'antd'
import type { ColumnsType } from 'antd/es/table'
import {
  PlusOutlined, EditOutlined, DeleteOutlined, ArrowLeftOutlined,
  BranchesOutlined, CodeOutlined, HomeOutlined,
} from '@ant-design/icons'
import { projectApi, Service } from '@/api/project'
import dayjs from 'dayjs'

const { Title } = Typography

const SERVICE_TYPE_MAP: Record<string, { color: string; label: string }> = {
  backend: { color: 'blue', label: '后端' },
  frontend: { color: 'green', label: '前端' },
  middleware: { color: 'orange', label: '中间件' },
}

const STATUS_MAP: Record<string, { color: string; label: string }> = {
  active: { color: 'green', label: '运行中' },
  inactive: { color: 'default', label: '未激活' },
}

const DEPLOY_TYPE_OPTIONS = [
  { value: 'helm', label: 'Helm' },
  { value: 'k8s_yaml', label: 'K8s YAML' },
  { value: 'kustomize', label: 'Kustomize' },
]

const SERVICE_TYPE_OPTIONS = [
  { value: 'backend', label: '后端' },
  { value: 'frontend', label: '前端' },
  { value: 'middleware', label: '中间件' },
]

export default function ServiceList() {
  const { projectId } = useParams<{ projectId: string }>()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const { message: msg } = App.useApp()
  const [form] = Form.useForm()

  const [modalOpen, setModalOpen] = useState(false)
  const [editingService, setEditingService] = useState<Service | null>(null)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)

  // Fetch project info for breadcrumb
  const { data: projectData } = useQuery({
    queryKey: ['project', projectId],
    queryFn: () => projectApi.get(projectId!),
    enabled: !!projectId,
  })

  // Fetch services
  const { data: servicesData, isLoading } = useQuery({
    queryKey: ['services', projectId, page, pageSize],
    queryFn: () => projectApi.listServices(projectId!, { page, page_size: pageSize }),
    enabled: !!projectId,
  })

  const projectName = (projectData as any)?.data?.name ?? '项目'
  const services: Service[] = (servicesData as any)?.data ?? []
  const total: number = (servicesData as any)?.pagination?.total ?? 0

  // Create service
  const createMutation = useMutation({
    mutationFn: (data: Partial<Service>) => projectApi.createService(projectId!, data),
    onSuccess: () => {
      msg.success('服务创建成功')
      queryClient.invalidateQueries({ queryKey: ['services', projectId] })
      closeModal()
    },
    onError: () => msg.error('服务创建失败'),
  })

  // Update service
  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Service> }) =>
      projectApi.updateService(id, data),
    onSuccess: () => {
      msg.success('服务更新成功')
      queryClient.invalidateQueries({ queryKey: ['services', projectId] })
      closeModal()
    },
    onError: () => msg.error('服务更新失败'),
  })

  // Delete service
  const deleteMutation = useMutation({
    mutationFn: (id: string) => projectApi.deleteService(id),
    onSuccess: () => {
      msg.success('服务已删除')
      queryClient.invalidateQueries({ queryKey: ['services', projectId] })
    },
    onError: () => msg.error('删除失败'),
  })

  const openCreate = useCallback(() => {
    setEditingService(null)
    form.resetFields()
    form.setFieldsValue({ branch: 'main', dockerfile_path: 'Dockerfile', build_context: '.', health_check_path: '/healthz' })
    setModalOpen(true)
  }, [form])

  const openEdit = useCallback((record: Service) => {
    setEditingService(record)
    form.setFieldsValue(record)
    setModalOpen(true)
  }, [form])

  const closeModal = useCallback(() => {
    setModalOpen(false)
    setEditingService(null)
    form.resetFields()
  }, [form])

  const handleSubmit = useCallback(async () => {
    const values = await form.validateFields()
    if (editingService) {
      updateMutation.mutate({ id: editingService.id, data: values })
    } else {
      createMutation.mutate(values)
    }
  }, [form, editingService, updateMutation, createMutation])

  const columns: ColumnsType<Service> = [
    {
      title: '服务名称',
      dataIndex: 'name',
      key: 'name',
      render: (text: string, record: Service) => (
        <a
          onClick={() => navigate(`/projects/${projectId}/services/${record.id}`)}
          style={{ fontWeight: 600 }}
        >
          {text}
        </a>
      ),
    },
    {
      title: '类型',
      dataIndex: 'service_type',
      key: 'service_type',
      width: 100,
      render: (val: string) => {
        const cfg = SERVICE_TYPE_MAP[val] ?? { color: 'default', label: val }
        return <Tag color={cfg.color}>{cfg.label}</Tag>
      },
    },
    {
      title: '语言',
      dataIndex: 'language',
      key: 'language',
      width: 100,
      render: (val: string) => val ? <Tag icon={<CodeOutlined />}>{val}</Tag> : '-',
    },
    {
      title: '仓库',
      dataIndex: 'repo_url',
      key: 'repo_url',
      ellipsis: true,
      width: 220,
      render: (val: string) =>
        val ? (
          <Tooltip title={val}>
            <span>{val}</span>
          </Tooltip>
        ) : '-',
    },
    {
      title: '分支',
      dataIndex: 'branch',
      key: 'branch',
      width: 110,
      render: (val: string) =>
        val ? (
          <Tag icon={<BranchesOutlined />} style={{ fontFamily: 'monospace' }}>
            {val}
          </Tag>
        ) : '-',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 90,
      render: (val: string) => {
        const cfg = STATUS_MAP[val] ?? { color: 'default', label: val }
        return <Tag color={cfg.color}>{cfg.label}</Tag>
      },
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
      width: 120,
      render: (_: unknown, record: Service) => (
        <Space size="small">
          <Button type="link" size="small" icon={<EditOutlined />} onClick={() => openEdit(record)}>
            编辑
          </Button>
          <Popconfirm
            title="确认删除"
            description={`确定要删除服务「${record.name}」吗？`}
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
          { title: projectName },
          { title: '服务列表' },
        ]}
      />

      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <Space>
          <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/projects')}>
            返回
          </Button>
          <Title level={4} style={{ margin: 0 }}>服务列表</Title>
        </Space>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
          添加服务
        </Button>
      </div>

      <Table<Service>
        rowKey="id"
        columns={columns}
        dataSource={services}
        pagination={{
          current: page,
          pageSize,
          total,
          showSizeChanger: true,
          showTotal: (t) => `共 ${t} 个服务`,
          onChange: (p, ps) => { setPage(p); setPageSize(ps) },
        }}
        locale={{ emptyText: '暂无服务，点击「添加服务」创建第一个服务' }}
      />

      <Modal
        title={editingService ? '编辑服务' : '添加服务'}
        open={modalOpen}
        onOk={handleSubmit}
        onCancel={closeModal}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
        okText={editingService ? '保存' : '创建'}
        cancelText="取消"
        width={600}
        destroyOnClose
      >
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item name="name" label="服务名称" rules={[{ required: true, message: '请输入服务名称' }]}>
            <Input placeholder="例如：user-service" />
          </Form.Item>
          <Form.Item name="service_type" label="类型" rules={[{ required: true, message: '请选择服务类型' }]}>
            <Select options={SERVICE_TYPE_OPTIONS} placeholder="选择服务类型" />
          </Form.Item>
          <Form.Item name="language" label="语言">
            <Input placeholder="例如：Go, Java, TypeScript" />
          </Form.Item>
          <Form.Item name="repo_url" label="仓库地址">
            <Input placeholder="https://github.com/org/repo" />
          </Form.Item>
          <Form.Item name="branch" label="分支">
            <Input placeholder="main" />
          </Form.Item>
          <Form.Item name="dockerfile_path" label="Dockerfile 路径">
            <Input placeholder="Dockerfile" />
          </Form.Item>
          <Form.Item name="build_context" label="构建上下文">
            <Input placeholder="." />
          </Form.Item>
          <Form.Item name="deploy_type" label="部署方式">
            <Select options={DEPLOY_TYPE_OPTIONS} placeholder="选择部署方式" />
          </Form.Item>
          <Form.Item name="health_check_path" label="健康检查路径">
            <Input placeholder="/healthz" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}