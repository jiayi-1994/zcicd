import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Table, Button, Tag, Modal, Form, Select, Typography, Space, App, Popconfirm, Switch, Input,
} from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { systemApi, NotifyChannel } from '@/api/system'
import type { ColumnsType } from 'antd/es/table'

const { Title } = Typography

const TYPE_MAP: Record<string, { color: string; label: string }> = {
  dingtalk: { color: 'blue', label: '钉钉' },
  wechat: { color: 'green', label: '企业微信' },
  slack: { color: 'purple', label: 'Slack' },
  email: { color: 'orange', label: '邮件' },
  webhook: { color: 'default', label: 'Webhook' },
}

export default function NotificationChannels() {
  const queryClient = useQueryClient()
  const { message } = App.useApp()
  const [form] = Form.useForm()

  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(20)
  const [modalOpen, setModalOpen] = useState(false)
  const [editingChannel, setEditingChannel] = useState<NotifyChannel | null>(null)

  // --- Data fetching ---
  const { data, isLoading } = useQuery({
    queryKey: ['notifyChannels', page, pageSize],
    queryFn: async () => {
      const res: any = await systemApi.listChannels({ page, page_size: pageSize })
      return res
    },
  })

  const channels: NotifyChannel[] = data?.data ?? []
  const total: number = data?.pagination?.total ?? 0

  // --- Mutations ---
  const createMutation = useMutation({
    mutationFn: (values: Partial<NotifyChannel>) => systemApi.createChannel(values),
    onSuccess: () => {
      message.success('通知渠道创建成功')
      queryClient.invalidateQueries({ queryKey: ['notifyChannels'] })
      closeModal()
    },
    onError: (err: any) => message.error(err?.message || '创建失败'),
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, values }: { id: string; values: Partial<NotifyChannel> }) =>
      systemApi.updateChannel(id, values),
    onSuccess: () => {
      message.success('通知渠道更新成功')
      queryClient.invalidateQueries({ queryKey: ['notifyChannels'] })
      closeModal()
    },
    onError: (err: any) => message.error(err?.message || '更新失败'),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => systemApi.deleteChannel(id),
    onSuccess: () => {
      message.success('通知渠道已删除')
      queryClient.invalidateQueries({ queryKey: ['notifyChannels'] })
    },
    onError: (err: any) => message.error(err?.message || '删除失败'),
  })

  const toggleMutation = useMutation({
    mutationFn: ({ id, enabled }: { id: string; enabled: boolean }) =>
      systemApi.updateChannel(id, { enabled }),
    onSuccess: () => {
      message.success('状态更新成功')
      queryClient.invalidateQueries({ queryKey: ['notifyChannels'] })
    },
    onError: (err: any) => message.error(err?.message || '更新失败'),
  })

  // --- Modal helpers ---
  const openCreate = () => {
    setEditingChannel(null)
    form.resetFields()
    form.setFieldsValue({ enabled: true })
    setModalOpen(true)
  }

  const openEdit = (channel: NotifyChannel) => {
    setEditingChannel(channel)
    form.setFieldsValue({
      name: channel.name,
      type: channel.type,
      enabled: channel.enabled,
    })
    setModalOpen(true)
  }

  const closeModal = () => {
    setModalOpen(false)
    setEditingChannel(null)
    form.resetFields()
  }

  const handleSubmit = async () => {
    const values = await form.validateFields()
    if (editingChannel) {
      updateMutation.mutate({ id: editingChannel.id, values })
    } else {
      createMutation.mutate(values)
    }
  }

  // --- Render helpers ---
  const getType = (type: string) => TYPE_MAP[type] || { color: 'default', label: type }

  const columns: ColumnsType<NotifyChannel> = [
    {
      title: '渠道名称',
      dataIndex: 'name',
      key: 'name',
      width: 200,
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      width: 120,
      render: (text: string) => {
        const type = getType(text)
        return <Tag color={type.color}>{type.label}</Tag>
      },
    },
    {
      title: '状态',
      dataIndex: 'enabled',
      key: 'enabled',
      width: 100,
      render: (enabled: boolean, record: NotifyChannel) => (
        <Switch
          checked={enabled}
          onChange={(checked) => toggleMutation.mutate({ id: record.id, enabled: checked })}
          loading={toggleMutation.isPending}
        />
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
      render: (text: string) => new Date(text).toLocaleString('zh-CN'),
    },
    {
      title: '操作',
      key: 'actions',
      width: 120,
      render: (_, record: NotifyChannel) => (
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
            description={`确定要删除通知渠道「${record.name}」吗？`}
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
      {/* Page header */}
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <Title level={3} style={{ margin: 0 }}>通知渠道</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
          新建渠道
        </Button>
      </div>

      {/* Table */}
      <Table
        columns={columns}
        dataSource={channels}
        rowKey="id"
        loading={isLoading}
        pagination={{
          current: page,
          pageSize: pageSize,
          total: total,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (t) => `共 ${t} 个渠道`,
          pageSizeOptions: [20, 50, 100],
          onChange: (p, ps) => {
            setPage(p)
            setPageSize(ps)
          },
        }}
      />

      {/* Create / Edit Modal */}
      <Modal
        title={editingChannel ? '编辑通知渠道' : '新建通知渠道'}
        open={modalOpen}
        onCancel={closeModal}
        onOk={handleSubmit}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
        okText={editingChannel ? '保存' : '创建'}
        cancelText="取消"
        destroyOnClose
        width={560}
      >
        <Form
          form={form}
          layout="vertical"
          style={{ marginTop: 16 }}
          initialValues={{ enabled: true }}
        >
          <Form.Item
            name="name"
            label="渠道名称"
            rules={[{ required: true, message: '请输入渠道名称' }]}
          >
            <Input placeholder="例如：运维钉钉群" maxLength={50} />
          </Form.Item>

          <Form.Item
            name="type"
            label="渠道类型"
            rules={[{ required: true, message: '请选择渠道类型' }]}
          >
            <Select
              options={[
                { value: 'dingtalk', label: '钉钉' },
                { value: 'wechat', label: '企业微信' },
                { value: 'slack', label: 'Slack' },
                { value: 'email', label: '邮件' },
                { value: 'webhook', label: 'Webhook' },
              ]}
            />
          </Form.Item>

          <Form.Item
            name="enabled"
            label="启用状态"
            valuePropName="checked"
          >
            <Switch />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
