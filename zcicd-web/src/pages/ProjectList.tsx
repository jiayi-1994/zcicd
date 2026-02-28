import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Card, Row, Col, Button, Input, Tag, Modal, Form, Select, Typography,
  Pagination, Skeleton, Empty, Popconfirm, Space, App, Tooltip,
} from 'antd'
import {
  PlusOutlined, SearchOutlined, EditOutlined, DeleteOutlined,
  EyeOutlined, GithubOutlined, LockOutlined, GlobalOutlined,
  BranchesOutlined, ClockCircleOutlined,
} from '@ant-design/icons'
import { projectApi, Project } from '@/api/project'
import dayjs from 'dayjs'

const { Title, Text, Paragraph } = Typography

const STATUS_MAP: Record<string, { color: string; label: string; dot: string }> = {
  active: { color: 'green', label: '活跃', dot: '#52c41a' },
  archived: { color: 'default', label: '已归档', dot: '#bfbfbf' },
  inactive: { color: 'orange', label: '未激活', dot: '#faad14' },
}

const VISIBILITY_MAP: Record<string, { color: string; label: string; icon: React.ReactNode }> = {
  public: { color: 'green', label: '公开', icon: <GlobalOutlined /> },
  private: { color: 'blue', label: '私有', icon: <LockOutlined /> },
}

