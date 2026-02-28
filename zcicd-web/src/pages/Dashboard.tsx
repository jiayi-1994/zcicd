import { useNavigate } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import {
  Card, Row, Col, Statistic, Typography, Tag, Table, Skeleton,
} from 'antd'
import {
  ProjectOutlined, CloudServerOutlined, RocketOutlined,
  InboxOutlined,
} from '@ant-design/icons'
import { projectApi } from '@/api/project'
import { deployApi } from '@/api/deploy'
import { systemApi } from '@/api/system'
import { artifactApi } from '@/api/artifact'

const { Title } = Typography

export default function Dashboard() {
  const navigate = useNavigate()

  const { data: projectData, isLoading: loadingProjects } = useQuery({
    queryKey: ['dashboard-projects'],
    queryFn: () => projectApi.list({ page: 1, page_size: 5 }) as Promise<any>,
  })

  const { data: deployData, isLoading: loadingDeploys } = useQuery({
    queryKey: ['dashboard-deploys'],
    queryFn: () => deployApi.list({ page: 1, page_size: 5 }) as Promise<any>,
  })

  const { data: clusterData, isLoading: loadingClusters } = useQuery({
    queryKey: ['dashboard-clusters'],
    queryFn: () => systemApi.listClusters({ page: 1, page_size: 100 }) as Promise<any>,
  })

  const { data: registryData, isLoading: loadingRegistries } = useQuery({
    queryKey: ['dashboard-registries'],
    queryFn: () => artifactApi.listRegistries({ page: 1, page_size: 100 }) as Promise<any>,
  })

  const projects = projectData?.data ?? []
  const deploys = deployData?.data ?? []
  const clusters = clusterData?.data ?? []
  const registries = registryData?.data ?? []

  const projectTotal = projectData?.pagination?.total ?? projects.length
  const deployTotal = deployData?.pagination?.total ?? deploys.length
  const clusterTotal = clusters.length
  const registryTotal = registries.length

  const recentProjectColumns = [
    { title: '项目名称', dataIndex: 'name', key: 'name' },
    {
      title: '描述', dataIndex: 'description', key: 'description',
      ellipsis: true,
      render: (text: string) => text || '-',
    },
    {
      title: '创建时间', dataIndex: 'created_at', key: 'created_at',
      width: 180,
      render: (t: string) => t ? new Date(t).toLocaleString('zh-CN') : '-',
    },
  ]

  const recentDeployColumns = [
    { title: '服务', dataIndex: 'service_id', key: 'service_id', width: 150 },
    {
      title: '策略', dataIndex: 'deploy_strategy', key: 'deploy_strategy', width: 120,
      render: (s: string) => {
        const map: Record<string, string> = { rolling: '滚动更新', 'blue-green': '蓝绿部署', canary: '金丝雀' }
        return map[s] || s
      },
    },
    {
      title: '自动同步', dataIndex: 'auto_sync', key: 'auto_sync', width: 100,
      render: (v: boolean) => <Tag color={v ? 'green' : 'default'}>{v ? '开启' : '关闭'}</Tag>,
    },
    {
      title: '需审批', dataIndex: 'require_approval', key: 'require_approval', width: 100,
      render: (v: boolean) => <Tag color={v ? 'orange' : 'default'}>{v ? '是' : '否'}</Tag>,
    },
  ]

  const isLoading = loadingProjects || loadingDeploys || loadingClusters || loadingRegistries

  return (
    <div style={{ padding: 24 }}>
      <Title level={3} style={{ marginBottom: 24 }}>仪表盘</Title>

      {/* Stats cards */}
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={12} lg={6}>
          <Card hoverable onClick={() => navigate('/projects')}>
            {isLoading ? <Skeleton active paragraph={{ rows: 1 }} /> : (
              <Statistic title="项目总数" value={projectTotal}
                prefix={<ProjectOutlined style={{ color: '#1677ff' }} />} />
            )}
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card hoverable onClick={() => navigate('/system/clusters')}>
            {isLoading ? <Skeleton active paragraph={{ rows: 1 }} /> : (
              <Statistic title="集群数量" value={clusterTotal}
                prefix={<CloudServerOutlined style={{ color: '#52c41a' }} />} />
            )}
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card hoverable>
            {isLoading ? <Skeleton active paragraph={{ rows: 1 }} /> : (
              <Statistic title="部署配置" value={deployTotal}
                prefix={<RocketOutlined style={{ color: '#722ed1' }} />} />
            )}
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card hoverable onClick={() => navigate('/artifacts/registries')}>
            {isLoading ? <Skeleton active paragraph={{ rows: 1 }} /> : (
              <Statistic title="镜像仓库" value={registryTotal}
                prefix={<InboxOutlined style={{ color: '#fa8c16' }} />} />
            )}
          </Card>
        </Col>
      </Row>

      {/* Recent tables */}
      <Row gutter={[16, 16]}>
        <Col xs={24} lg={12}>
          <Card title="最近项目" extra={<a onClick={() => navigate('/projects')}>查看全部</a>}>
            <Table
              columns={recentProjectColumns}
              dataSource={projects}
              rowKey="id"
              loading={loadingProjects}
              pagination={false}
              size="small"
              onRow={(record: any) => ({
                style: { cursor: 'pointer' },
                onClick: () => navigate(`/projects/${record.id}`),
              })}
            />
          </Card>
        </Col>
        <Col xs={24} lg={12}>
          <Card title="最近部署" extra={<a onClick={() => navigate('/projects')}>查看全部</a>}>
            <Table
              columns={recentDeployColumns}
              dataSource={deploys}
              rowKey="id"
              loading={loadingDeploys}
              pagination={false}
              size="small"
            />
          </Card>
        </Col>
      </Row>
    </div>
  )
}
