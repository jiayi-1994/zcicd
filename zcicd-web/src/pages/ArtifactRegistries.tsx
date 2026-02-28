import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Table, Button, Input, Tag, Modal, Form, Select, Space, App, Popconfirm, Typography,
} from 'antd'
import {
  PlusOutlined, SearchOutlined, EditOutlined, DeleteOutlined,
} from '@ant-design/icons'
import { artifactApi, ImageRegistry } from '@/api/artifact'

const { Title } = Typography

const REGISTRY_TYPE_MAP: Record<string, { color: string; label: string }> = {
  harbor: { color: 'blue', label: 'Harbor' },
  dockerhub: { color: 'cyan', label: 'Docker Hub' },
  acr: { color: 'orange', label: 'ACR' },
  ecr: { color: 'green', label: 'ECR' },
  ghcr: { color: 'purple', label: 'GHCR' },
}

export default function ArtifactRegistries() {
  const queryClient = useQueryClient()
  const { message } = App.useApp()
  const [form] = Form.useForm()

  const [keyword, setKeyword] = useState('')
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [modalOpen, setModalOpen] = useState(false)
  const [editingRegistry, setEditingRegistry] = useState<ImageRegistry | null>(null)

  // --- Data fetching ---
  const { data, isLoading } = useQuery({
    queryKey: ['registries', page, pageSize, keyword],
    queryFn: async () => {
      const res: any = await artifactApi.listRegistries({ page, page_size: pageSize })
      return res
    },
  })

  const registries: ImageRegistry[] = data?.data ?? []
  const total: number = data?.pagination?.total ?? 0

  // --- Mutations ---
  const createMutation = useMutation({
    mutationFn: (values: Partial<ImageRegistry>) => artifactApi.createRegistry(values),
    onSuccess: () => {
      message.success('镜像仓库创建成功')
      queryClient.invalidateQueries({ queryKey: ['registries'] })
      closeModal()
    },
    onError: (err: any) => message.error(err?.message || '创建失败'),
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, values }: { id: string; values: Partial<ImageRegistry> }) =>
      artifactApi.updateRegistry(id, values),
    onSuccess: () => {
      message.success('镜像仓库更新成功')
      queryClient.invalidateQueries({ queryKey: ['registries'] })
      closeModal()
    },
    onError: (err: any) => message.error(err?.message || '更新失败'),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => artifactApi.deleteRegistry(id),
    onSuccess: () => {
      message.success('镜像仓库已删除')
      queryClient.invalidateQueries({ queryKey: ['registries'] })
    },
    onError: (err: any) => message.error(err?.message || '删除失败'),
  })

  // --- Modal helpers ---
  const openCreate = () => {
    setEditingRegistry(null)
    form.resetFields()
    form.setFieldsValue({ type: 'harbor' })
    setModalOpen(true)
  }

  const openEdit = (registry: ImageRegistry) => {
    setEditingRegistry(registry)
    form.setFieldsValue({
      name: registry.name,
      type: registry.type,
      endpoint: registry.endpoint,
      username: registry.username,
    })
    setModalOpen(true)
  }

  const closeModal = () => {
    setModalOpen(false)
    setEditingRegistry(null)
    form.resetFields()
  }

  const handleSubmit = async () => {
    const values = await form.validateFields()
    if (editingRegistry) {
      updateMutation.mutate({ id: editingRegistry.id, values })
    } else {
      createMutation.mutate(values)
    }
  }

  const handleSearch = (value: string) => {
    setKeyword(value)
    setPage(1)
  }

  const getRegistryType = (type: string) => REGISTRY_TYPE_MAP[type] || REGISTRY_TYPE_MAP.harbor

  // --- Table columns ---
  const columns = [
    {
      title: '名称',
      dataIndex: 'name',
      key: 'name',
      width: 200,
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      width: 120,
      render: (type: string) => {
        const t = getRegistryType(type)
        return <Tag color={t.color}>{t.label}</Tag>
      },
    },
    {
      title: '端点地址',
      dataIndex: 'endpoint',
      key: 'endpoint',
      ellipsis: true,
    },
    {
      title: '默认仓库',
      dataIndex: 'is_default',
      key: 'is_default',
      width: 100,
      render: (isDefault: boolean) => (
        isDefault ? <Tag color="gold">默认</Tag> : null
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: string) => (
        <Tag color={status === 'active' ? 'green' : 'default'}>
          {status === 'active' ? '正常' : '异常'}
        </Tag>
      ),
    },
    {
      title: '操作',
      key: 'actions',
      width: 150,
      render: (_: any, record: ImageRegistry) => (
        <Space size="small">
          <Button
            type="link"
            size="small"
            icon={<EditOutlined />}
            onClick={() => openEdit(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确认删除"
            description={`确定要删除镜像仓库「${record.name}」吗？`}
            onConfirm={() => deleteMutation.mutate(record.id)}
            okText="删除"
            cancelText="取消"
            okButtonProps={{ danger: true }}
          >
            <Button
              type="link"
              size="small"
              danger
              icon={<DeleteOutlined />}
            >
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div style={{ padding: 24 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <Title level={3} style={{ margin: 0 }}>镜像仓库管理</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
          新建镜像仓库
        </Button>
      </div>

      <div style={{ marginBottom: 16 }}>
        <Input.Search
          placeholder="搜索镜像仓库..."
          allowClear
          enterButton={<><SearchOutlined /> 搜索</>}
          onSearch={handleSearch}
          style={{ maxWidth: 400 }}
        />
      </div>

      <Table
        columns={columns}
        dataSource={registries}
        rowKey="id"
        loading={isLoading}
        pagination={{
          current: page,
          pageSize: pageSize,
          total: total,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (t) => `共 ${t} 条`,
          onChange: (p, ps) => { setPage(p); setPageSize(ps) },
        }}
      />

      <Modal
        title={editingRegistry ? '编辑镜像仓库' : '新建镜像仓库'}
        open={modalOpen}
        onCancel={closeModal}
        onOk={handleSubmit}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
        okText={editingRegistry ? '保存' : '创建'}
        cancelText="取消"
        destroyOnClose
        width={600}
      >
        <Form
          form={form}
          layout="vertical"
          style={{ marginTop: 16 }}
          initialValues={{ type: 'harbor' }}
        >
          <Form.Item
            name="name"
            label="仓库名称"
            rules={[{ required: true, message: '请输入仓库名称' }]}
          >
            <Input placeholder="例如：生产环境Harbor" />
          </Form.Item>

          <Form.Item
            name="type"
            label="仓库类型"
            rules={[{ required: true, message: '请选择仓库类型' }]}
          >
            <Select
              options={[
                { value: 'harbor', label: 'Harbor' },
                { value: 'dockerhub', label: 'Docker Hub' },
                { value: 'acr', label: 'ACR' },
                { value: 'ecr', label: 'ECR' },
                { value: 'ghcr', label: 'GHCR' },
              ]}
            />
          </Form.Item>

          <Form.Item
            name="endpoint"
            label="端点地址"
            rules={[{ required: true, message: '请输入端点地址' }]}
          >
            <Input placeholder="https://harbor.example.com" />
          </Form.Item>

          <Form.Item
            name="username"
            label="用户名"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input placeholder="admin" />
          </Form.Item>

          <Form.Item
            name="password"
            label="密码"
            rules={[{ required: !editingRegistry, message: '请输入密码' }]}
          >
            <Input.Password placeholder={editingRegistry ? '留空则不修改' : '请输入密码'} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