export default function ProjectList() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const { message } = App.useApp()
  const [form] = Form.useForm()

  const [keyword, setKeyword] = useState('')
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(12)
  const [modalOpen, setModalOpen] = useState(false)
  const [editingProject, setEditingProject] = useState<Project | null>(null)

  // --- Data fetching ---
  const { data, isLoading } = useQuery({
    queryKey: ['projects', page, pageSize, keyword],
    queryFn: async () => {
      const res: any = await projectApi.list({ page, page_size: pageSize, keyword: keyword || undefined })
      return res
    },
  })

  const projects: Project[] = data?.data ?? []
  const total: number = data?.pagination?.total ?? 0

  // --- Mutations ---
  const createMutation = useMutation({
    mutationFn: (values: Partial<Project>) => projectApi.create(values),
    onSuccess: () => {
      message.success('项目创建成功')
      queryClient.invalidateQueries({ queryKey: ['projects'] })
      closeModal()
    },
    onError: (err: any) => message.error(err?.message || '创建失败'),
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, values }: { id: string; values: Partial<Project> }) =>
      projectApi.update(id, values),
    onSuccess: () => {
      message.success('项目更新成功')
      queryClient.invalidateQueries({ queryKey: ['projects'] })
      closeModal()
    },
    onError: (err: any) => message.error(err?.message || '更新失败'),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => projectApi.delete(id),
    onSuccess: () => {
      message.success('项目已删除')
      queryClient.invalidateQueries({ queryKey: ['projects'] })
    },
    onError: (err: any) => message.error(err?.message || '删除失败'),
  })

  // --- Modal helpers ---
  const openCreate = () => {
    setEditingProject(null)
    form.resetFields()
    form.setFieldsValue({ visibility: 'private', default_branch: 'main' })
    setModalOpen(true)
  }

  const openEdit = (project: Project) => {
    setEditingProject(project)
    form.setFieldsValue({
      name: project.name,
      description: project.description,
      repo_url: project.repo_url,
      default_branch: project.default_branch,
      visibility: project.visibility,
    })
    setModalOpen(true)
  }

  const closeModal = () => {
    setModalOpen(false)
    setEditingProject(null)
    form.resetFields()
  }

  const handleSubmit = async () => {
    const values = await form.validateFields()
    if (editingProject) {
      updateMutation.mutate({ id: editingProject.id, values })
    } else {
      createMutation.mutate(values)
    }
  }

  const handleSearch = (value: string) => {
    setKeyword(value)
    setPage(1)
  }

  // --- Render helpers ---
  const getStatus = (status: string) => STATUS_MAP[status] || STATUS_MAP.active
  const getVisibility = (v: string) => VISIBILITY_MAP[v] || VISIBILITY_MAP.private

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

  const renderProjectCard = (project: Project) => {
    const status = getStatus(project.status)
    const visibility = getVisibility(project.visibility)

    return (
      <Col xs={24} sm={12} lg={8} xl={6} key={project.id}>
        <Card
          hoverable
          style={{ height: '100%', display: 'flex', flexDirection: 'column' }}
          styles={{ body: { flex: 1, display: 'flex', flexDirection: 'column' } }}
          actions={[
            <Tooltip title="查看详情" key="view">
              <EyeOutlined onClick={() => navigate(`/projects/${project.id}`)} />
            </Tooltip>,
            <Tooltip title="编辑" key="edit">
              <EditOutlined onClick={() => openEdit(project)} />
            </Tooltip>,
            <Popconfirm
              key="delete"
              title="确认删除"
              description={`确定要删除项目「${project.name}」吗？此操作不可恢复。`}
              onConfirm={() => deleteMutation.mutate(project.id)}
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
            <span
              style={{
                width: 8, height: 8, borderRadius: '50%',
                backgroundColor: status.dot, flexShrink: 0,
              }}
            />
            <Text strong ellipsis style={{ flex: 1, fontSize: 16 }}>
              {project.name}
            </Text>
          </div>

          {/* Tags */}
          <Space size={4} style={{ marginBottom: 8 }}>
            <Tag color={status.color}>{status.label}</Tag>
            <Tag icon={visibility.icon} color={visibility.color}>{visibility.label}</Tag>
          </Space>

          {/* Identifier + Description */}
          <Text code style={{ fontSize: 12, marginBottom: 4, display: 'inline-block' }}>
            {project.identifier}
          </Text>
          <Paragraph
            type="secondary"
            ellipsis={{ rows: 2 }}
            style={{ fontSize: 13, marginBottom: 12, flex: 1 }}
          >
            {project.description || '暂无描述'}
          </Paragraph>

          {/* Meta info */}
          <div style={{ fontSize: 12, color: '#8c8c8c', display: 'flex', flexDirection: 'column', gap: 4 }}>
            {project.repo_url && (
              <Tooltip title={project.repo_url}>
                <Space size={4}>
                  <GithubOutlined />
                  <Text type="secondary" ellipsis style={{ fontSize: 12, maxWidth: 180 }}>
                    {project.repo_url.replace(/^https?:\/\/(www\.)?/, '')}
                  </Text>
                </Space>
              </Tooltip>
            )}
            <Space size={12}>
              {project.default_branch && (
                <Space size={4}>
                  <BranchesOutlined />
                  <Text type="secondary" style={{ fontSize: 12 }}>{project.default_branch}</Text>
                </Space>
              )}
              <Space size={4}>
                <ClockCircleOutlined />
                <Text type="secondary" style={{ fontSize: 12 }}>
                  {dayjs(project.created_at).format('YYYY-MM-DD')}
                </Text>
              </Space>
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
        <Title level={3} style={{ margin: 0 }}>项目管理</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
          新建项目
        </Button>
      </div>

      {/* Search bar */}
      <div style={{ marginBottom: 24 }}>
        <Input.Search
          placeholder="搜索项目名称或标识..."
          allowClear
          enterButton={<><SearchOutlined /> 搜索</>}
          size="large"
          onSearch={handleSearch}
          style={{ maxWidth: 480 }}
        />
      </div>

      {/* Project cards */}
      {isLoading ? (
        renderSkeletons()
      ) : projects.length === 0 ? (
        <Empty
          description="暂无项目"
          style={{ marginTop: 80 }}
        >
          <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
            创建第一个项目
          </Button>
        </Empty>
      ) : (
        <>
          <Row gutter={[16, 16]}>
            {projects.map(renderProjectCard)}
          </Row>
          <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: 24 }}>
            <Pagination
              current={page}
              pageSize={pageSize}
              total={total}
              showSizeChanger
              showQuickJumper
              showTotal={(t) => `共 ${t} 个项目`}
              pageSizeOptions={[12, 24, 48]}
              onChange={(p, ps) => { setPage(p); setPageSize(ps) }}
            />
          </div>
        </>
      )}
      {/* Create / Edit Modal */}
      <Modal
        title={editingProject ? '编辑项目' : '新建项目'}
        open={modalOpen}
        onCancel={closeModal}
        onOk={handleSubmit}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
        okText={editingProject ? '保存' : '创建'}
        cancelText="取消"
        destroyOnClose
        width={560}
      >
        <Form
          form={form}
          layout="vertical"
          style={{ marginTop: 16 }}
          initialValues={{ visibility: 'private', default_branch: 'main' }}
        >
          <Form.Item
            name="name"
            label="项目名称"
            rules={[{ required: true, message: '请输入项目名称' }]}
          >
            <Input placeholder="例如：用户中心" maxLength={50} />
          </Form.Item>

          {!editingProject && (
            <Form.Item
              name="identifier"
              label="项目标识"
              rules={[
                { required: true, message: '请输入项目标识' },
                { pattern: /^[a-z][a-z0-9-]*$/, message: '仅支持小写字母、数字和连字符，且以字母开头' },
              ]}
              extra="创建后不可修改，用于系统内部唯一标识"
            >
              <Input placeholder="例如：user-center" maxLength={50} />
            </Form.Item>
          )}

          <Form.Item name="description" label="项目描述">
            <Input.TextArea rows={3} placeholder="简要描述项目用途..." maxLength={500} showCount />
          </Form.Item>

          <Form.Item name="repo_url" label="仓库地址">
            <Input prefix={<GithubOutlined />} placeholder="https://github.com/org/repo" />
          </Form.Item>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="default_branch" label="默认分支">
                <Input prefix={<BranchesOutlined />} placeholder="main" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="visibility"
                label="可见性"
                rules={[{ required: true, message: '请选择可见性' }]}
              >
                <Select
                  options={[
                    { value: 'private', label: '私有' },
                    { value: 'public', label: '公开' },
                  ]}
                />
              </Form.Item>
            </Col>
          </Row>
        </Form>
      </Modal>
    </div>
  )
}