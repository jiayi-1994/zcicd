import { useState, useCallback } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Table, Button, Tag, Modal, Form, Input, Select, Typography,
  Space, Popconfirm, App, Breadcrumb, Tooltip, Skeleton, Switch,
} from 'antd'
import type { ColumnsType } from 'antd/es/table'
import {
  PlusOutlined, EditOutlined, DeleteOutlined, ArrowLeftOutlined,
  PlayCircleOutlined, HomeOutlined, BranchesOutlined,
} from '@ant-design/icons'
import { workflowApi, BuildConfig, BuildTemplate } from '@/api/workflow'
import { projectApi } from '@/api/project'
import dayjs from 'dayjs'

const { Title } = Typography

const TAG_STRATEGY_OPTIONS = [
  { value: 'commit_sha', label: 'Commit SHA' },
  { value: 'branch_latest', label: '分支最新' },
  { value: 'semver', label: '语义化版本' },
  { value: 'datetime', label: '日期时间' },
]

export default function BuildList() {
  const { projectId } = useParams<{ projectId: string }>()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const { message: msg } = App.useApp()
  const [form] = Form.useForm()

  const [modalOpen, setModalOpen] = useState(false)
  const [editingConfig, setEditingConfig] = useState<BuildConfig | null>(null)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)

  const { data: projectData } = useQuery({
    queryKey: ['project', projectId],
    queryFn: () => projectApi.get(projectId!),
    enabled: !!projectId,
  })

  const { data: configsData, isLoading } = useQuery({
    queryKey: ['buildConfigs', projectId, page, pageSize],
    queryFn: () => workflowApi.listBuildConfigs({ project_id: projectId!, page, page_size: pageSize }),
    enabled: !!projectId,
  })

  const { data: templatesData } = useQuery({
    queryKey: ['buildTemplates'],
    queryFn: () => workflowApi.listBuildTemplates(),
  })

  const projectName = (projectData as any)?.data?.name ?? '项目'
  const configs: BuildConfig[] = (configsData as any)?.data ?? []
  const total: number = (configsData as any)?.pagination?.total ?? 0
  const templates: BuildTemplate[] = (templatesData as any)?.data ?? []

  const createMutation = useMutation({
    mutationFn: (data: Partial<BuildConfig>) => workflowApi.createBuildConfig({ ...data, project_id: projectId! }),
    onSuccess: () => {
      msg.success('构建配置创建成功')
      queryClient.invalidateQueries({ queryKey: ['buildConfigs', projectId] })
      closeModal()
    },
    onError: () => msg.error('构建配置创建失败'),
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<BuildConfig> }) =>
      workflowApi.updateBuildConfig(id, data),
    onSuccess: () => {
      msg.success('构建配置更新成功')
      queryClient.invalidateQueries({ queryKey: ['buildConfigs', projectId] })
      closeModal()
    },
    onError: () => msg.error('构建配置更新失败'),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => workflowApi.deleteBuildConfig(id),
    onSuccess: () => {
      msg.success('构建配置已删除')
      queryClient.invalidateQueries({ queryKey: ['buildConfigs', projectId] })
    },
    onError: () => msg.error('删除失败'),
  })

  const triggerMutation = useMutation({
    mutationFn: (id: string) => workflowApi.triggerBuild(id),
    onSuccess: () => {
      msg.success('构建已触发')
      queryClient.invalidateQueries({ queryKey: ['buildConfigs', projectId] })
    },
    onError: () => msg.error('触发构建失败'),
  })

  const openCreate = useCallback(() => {
    setEditingConfig(null)
    form.resetFields()
    form.setFieldsValue({ branch: 'main', dockerfile_path: 'Dockerfile', docker_context: '.', tag_strategy: 'commit_sha', cache_enabled: true })
    setModalOpen(true)
  }, [form])

  const openEdit = useCallback((record: BuildConfig) => {
    setEditingConfig(record)
    form.setFieldsValue(record)
    setModalOpen(true)
  }, [form])

  const closeModal = useCallback(() => {
    setModalOpen(false)
    setEditingConfig(null)
    form.resetFields()
  }, [form])

  const handleSubmit = useCallback(async () => {
    const values = await form.validateFields()
    if (editingConfig) {
      updateMutation.mutate({ id: editingConfig.id, data: values })
    } else {
      createMutation.mutate(values)
    }
  }, [form, editingConfig, updateMutation, createMutation])

  const columns: ColumnsType<BuildConfig> = [
    {
      title: '配置名称',
      dataIndex: 'name',
      key: 'name',
      render: (text: string, record: BuildConfig) => (
        <a onClick={() => navigate(`/projects/${projectId}/builds/${record.id}/runs`)} style={{ fontWeight: 600 }}>
          {text}
        </a>
      ),
    },
    {
      title: '仓库',
      dataIndex: 'repo_url',
      key: 'repo_url',
      ellipsis: true,
      width: 200,
      render: (val: string) => val ? <Tooltip title={val}><span>{val}</span></Tooltip> : '-',
    },
    {
      title: '分支',
      dataIndex: 'branch',
      key: 'branch',
      width: 120,
      render: (val: string) => val ? <Tag icon={<BranchesOutlined />} style={{ fontFamily: 'monospace' }}>{val}</Tag> : '-',
    },
    {
      title: '镜像仓库',
      dataIndex: 'image_repo',
      key: 'image_repo',
      ellipsis: true,
      width: 180,
      render: (val: string) => val || '-',
    },
    {
      title: 'Tag 策略',
      dataIndex: 'tag_strategy',
      key: 'tag_strategy',
      width: 120,
      render: (val: string) => {
        const opt = TAG_STRATEGY_OPTIONS.find(o => o.value === val)
        return <Tag>{opt?.label ?? val}</Tag>
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
      width: 180,
      render: (_: unknown, record: BuildConfig) => (
        <Space size="small">
          <Button type="link" size="small" icon={<PlayCircleOutlined />} onClick={() => triggerMutation.mutate(record.id)}>
            构建
          </Button>
          <Button type="link" size="small" icon={<EditOutlined />} onClick={() => openEdit(record)}>
            编辑
          </Button>
          <Popconfirm
            title="确认删除"
            description={`确定要删除构建配置「${record.name}」吗？`}
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
          { title: '构建管理' },
        ]}
      />

      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <Space>
          <Button icon={<ArrowLeftOutlined />} onClick={() => navigate(`/projects/${projectId}`)}>
            返回
          </Button>
          <Title level={4} style={{ margin: 0 }}>构建管理</Title>
        </Space>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
          新建构建配置
        </Button>
      </div>

      <Table<BuildConfig>
        rowKey="id"
        columns={columns}
        dataSource={configs}
        pagination={{
          current: page,
          pageSize,
          total,
          showSizeChanger: true,
          showTotal: (t) => `共 ${t} 个构建配置`,
          onChange: (p, ps) => { setPage(p); setPageSize(ps) },
        }}
        locale={{ emptyText: '暂无构建配置，点击「新建构建配置」创建' }}
      />

      <Modal
        title={editingConfig ? '编辑构建配置' : '新建构建配置'}
        open={modalOpen}
        onOk={handleSubmit}
        onCancel={closeModal}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
        okText={editingConfig ? '保存' : '创建'}
        cancelText="取消"
        width={600}
        destroyOnClose
      >
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item name="name" label="配置名称" rules={[{ required: true, message: '请输入配置名称' }]}>
            <Input placeholder="例如：user-service-build" />
          </Form.Item>
          <Form.Item name="template_id" label="构建模板">
            <Select
              placeholder="选择构建模板"
              allowClear
              options={templates.map(t => ({ value: t.id, label: `${t.name} (${t.language}/${t.framework})` }))}
            />
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
          <Form.Item name="docker_context" label="Docker 上下文">
            <Input placeholder="." />
          </Form.Item>
          <Form.Item name="image_repo" label="镜像仓库">
            <Input placeholder="registry.example.com/project/service" />
          </Form.Item>
          <Form.Item name="tag_strategy" label="Tag 策略">
            <Select options={TAG_STRATEGY_OPTIONS} placeholder="选择 Tag 策略" />
          </Form.Item>
          <Form.Item name="cache_enabled" label="启用缓存" valuePropName="checked">
            <Switch />
          </Form.Item>
          <Form.Item name="build_script" label="构建脚本">
            <Input.TextArea rows={3} placeholder="自定义构建脚本（可选）" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
