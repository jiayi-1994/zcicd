import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Card, Row, Col, Button, Input, Tag, Modal, Form, Select, Typography,
  Pagination, Skeleton, Empty, Popconfirm, Space, App, Tooltip, Badge,
} from 'antd'
import {
  PlusOutlined, SearchOutlined, EditOutlined, DeleteOutlined,
  CloudServerOutlined, ApiOutlined, CheckCircleOutlined, CloseCircleOutlined,
} from '@ant-design/icons'
import { systemApi, Cluster } from '@/api/system'

const { Title, Text, Paragraph } = Typography

const PROVIDER_MAP: Record<string, { color: string; label: string }> = {
  'self-managed': { color: 'blue', label: '自建' },
  'eks': { color: 'orange', label: 'EKS' },
  'aks': { color: 'cyan', label: 'AKS' },
  'gke': { color: 'green', label: 'GKE' },
}

const STATUS_MAP: Record<string, { color: 'success' | 'error'; label: string; icon: React.ReactNode }> = {
  connected: { color: 'success', label: '已连接', icon: <CheckCircleOutlined /> },
  disconnected: { color: 'error', label: '未连接', icon: <CloseCircleOutlined /> },
}

export default function ClusterList() {
  const queryClient = useQueryClient()
  const { message } = App.useApp()
  const [form] = Form.useForm()

  const [keyword, setKeyword] = useState('')
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(12)
  const [modalOpen, setModalOpen] = useState(false)
  const [editingCluster, setEditingCluster] = useState<Cluster | null>(null)

  // --- Data fetching ---
  const { data, isLoading } = useQuery({
    queryKey: ['clusters', page, pageSize, keyword],
    queryFn: async () => {
      const res: any = await systemApi.listClusters({ page, page_size: pageSize })
      return res
    },
  })

  const clusters: Cluster[] = data?.data ?? []
  const total: number = data?.pagination?.total ?? 0

  // --- Mutations ---
  const createMutation = useMutation({
    mutationFn: (values: Partial<Cluster>) => systemApi.createCluster(values),
    onSuccess: () => {
      message.success('集群创建成功')
      queryClient.invalidateQueries({ queryKey: ['clusters'] })
      closeModal()
    },
    onError: (err: any) => message.error(err?.message || '创建失败'),
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, values }: { id: string; values: Partial<Cluster> }) =>
      systemApi.updateCluster(id, values),
    onSuccess: () => {
      message.success('集群更新成功')
      queryClient.invalidateQueries({ queryKey: ['clusters'] })
      closeModal()
    },
    onError: (err: any) => message.error(err?.message || '更新失败'),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => systemApi.deleteCluster(id),
    onSuccess: () => {
      message.success('集群已删除')
      queryClient.invalidateQueries({ queryKey: ['clusters'] })
    },
    onError: (err: any) => message.error(err?.message || '删除失败'),
  })

  // --- Modal helpers ---
  const openCreate = () => {
    setEditingCluster(null)
    form.resetFields()
    form.setFieldsValue({ provider: 'self-managed' })
    setModalOpen(true)
  }

  const openEdit = (cluster: Cluster) => {
    setEditingCluster(cluster)
    form.setFieldsValue({
      name: cluster.name,
      display_name: cluster.display_name,
      description: cluster.description,
      provider: cluster.provider,
      api_server_url: cluster.api_server_url,
    })
    setModalOpen(true)
  }

  const closeModal = () => {
    setModalOpen(false)
    setEditingCluster(null)
    form.resetFields()
  }

  const handleSubmit = async () => {
    const values = await form.validateFields()
    if (editingCluster) {
      updateMutation.mutate({ id: editingCluster.id, values })
    } else {
      createMutation.mutate(values)
    }
  }

  const handleSearch = (value: string) => {
    setKeyword(value)
    setPage(1)
  }

  // --- Render helpers ---
  const getProvider = (provider: string) => PROVIDER_MAP[provider] || { color: 'default', label: provider }
  const getStatus = (status: string) => STATUS_MAP[status] || STATUS_MAP.disconnected

  const renderSkeletons = () => (
    <Row gutter={[16, 16]}>
      {Array.from({ length: 6 }).map((_, i) => (
        <Col xs={24} sm={12} lg={8} xl={6} key={i}>
          <Card>
            <Skeleton active paragraph={{ rows: 4 }} />
          </Card>
        </Col>
      ))}
    </Row>
  )

  const renderClusterCard = (cluster: Cluster) => {
    const provider = getProvider(cluster.provider)
    const status = getStatus(cluster.status)

    return (
      <Col xs={24} sm={12} lg={8} xl={6} key={cluster.id}>
        <Card
          hoverable
          style={{ height: '100%', display: 'flex', flexDirection: 'column' }}
          styles={{ body: { flex: 1, display: 'flex', flexDirection: 'column' } }}
          actions={[
            <Tooltip title="编辑" key="edit">
              <EditOutlined onClick={() => openEdit(cluster)} />
            </Tooltip>,
            <Popconfirm
              key="delete"
              title="确认删除"
              description={`确定要删除集群「${cluster.display_name}」吗？此操作不可恢复。`}
              onConfirm={() => deleteMutation.mutate(cluster.id)}
              okText="删除"
              cancelText="取消"
              okButtonProps={{ danger: true }}
            >
              <Tooltip title="删除">
                <DeleteOutlined style={{ color: '#ff4d4f' }} />
              </Tooltip>
            </Popconfirm>,
          ]}
        >
          {/* Title row */}
          <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 8 }}>
            <CloudServerOutlined style={{ fontSize: 18, color: '#1890ff' }} />
            <Text strong ellipsis style={{ flex: 1, fontSize: 16 }}>
              {cluster.display_name}
            </Text>
          </div>

          {/* Tags */}
          <Space size={4} style={{ marginBottom: 8 }}>
            <Tag color={provider.color}>{provider.label}</Tag>
            <Badge status={status.color} text={status.label} />
          </Space>

          {/* Name code */}
          <Text code style={{ fontSize: 12, marginBottom: 4, display: 'inline-block' }}>
            {cluster.name}
          </Text>

          {/* Description */}
          <Paragraph
            type="secondary"
            ellipsis={{ rows: 2 }}
            style={{ fontSize: 13, marginBottom: 12, flex: 1 }}
          >
            {cluster.description || '暂无描述'}
          </Paragraph>

          {/* Meta info */}
          <div style={{ fontSize: 12, color: '#8c8c8c', display: 'flex', flexDirection: 'column', gap: 4 }}>
            <Tooltip title={cluster.api_server_url}>
              <Space size={4}>
                <ApiOutlined />
                <Text type="secondary" ellipsis style={{ fontSize: 12, maxWidth: 180 }}>
                  {cluster.api_server_url}
                </Text>
              </Space>
            </Tooltip>
            <Space size={12}>
              <Text type="secondary" style={{ fontSize: 12 }}>
                节点: {cluster.node_count}
              </Text>
              <Text type="secondary" style={{ fontSize: 12 }}>
                版本: {cluster.version}
              </Text>
            </Space>
          </div>
        </Card>
      </Col>
    )
  }

  return (
    <div style={{ padding: 24 }}>
      {/* Page header */}
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <Title level={3} style={{ margin: 0 }}>集群管理</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
          新建集群
        </Button>
      </div>

      {/* Search bar */}
      <div style={{ marginBottom: 24 }}>
        <Input.Search
          placeholder="搜索集群名称..."
          allowClear
          enterButton={<><SearchOutlined /> 搜索</>}
          size="large"
          onSearch={handleSearch}
          style={{ maxWidth: 480 }}
        />
      </div>

      {/* Cluster cards */}
      {isLoading ? (
        renderSkeletons()
      ) : clusters.length === 0 ? (
        <Empty
          description="暂无集群"
          style={{ marginTop: 80 }}
        >
          <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
            创建第一个集群
          </Button>
        </Empty>
      ) : (
        <>
          <Row gutter={[16, 16]}>
            {clusters.map(renderClusterCard)}
          </Row>
          <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: 24 }}>
            <Pagination
              current={page}
              pageSize={pageSize}
              total={total}
              showSizeChanger
              showQuickJumper
              showTotal={(t) => `共 ${t} 个集群`}
              pageSizeOptions={[12, 24, 48]}
              onChange={(p, ps) => { setPage(p); setPageSize(ps) }}
            />
          </div>
        </>
      )}

      {/* Create / Edit Modal */}
      <Modal
        title={editingCluster ? '编辑集群' : '新建集群'}
        open={modalOpen}
        onCancel={closeModal}
        onOk={handleSubmit}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
        okText={editingCluster ? '保存' : '创建'}
        cancelText="取消"
        destroyOnClose
        width={560}
      >
        <Form
          form={form}
          layout="vertical"
          style={{ marginTop: 16 }}
          initialValues={{ provider: 'self-managed' }}
        >
          <Form.Item
            name="name"
            label="集群名称"
            rules={[{ required: true, message: '请输入集群名称' }]}
          >
            <Input placeholder="例如：prod-k8s" maxLength={50} />
          </Form.Item>

          <Form.Item
            name="display_name"
            label="显示名称"
            rules={[{ required: true, message: '请输入显示名称' }]}
          >
            <Input placeholder="例如：生产环境集群" maxLength={50} />
          </Form.Item>

          <Form.Item name="description" label="集群描述">
            <Input.TextArea rows={3} placeholder="简要描述集群用途..." maxLength={500} showCount />
          </Form.Item>

          <Form.Item
            name="provider"
            label="提供商"
            rules={[{ required: true, message: '请选择提供商' }]}
          >
            <Select
              options={[
                { value: 'self-managed', label: '自建' },
                { value: 'eks', label: 'AWS EKS' },
                { value: 'aks', label: 'Azure AKS' },
                { value: 'gke', label: 'Google GKE' },
              ]}
            />
          </Form.Item>

          <Form.Item
            name="api_server_url"
            label="API Server 地址"
            rules={[{ required: true, message: '请输入 API Server 地址' }]}
          >
            <Input prefix={<ApiOutlined />} placeholder="https://k8s.example.com:6443" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
