import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import {
  Table, Button, Tag, Select, Typography, Space, DatePicker, Card,
} from 'antd'
import { ReloadOutlined } from '@ant-design/icons'
import { systemApi, AuditLog } from '@/api/system'
import type { ColumnsType } from 'antd/es/table'
import dayjs, { Dayjs } from 'dayjs'

const { Title } = Typography
const { RangePicker } = DatePicker

const ACTION_MAP: Record<string, { color: string; label: string }> = {
  create: { color: 'green', label: '创建' },
  update: { color: 'blue', label: '更新' },
  delete: { color: 'red', label: '删除' },
  deploy: { color: 'purple', label: '部署' },
  approve: { color: 'orange', label: '审批' },
}

export default function AuditLogs() {
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(20)
  const [action, setAction] = useState<string | undefined>(undefined)
  const [resourceType, setResourceType] = useState<string | undefined>(undefined)
  const [dateRange, setDateRange] = useState<[Dayjs | null, Dayjs | null] | null>(null)

  // --- Data fetching ---
  const { data, isLoading, refetch } = useQuery({
    queryKey: ['auditLogs', page, pageSize, action, resourceType, dateRange],
    queryFn: async () => {
      const res: any = await systemApi.listAuditLogs({
        page,
        page_size: pageSize,
        action: action || undefined,
        resource_type: resourceType || undefined,
      })
      return res
    },
  })

  const logs: AuditLog[] = data?.data ?? []
  const total: number = data?.pagination?.total ?? 0

  // --- Render helpers ---
  const getAction = (action: string) => ACTION_MAP[action] || { color: 'default', label: action }

  const columns: ColumnsType<AuditLog> = [
    {
      title: '时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
      render: (text: string) => dayjs(text).format('YYYY-MM-DD HH:mm:ss'),
    },
    {
      title: '用户',
      dataIndex: 'username',
      key: 'username',
      width: 120,
    },
    {
      title: '操作',
      dataIndex: 'action',
      key: 'action',
      width: 100,
      render: (text: string) => {
        const action = getAction(text)
        return <Tag color={action.color}>{action.label}</Tag>
      },
    },
    {
      title: '资源类型',
      dataIndex: 'resource_type',
      key: 'resource_type',
      width: 120,
      render: (text: string) => <Tag>{text}</Tag>,
    },
    {
      title: '资源名称',
      dataIndex: 'resource_name',
      key: 'resource_name',
      ellipsis: true,
    },
    {
      title: 'IP 地址',
      dataIndex: 'ip_address',
      key: 'ip_address',
      width: 140,
    },
  ]

  return (
    <div style={{ padding: 24 }}>
      {/* Page header */}
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <Title level={3} style={{ margin: 0 }}>审计日志</Title>
        <Button icon={<ReloadOutlined />} onClick={() => refetch()}>
          刷新
        </Button>
      </div>

      {/* Filter bar */}
      <Card style={{ marginBottom: 16 }}>
        <Space size={12} wrap>
          <Select
            placeholder="操作类型"
            allowClear
            style={{ width: 120 }}
            value={action}
            onChange={setAction}
            options={[
              { value: 'create', label: '创建' },
              { value: 'update', label: '更新' },
              { value: 'delete', label: '删除' },
              { value: 'deploy', label: '部署' },
              { value: 'approve', label: '审批' },
            ]}
          />
          <Select
            placeholder="资源类型"
            allowClear
            style={{ width: 120 }}
            value={resourceType}
            onChange={setResourceType}
            options={[
              { value: 'project', label: '项目' },
              { value: 'pipeline', label: '流水线' },
              { value: 'cluster', label: '集群' },
              { value: 'integration', label: '集成' },
            ]}
          />
          <RangePicker
            value={dateRange}
            onChange={setDateRange}
            format="YYYY-MM-DD"
            style={{ width: 260 }}
          />
        </Space>
      </Card>

      {/* Table */}
      <Table
        columns={columns}
        dataSource={logs}
        rowKey="id"
        loading={isLoading}
        pagination={{
          current: page,
          pageSize: pageSize,
          total: total,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (t) => `共 ${t} 条记录`,
          pageSizeOptions: [20, 50, 100],
          onChange: (p, ps) => {
            setPage(p)
            setPageSize(ps)
          },
        }}
      />
    </div>
  )
}
