import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Table, Button, Input, Modal, Form, Space, App, Popconfirm, Typography,
} from 'antd'
import {
  PlusOutlined, SearchOutlined, EyeOutlined, DeleteOutlined,
} from '@ant-design/icons'
import { artifactApi, HelmChart } from '@/api/artifact'

const { Title } = Typography

export default function ArtifactCharts() {
  const queryClient = useQueryClient()
  const { message } = App.useApp()
  const [form] = Form.useForm()

  const [keyword, setKeyword] = useState('')
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [modalOpen, setModalOpen] = useState(false)

  // --- Data fetching ---
  const { data, isLoading } = useQuery({
    queryKey: ['charts', page, pageSize, keyword],
    queryFn: async () => {
      const res: any = await artifactApi.listCharts({ page, page_size: pageSize })
      return res
    },
  })

  const charts: HelmChart[] = data?.data ?? []
  const total: number = data?.pagination?.total ?? 0

  // --- Mutations ---
  const createMutation = useMutation({
    mutationFn: (values: Partial<HelmChart>) => artifactApi.createChart(values),
    onSuccess: () => {
      message.success('Helm Chart创建成功')
      queryClient.invalidateQueries({ queryKey: ['charts'] })
      closeModal()
    },
    onError: (err: any) => message.error(err?.message || '创建失败'),
  })

  const deleteMutation = useMutation({
    mutationFn: (name: string) => artifactApi.deleteChart(name),
    onSuccess: () => {
      message.success('Helm Chart已删除')
      queryClient.invalidateQueries({ queryKey: ['charts'] })
    },
    onError: (err: any) => message.error(err?.message || '删除失败'),
  })

  // --- Modal helpers ---
  const openCreate = () => {
    form.resetFields()
    setModalOpen(true)
  }

  const closeModal = () => {
    setModalOpen(false)
    form.resetFields()
  }

  const handleSubmit = async () => {
    const values = await form.validateFields()
    createMutation.mutate(values)
  }

  const handleSearch = (value: string) => {
    setKeyword(value)
    setPage(1)
  }

  // --- Table columns ---
  const columns = [
    {
      title: 'Chart名称',
      dataIndex: 'name',
      key: 'name',
      width: 200,
    },
    {
      title: '版本',
      dataIndex: 'version',
      key: 'version',
      width: 120,
    },
    {
      title: '应用版本',
      dataIndex: 'app_version',
      key: 'app_version',
      width: 120,
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
    },
    {
      title: '仓库地址',
      dataIndex: 'repo_url',
      key: 'repo_url',
      ellipsis: true,
    },
    {
      title: '操作',
      key: 'actions',
      width: 150,
      render: (_: any, record: HelmChart) => (
        <Space size="small">
          <Button
            type="link"
            size="small"
            icon={<EyeOutlined />}
            onClick={() => message.info('查看功能待实现')}
          >
            查看
          </Button>
          <Popconfirm
            title="确认删除"
            description={`确定要删除Chart「${record.name}」吗？`}
            onConfirm={() => deleteMutation.mutate(record.name)}
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
        <Title level={3} style={{ margin: 0 }}>Helm Chart管理</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
          新建Chart
        </Button>
      </div>

      <div style={{ marginBottom: 16 }}>
        <Input.Search
          placeholder="搜索Chart..."
          allowClear
          enterButton={<><SearchOutlined /> 搜索</>}
          onSearch={handleSearch}
          style={{ maxWidth: 400 }}
        />
      </div>

      <Table
        columns={columns}
        dataSource={charts}
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
        title="新建Helm Chart"
        open={modalOpen}
        onCancel={closeModal}
        onOk={handleSubmit}
        confirmLoading={createMutation.isPending}
        okText="创建"
        cancelText="取消"
        destroyOnClose
        width={600}
      >
        <Form
          form={form}
          layout="vertical"
          style={{ marginTop: 16 }}
        >
          <Form.Item
            name="name"
            label="Chart名称"
            rules={[{ required: true, message: '请输入Chart名称' }]}
          >
            <Input placeholder="例如：nginx" />
          </Form.Item>

          <Form.Item
            name="version"
            label="Chart版本"
            rules={[{ required: true, message: '请输入Chart版本' }]}
          >
            <Input placeholder="例如：1.0.0" />
          </Form.Item>

          <Form.Item
            name="app_version"
            label="应用版本"
            rules={[{ required: true, message: '请输入应用版本' }]}
          >
            <Input placeholder="例如：1.21.0" />
          </Form.Item>

          <Form.Item
            name="description"
            label="描述"
          >
            <Input.TextArea rows={3} placeholder="简要描述Chart用途..." />
          </Form.Item>

          <Form.Item
            name="repo_url"
            label="仓库地址"
            rules={[{ required: true, message: '请输入仓库地址' }]}
          >
            <Input placeholder="https://charts.example.com" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
