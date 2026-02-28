import { useState } from 'react'
import { useParams } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Card, Row, Col, Button, Input, Tag, Modal, Form, Select, Typography,
  Skeleton, Empty, Popconfirm, Space, App, Tooltip, Switch,
} from 'antd'
import {
  PlusOutlined, SearchOutlined, EditOutlined, DeleteOutlined,
  CloudOutlined, ClusterOutlined,
} from '@ant-design/icons'
import { projectApi, Environment } from '@/api/project'

const { Title, Text } = Typography

const ENV_TYPE_MAP: Record<string, { color: string; label: string }> = {
  dev: { color: 'blue', label: '开发' },
  test: { color: 'orange', label: '测试' },
  staging: { color: 'purple', label: '预发布' },
  prod: { color: 'red', label: '生产' },
}

export default function EnvironmentList() {
  const { projectId } = useParams<{ projectId: string }>()
  const queryClient = useQueryClient()
  const { message } = App.useApp()
  const [form] = Form.useForm()

  const [keyword, setKeyword] = useState('')
  const [modalOpen, setModalOpen] = useState(false)
  const [editingEnv, setEditingEnv] = useState<Environment | null>(null)

  // --- Data fetching ---
  const { data, isLoading } = useQuery({
    queryKey: ['environments', projectId],
    queryFn: async () => {
      const res: any = await projectApi.listEnvironments(projectId!)
      return res
    },
    enabled: !!projectId,
  })

  const environments: Environment[] = data?.data ?? []
  const filteredEnvs = keyword
    ? environments.filter(env =>
        env.name.toLowerCase().includes(keyword.toLowerCase()) ||
        env.namespace.toLowerCase().includes(keyword.toLowerCase())
      )
    : environments

  // --- Mutations ---
  const createMutation = useMutation({
    mutationFn: (values: Partial<Environment>) => projectApi.createEnvironment(projectId!, values),
    onSuccess: () => {
      message.success('环境创建成功')
      queryClient.invalidateQueries({ queryKey: ['environments'] })
      closeModal()
    },
    onError: (err: any) => message.error(err?.message || '创建失败'),
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, values }: { id: string; values: Partial<Environment> }) =>
      projectApi.updateEnvironment(id, values),
    onSuccess: () => {
      message.success('环境更新成功')
      queryClient.invalidateQueries({ queryKey: ['environments'] })
      closeModal()
    },
    onError: (err: any) => message.error(err?.message || '更新失败'),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => projectApi.deleteEnvironment(id),
    onSuccess: () => {
      message.success('环境已删除')
      queryClient.invalidateQueries({ queryKey: ['environments'] })
    },
    onError: (err: any) => message.error(err?.message || '删除失败'),
  })

  // --- Modal helpers ---
  const openCreate = () => {
    setEditingEnv(null)
    form.resetFields()
    form.setFieldsValue({ env_type: 'dev', auto_deploy: false })
    setModalOpen(true)
  }

  const openEdit = (env: Environment) => {
    setEditingEnv(env)
    form.setFieldsValue({
      name: env.name,
      env_type: env.env_type,
      namespace: env.namespace,
      cluster: env.cluster,
      auto_deploy: env.auto_deploy,
    })
    setModalOpen(true)
  }

  const closeModal = () => {
    setModalOpen(false)
    setEditingEnv(null)
    form.resetFields()
  }

  const handleSubmit = async () => {
    const values = await form.validateFields()
    if (editingEnv) {
      updateMutation.mutate({ id: editingEnv.id, values })
    } else {
      createMutation.mutate(values)
    }
  }

  const handleSearch = (value: string) => {
    setKeyword(value)
  }

  // --- Render helpers ---
  const getEnvType = (type: string) => ENV_TYPE_MAP[type] || { color: 'default', label: type }

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

  const renderEnvironmentCard = (env: Environment) => {
    const envType = getEnvType(env.env_type)

    return (
      <Col xs={24} sm={12} lg={8} xl={6} key={env.id}>
        <Card
          hoverable
          style={{ height: '100%', display: 'flex', flexDirection: 'column' }}
          styles={{ body: { flex: 1, display: 'flex', flexDirection: 'column' } }}
          actions={[
            <Tooltip title="编辑" key="edit">
              <EditOutlined onClick={() => openEdit(env)} />
            </Tooltip>,
            <Popconfirm
              key="delete"
              title="确认删除"
              description={`确定要删除环境「${env.name}」吗？此操作不可恢复。`}
              onConfirm={() => deleteMutation.mutate(env.id)}
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
          {/* Title */}
          <Text strong ellipsis style={{ fontSize: 16, marginBottom: 8, display: 'block' }}>
            {env.name}
          </Text>

          {/* Tags */}
          <Space size={4} style={{ marginBottom: 12 }}>
            <Tag color={envType.color}>{envType.label}</Tag>
            {env.auto_deploy && <Tag color="green">自动部署</Tag>}
          </Space>

          {/* Meta info */}
          <div style={{ fontSize: 13, color: '#595959', display: 'flex', flexDirection: 'column', gap: 6, flex: 1 }}>
            <Space size={4}>
              <CloudOutlined />
              <Text type="secondary" style={{ fontSize: 13 }}>命名空间：{env.namespace}</Text>
            </Space>
            <Space size={4}>
              <ClusterOutlined />
              <Text type="secondary" style={{ fontSize: 13 }}>集群：{env.cluster}</Text>
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
        <Title level={3} style={{ margin: 0 }}>环境管理</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
          新建环境
        </Button>
      </div>

      {/* Search bar */}
      <div style={{ marginBottom: 24 }}>
        <Input.Search
          placeholder="搜索环境名称或命名空间..."
          allowClear
          enterButton={<><SearchOutlined /> 搜索</>}
          size="large"
          onSearch={handleSearch}
          style={{ maxWidth: 480 }}
        />
      </div>

      {/* Environment cards */}
      {isLoading ? (
        renderSkeletons()
      ) : filteredEnvs.length === 0 ? (
        <Empty
          description={keyword ? '未找到匹配的环境' : '暂无环境'}
          style={{ marginTop: 80 }}
        >
          {!keyword && (
            <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
              创建第一个环境
            </Button>
          )}
        </Empty>
      ) : (
        <Row gutter={[16, 16]}>
          {filteredEnvs.map(renderEnvironmentCard)}
        </Row>
      )}

      {/* Create / Edit Modal */}
      <Modal
        title={editingEnv ? '编辑环境' : '新建环境'}
        open={modalOpen}
        onCancel={closeModal}
        onOk={handleSubmit}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
        okText={editingEnv ? '保存' : '创建'}
        cancelText="取消"
        destroyOnClose
        width={560}
      >
        <Form
          form={form}
          layout="vertical"
          style={{ marginTop: 16 }}
          initialValues={{ env_type: 'dev', auto_deploy: false }}
        >
          <Form.Item
            name="name"
            label="环境名称"
            rules={[{ required: true, message: '请输入环境名称' }]}
          >
            <Input placeholder="例如：开发环境" maxLength={50} />
          </Form.Item>

          <Form.Item
            name="env_type"
            label="环境类型"
            rules={[{ required: true, message: '请选择环境类型' }]}
          >
            <Select
              options={[
                { value: 'dev', label: '开发' },
                { value: 'test', label: '测试' },
                { value: 'staging', label: '预发布' },
                { value: 'prod', label: '生产' },
              ]}
            />
          </Form.Item>

          <Form.Item
            name="namespace"
            label="命名空间"
            rules={[{ required: true, message: '请输入命名空间' }]}
          >
            <Input placeholder="例如：app-dev" maxLength={100} />
          </Form.Item>

          <Form.Item
            name="cluster"
            label="集群"
            rules={[{ required: true, message: '请输入集群名称' }]}
          >
            <Input placeholder="例如：k8s-cluster-1" maxLength={100} />
          </Form.Item>

          <Form.Item name="auto_deploy" label="自动部署" valuePropName="checked">
            <Switch />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
