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
import { qualityApi, ScanConfig } from '@/api/quality'

const { Title } = Typography

const SCAN_TYPE_MAP: Record<string, { color: string; label: string }> = {
  sonar: { color: 'blue', label: 'SonarQube' },
  sast: { color: 'red', label: 'SAST' },
  dependency: { color: 'orange', label: '依赖扫描' },
}

export default function ScanList() {
  const { projectId } = useParams<{ projectId: string }>()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const { message } = App.useApp()
  const [form] = Form.useForm()

  const [keyword, setKeyword] = useState('')
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [modalOpen, setModalOpen] = useState(false)
  const [editingScan, setEditingScan] = useState<ScanConfig | null>(null)

  // --- Data fetching ---
  const { data, isLoading } = useQuery({
    queryKey: ['scans', projectId, page, pageSize, keyword],
    queryFn: async () => {
      const res: any = await qualityApi.listScans({
        project_id: projectId,
        page,
        page_size: pageSize
      })
      return res
    },
    enabled: !!projectId,
  })

  const scans: ScanConfig[] = data?.data ?? []
  const total: number = data?.pagination?.total ?? 0

  // --- Mutations ---
  const createMutation = useMutation({
    mutationFn: (values: Partial<ScanConfig>) => qualityApi.createScan(projectId!, values),
    onSuccess: () => {
      message.success('扫描配置创建成功')
      queryClient.invalidateQueries({ queryKey: ['scans'] })
      closeModal()
    },
    onError: (err: any) => message.error(err?.message || '创建失败'),
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, values }: { id: string; values: Partial<ScanConfig> }) =>
      qualityApi.updateScan(projectId!, id, values),
    onSuccess: () => {
      message.success('扫描配置更新成功')
      queryClient.invalidateQueries({ queryKey: ['scans'] })
      closeModal()
    },
    onError: (err: any) => message.error(err?.message || '更新失败'),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => qualityApi.deleteScan(projectId!, id),
    onSuccess: () => {
      message.success('扫描配置已删除')
      queryClient.invalidateQueries({ queryKey: ['scans'] })
    },
    onError: (err: any) => message.error(err?.message || '删除失败'),
  })

  const triggerMutation = useMutation({
    mutationFn: (id: string) => qualityApi.triggerScan(projectId!, id),
    onSuccess: () => {
      message.success('扫描已触发')
    },
    onError: (err: any) => message.error(err?.message || '触发失败'),
  })

  // --- Modal helpers ---
  const openCreate = () => {
    setEditingScan(null)
    form.resetFields()
    form.setFieldsValue({
      scan_type: 'sonar',
      branch: 'main',
      enabled: true
    })
    setModalOpen(true)
  }

  const openEdit = (scan: ScanConfig) => {
    setEditingScan(scan)
    form.setFieldsValue({
      name: scan.name,
      scan_type: scan.scan_type,
      sonar_project_key: scan.sonar_project_key,
      repo_url: scan.repo_url,
      branch: scan.branch,
      enabled: scan.enabled,
    })
    setModalOpen(true)
  }

  const closeModal = () => {
    setModalOpen(false)
    setEditingScan(null)
    form.resetFields()
  }

  const handleSubmit = async () => {
    const values = await form.validateFields()
    if (editingScan) {
      updateMutation.mutate({ id: editingScan.id, values })
    } else {
      createMutation.mutate(values)
    }
  }

  const handleSearch = (value: string) => {
    setKeyword(value)
    setPage(1)
  }

  const getScanType = (type: string) => SCAN_TYPE_MAP[type] || SCAN_TYPE_MAP.sonar

  // --- Table columns ---
  const columns = [
    {
      title: '名称',
      dataIndex: 'name',
      key: 'name',
      width: 200,
    },
    {
      title: '扫描类型',
      dataIndex: 'scan_type',
      key: 'scan_type',
      width: 120,
      render: (type: string) => {
        const t = getScanType(type)
        return <Tag color={t.color}>{t.label}</Tag>
      },
    },
    {
      title: 'SonarQube项目Key',
      dataIndex: 'sonar_project_key',
      key: 'sonar_project_key',
      width: 180,
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
      render: (_: any, record: ScanConfig) => (
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
            onClick={() => navigate(`/projects/${projectId}/scans/${record.id}/runs`)}
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
            description={`确定要删除扫描配置「${record.name}」吗？`}
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
        <Title level={3} style={{ margin: 0 }}>代码扫描配置</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
          新建扫描配置
        </Button>
      </div>

      <div style={{ marginBottom: 16 }}>
        <Input.Search
          placeholder="搜索扫描配置..."
          allowClear
          enterButton={<><SearchOutlined /> 搜索</>}
          onSearch={handleSearch}
          style={{ maxWidth: 400 }}
        />
      </div>

      <Table
        columns={columns}
        dataSource={scans}
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
        title={editingScan ? '编辑扫描配置' : '新建扫描配置'}
        open={modalOpen}
        onCancel={closeModal}
        onOk={handleSubmit}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
        okText={editingScan ? '保存' : '创建'}
        cancelText="取消"
        destroyOnClose
        width={600}
      >
        <Form
          form={form}
          layout="vertical"
          style={{ marginTop: 16 }}
          initialValues={{ scan_type: 'sonar', branch: 'main', enabled: true }}
        >
          <Form.Item
            name="name"
            label="配置名称"
            rules={[{ required: true, message: '请输入配置名称' }]}
          >
            <Input placeholder="例如：代码质量扫描" />
          </Form.Item>

          <Form.Item
            name="scan_type"
            label="扫描类型"
            rules={[{ required: true, message: '请选择扫描类型' }]}
          >
            <Select
              options={[
                { value: 'sonar', label: 'SonarQube' },
                { value: 'sast', label: 'SAST' },
                { value: 'dependency', label: '依赖扫描' },
              ]}
            />
          </Form.Item>

          <Form.Item
            name="sonar_project_key"
            label="SonarQube项目Key"
          >
            <Input placeholder="例如：my-project" />
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

          <Form.Item name="enabled" label="启用状态" valuePropName="checked">
            <Switch />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
