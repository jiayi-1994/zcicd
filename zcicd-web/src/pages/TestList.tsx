import { useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Table, Button, Input, Tag, Modal, Form, Select, Switch, Space, App, Popconfirm, Typography,
} from 'antd'
import {
  PlusOutlined, SearchOutlined, EditOutlined, DeleteOutlined,
  PlayCircleOutlined, EyeOutlined,
} from '@ant-design/icons'
import { qualityApi, TestConfig } from '@/api/quality'

const { Title } = Typography

const TEST_TYPE_MAP: Record<string, { color: string; label: string }> = {
  unit: { color: 'blue', label: '单元测试' },
  integration: { color: 'green', label: '集成测试' },
  e2e: { color: 'purple', label: 'E2E测试' },
  performance: { color: 'orange', label: '性能测试' },
}

export default function TestList() {
  const { projectId } = useParams<{ projectId: string }>()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const { message } = App.useApp()
  const [form] = Form.useForm()

  const [keyword, setKeyword] = useState('')
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [modalOpen, setModalOpen] = useState(false)
  const [editingTest, setEditingTest] = useState<TestConfig | null>(null)

  // --- Data fetching ---
  const { data, isLoading } = useQuery({
    queryKey: ['tests', projectId, page, pageSize, keyword],
    queryFn: async () => {
      const res: any = await qualityApi.listTests({
        project_id: projectId,
        page,
        page_size: pageSize
      })
      return res
    },
    enabled: !!projectId,
  })

  const tests: TestConfig[] = data?.data ?? []
  const total: number = data?.pagination?.total ?? 0

  // --- Mutations ---
  const createMutation = useMutation({
    mutationFn: (values: Partial<TestConfig>) => qualityApi.createTest(projectId!, values),
    onSuccess: () => {
      message.success('测试配置创建成功')
      queryClient.invalidateQueries({ queryKey: ['tests'] })
      closeModal()
    },
    onError: (err: any) => message.error(err?.message || '创建失败'),
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, values }: { id: string; values: Partial<TestConfig> }) =>
      qualityApi.updateTest(projectId!, id, values),
    onSuccess: () => {
      message.success('测试配置更新成功')
      queryClient.invalidateQueries({ queryKey: ['tests'] })
      closeModal()
    },
    onError: (err: any) => message.error(err?.message || '更新失败'),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => qualityApi.deleteTest(projectId!, id),
    onSuccess: () => {
      message.success('测试配置已删除')
      queryClient.invalidateQueries({ queryKey: ['tests'] })
    },
    onError: (err: any) => message.error(err?.message || '删除失败'),
  })

  const triggerMutation = useMutation({
    mutationFn: (id: string) => qualityApi.triggerTest(projectId!, id),
    onSuccess: () => {
      message.success('测试已触发')
    },
    onError: (err: any) => message.error(err?.message || '触发失败'),
  })

  // --- Modal helpers ---
  const openCreate = () => {
    setEditingTest(null)
    form.resetFields()
    form.setFieldsValue({
      test_type: 'unit',
      branch: 'main',
      timeout: 300,
      enabled: true
    })
    setModalOpen(true)
  }

  const openEdit = (test: TestConfig) => {
    setEditingTest(test)
    form.setFieldsValue({
      name: test.name,
      test_type: test.test_type,
      framework: test.framework,
      repo_url: test.repo_url,
      branch: test.branch,
      test_command: test.test_command,
      timeout: test.timeout,
      enabled: test.enabled,
    })
    setModalOpen(true)
  }

  const closeModal = () => {
    setModalOpen(false)
    setEditingTest(null)
    form.resetFields()
  }

  const handleSubmit = async () => {
    const values = await form.validateFields()
    if (editingTest) {
      updateMutation.mutate({ id: editingTest.id, values })
    } else {
      createMutation.mutate(values)
    }
  }

  const handleSearch = (value: string) => {
    setKeyword(value)
    setPage(1)
  }

  const getTestType = (type: string) => TEST_TYPE_MAP[type] || TEST_TYPE_MAP.unit

  // --- Table columns ---
  const columns = [
    {
      title: '名称',
      dataIndex: 'name',
      key: 'name',
      width: 200,
    },
    {
      title: '测试类型',
      dataIndex: 'test_type',
      key: 'test_type',
      width: 120,
      render: (type: string) => {
        const t = getTestType(type)
        return <Tag color={t.color}>{t.label}</Tag>
      },
    },
    {
      title: '测试框架',
      dataIndex: 'framework',
      key: 'framework',
      width: 120,
    },
    {
      title: '仓库地址',
      dataIndex: 'repo_url',
      key: 'repo_url',
      ellipsis: true,
    },
    {
      title: '启用状态',
      dataIndex: 'enabled',
      key: 'enabled',
      width: 100,
      render: (enabled: boolean) => (
        <Tag color={enabled ? 'green' : 'default'}>
          {enabled ? '已启用' : '已禁用'}
        </Tag>
      ),
    },
    {
      title: '操作',
      key: 'actions',
      width: 200,
      render: (_: any, record: TestConfig) => (
        <Space size="small">
          <Button
            type="link"
            size="small"
            icon={<PlayCircleOutlined />}
            onClick={() => triggerMutation.mutate(record.id)}
          >
            触发
          </Button>
          <Button
            type="link"
            size="small"
            icon={<EyeOutlined />}
            onClick={() => navigate(`/projects/${projectId}/tests/${record.id}/runs`)}
          >
            查看运行
          </Button>
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
            description={`确定要删除测试配置「${record.name}」吗？`}
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
        <Title level={3} style={{ margin: 0 }}>测试配置</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
          新建测试配置
        </Button>
      </div>

      <div style={{ marginBottom: 16 }}>
        <Input.Search
          placeholder="搜索测试配置..."
          allowClear
          enterButton={<><SearchOutlined /> 搜索</>}
          onSearch={handleSearch}
          style={{ maxWidth: 400 }}
        />
      </div>

      <Table
        columns={columns}
        dataSource={tests}
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
        title={editingTest ? '编辑测试配置' : '新建测试配置'}
        open={modalOpen}
        onCancel={closeModal}
        onOk={handleSubmit}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
        okText={editingTest ? '保存' : '创建'}
        cancelText="取消"
        destroyOnClose
        width={600}
      >
        <Form
          form={form}
          layout="vertical"
          style={{ marginTop: 16 }}
          initialValues={{ test_type: 'unit', branch: 'main', timeout: 300, enabled: true }}
        >
          <Form.Item
            name="name"
            label="配置名称"
            rules={[{ required: true, message: '请输入配置名称' }]}
          >
            <Input placeholder="例如：单元测试" />
          </Form.Item>

          <Form.Item
            name="test_type"
            label="测试类型"
            rules={[{ required: true, message: '请选择测试类型' }]}
          >
            <Select
              options={[
                { value: 'unit', label: '单元测试' },
                { value: 'integration', label: '集成测试' },
                { value: 'e2e', label: 'E2E测试' },
                { value: 'performance', label: '性能测试' },
              ]}
            />
          </Form.Item>

          <Form.Item
            name="framework"
            label="测试框架"
            rules={[{ required: true, message: '请输入测试框架' }]}
          >
            <Input placeholder="例如：jest, pytest, cypress" />
          </Form.Item>

          <Form.Item
            name="repo_url"
            label="仓库地址"
            rules={[{ required: true, message: '请输入仓库地址' }]}
          >
            <Input placeholder="https://github.com/org/repo" />
          </Form.Item>

          <Form.Item
            name="branch"
            label="分支"
            rules={[{ required: true, message: '请输入分支' }]}
          >
            <Input placeholder="main" />
          </Form.Item>

          <Form.Item
            name="test_command"
            label="测试命令"
            rules={[{ required: true, message: '请输入测试命令' }]}
          >
            <Input placeholder="例如：npm test" />
          </Form.Item>

          <Form.Item
            name="timeout"
            label="超时时间（秒）"
            rules={[{ required: true, message: '请输入超时时间' }]}
          >
            <Input type="number" placeholder="300" />
          </Form.Item>

          <Form.Item name="enabled" label="启用状态" valuePropName="checked">
            <Switch />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
