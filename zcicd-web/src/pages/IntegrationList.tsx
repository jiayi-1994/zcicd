import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Card, Row, Col, Button, Input, Tag, Modal, Form, Select, Typography,
  Pagination, Skeleton, Empty, Popconfirm, Space, App, Tooltip, Badge,
} from 'antd'
import {
  PlusOutlined, SearchOutlined, EditOutlined, DeleteOutlined,
  ApiOutlined, GithubOutlined, DatabaseOutlined,
  BugOutlined, BellOutlined,
} from '@ant-design/icons'
import { systemApi, Integration } from '@/api/system'

const { Title, Text } = Typography

const TYPE_MAP: Record<string, { color: string; label: string; icon: React.ReactNode }> = {
  git: { color: 'blue', label: 'Git', icon: <GithubOutlined /> },
  registry: { color: 'purple', label: '镜像仓库', icon: <DatabaseOutlined /> },
  sonar: { color: 'orange', label: '代码质量', icon: <BugOutlined /> },
  jira: { color: 'cyan', label: '项目管理', icon: <ApiOutlined /> },
  notify: { color: 'green', label: '通知', icon: <BellOutlined /> },
}

const PROVIDER_MAP: Record<string, { color: string; label: string }> = {
  github: { color: 'default', label: 'GitHub' },
  gitlab: { color: 'orange', label: 'GitLab' },
  harbor: { color: 'blue', label: 'Harbor' },
  sonarqube: { color: 'cyan', label: 'SonarQube' },
  jira: { color: 'blue', label: 'Jira' },
}

const STATUS_MAP: Record<string, { color: 'success' | 'error' | 'default'; label: string }> = {
  active: { color: 'success', label: '正常' },
  error: { color: 'error', label: '异常' },
  inactive: { color: 'default', label: '未激活' },
}

export default function IntegrationList() {
  const queryClient = useQueryClient()
  const { message } = App.useApp()
  const [form] = Form.useForm()

  const [keyword, setKeyword] = useState('')
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(12)
  const [modalOpen, setModalOpen] = useState(false)
  const [editingIntegration, setEditingIntegration] = useState<Integration | null>(null)

  // --- Data fetching ---
  const { data, isLoading } = useQuery({
    queryKey: ['integrations', page, pageSize, keyword],
    queryFn: async () => {
      const res: any = await systemApi.listIntegrations({ page, page_size: pageSize })
      return res
    },
  })

  const integrations: Integration[] = data?.data ?? []
  const total: number = data?.pagination?.total ?? 0

  // --- Mutations ---
  const createMutation = useMutation({
    mutationFn: (values: Partial<Integration>) => systemApi.createIntegration(values),
    onSuccess: () => {
      message.success('集成创建成功')
      queryClient.invalidateQueries({ queryKey: ['integrations'] })
      closeModal()
    },
    onError: (err: any) => message.error(err?.message || '创建失败'),
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, values }: { id: string; values: Partial<Integration> }) =>
      systemApi.updateIntegration(id, values),
    onSuccess: () => {
      message.success('集成更新成功')
      queryClient.invalidateQueries({ queryKey: ['integrations'] })
      closeModal()
    },
    onError: (err: any) => message.error(err?.message || '更新失败'),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => systemApi.deleteIntegration(id),
    onSuccess: () => {
      message.success('集成已删除')
      queryClient.invalidateQueries({ queryKey: ['integrations'] })
    },
    onError: (err: any) => message.error(err?.message || '删除失败'),
  })

  // --- Modal helpers ---
  const openCreate = () => {
    setEditingIntegration(null)
    form.resetFields()
    setModalOpen(true)
  }

  const openEdit = (integration: Integration) => {
    setEditingIntegration(integration)
    form.setFieldsValue({
      name: integration.name,
      type: integration.type,
      provider: integration.provider,
    })
    setModalOpen(true)
  }

  const closeModal = () => {
    setModalOpen(false)
    setEditingIntegration(null)
    form.resetFields()
  }

  const handleSubmit = async () => {
    const values = await form.validateFields()
    if (editingIntegration) {
      updateMutation.mutate({ id: editingIntegration.id, values })
    } else {
      createMutation.mutate(values)
    }
  }

  const handleSearch = (value: string) => {
    setKeyword(value)
    setPage(1)
  }

  // --- Render helpers ---
  const getType = (type: string) => TYPE_MAP[type] || { color: 'default', label: type, icon: <ApiOutlined /> }
  const getProvider = (provider: string) => PROVIDER_MAP[provider] || { color: 'default', label: provider }
  const getStatus = (status: string) => STATUS_MAP[status] || STATUS_MAP.inactive

  const renderSkeletons = () => (
    <Row gutter={[16, 16]}>
      {Array.from({ length: 6 }).map((_, i) => (
        <Col xs={24} sm={12} lg={8} xl={6} key={i}>
          <Card>
            <Skeleton active paragraph={{ rows: 3 }} />
          </Card>
        </Col>
      ))}
    </Row>
  )

  const renderIntegrationCard = (integration: Integration) => {
    const type = getType(integration.type)
    const provider = getProvider(integration.provider)
    const status = getStatus(integration.status)

    return (
      <Col xs={24} sm={12} lg={8} xl={6} key={integration.id}>
        <Card
          hoverable
          style={{ height: '100%', display: 'flex', flexDirection: 'column' }}
          styles={{ body: { flex: 1, display: 'flex', flexDirection: 'column' } }}
          actions={[
            <Tooltip title="编辑" key="edit">
              <EditOutlined onClick={() => openEdit(integration)} />
            </Tooltip>,
            <Popconfirm
              key="delete"
              title="确认删除"
              description={`确定要删除集成「${integration.name}」吗？此操作不可恢复。`}
              onConfirm={() => deleteMutation.mutate(integration.id)}
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
            {type.icon}
            <Text strong ellipsis style={{ flex: 1, fontSize: 16 }}>
              {integration.name}
            </Text>
          </div>

          {/* Tags */}
          <Space size={4} style={{ marginBottom: 12 }}>
            <Tag color={type.color}>{type.label}</Tag>
            <Tag color={provider.color}>{provider.label}</Tag>
            <Badge status={status.color} text={status.label} />
          </Space>

          {/* Meta info */}
          <div style={{ fontSize: 12, color: '#8c8c8c', marginTop: 'auto' }}>
            {integration.last_check_at && (
              <Text type="secondary" style={{ fontSize: 12 }}>
                最后检查: {new Date(integration.last_check_at).toLocaleString('zh-CN')}
              </Text>
            )}
          </div>
        </Card>
      </Col>
    )
  }

  return (
    <div style={{ padding: 24 }}>
      {/* Page header */}
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <Title level={3} style={{ margin: 0 }}>外部集成</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
          新建集成
        </Button>
      </div>

      {/* Search bar */}
      <div style={{ marginBottom: 24 }}>
        <Input.Search
          placeholder="搜索集成名称..."
          allowClear
          enterButton={<><SearchOutlined /> 搜索</>}
          size="large"
          onSearch={handleSearch}
          style={{ maxWidth: 480 }}
        />
      </div>

      {/* Integration cards */}
      {isLoading ? (
        renderSkeletons()
      ) : integrations.length === 0 ? (
        <Empty
          description="暂无集成"
          style={{ marginTop: 80 }}
        >
          <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
            创建第一个集成
          </Button>
        </Empty>
      ) : (
        <>
          <Row gutter={[16, 16]}>
            {integrations.map(renderIntegrationCard)}
          </Row>
          <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: 24 }}>
            <Pagination
              current={page}
              pageSize={pageSize}
              total={total}
              showSizeChanger
              showQuickJumper
              showTotal={(t) => `共 ${t} 个集成`}
              pageSizeOptions={[12, 24, 48]}
              onChange={(p, ps) => { setPage(p); setPageSize(ps) }}
            />
          </div>
        </>
      )}

      {/* Create / Edit Modal */}
      <Modal
        title={editingIntegration ? '编辑集成' : '新建集成'}
        open={modalOpen}
        onCancel={closeModal}
        onOk={handleSubmit}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
        okText={editingIntegration ? '保存' : '创建'}
        cancelText="取消"
        destroyOnClose
        width={560}
      >
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item
            name="name"
            label="集成名称"
            rules={[{ required: true, message: '请输入集成名称' }]}
          >
            <Input placeholder="例如：GitHub 主仓库" maxLength={50} />
          </Form.Item>

          <Form.Item
            name="type"
            label="集成类型"
            rules={[{ required: true, message: '请选择集成类型' }]}
          >
            <Select
              options={[
                { value: 'git', label: 'Git 仓库' },
                { value: 'registry', label: '镜像仓库' },
                { value: 'sonar', label: '代码质量' },
                { value: 'jira', label: '项目管理' },
                { value: 'notify', label: '通知服务' },
              ]}
            />
          </Form.Item>

          <Form.Item
            name="provider"
            label="提供商"
            rules={[{ required: true, message: '请选择提供商' }]}
          >
            <Select
              options={[
                { value: 'github', label: 'GitHub' },
                { value: 'gitlab', label: 'GitLab' },
                { value: 'harbor', label: 'Harbor' },
                { value: 'sonarqube', label: 'SonarQube' },
                { value: 'jira', label: 'Jira' },
              ]}
            />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
