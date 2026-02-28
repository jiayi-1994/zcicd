import { useParams, useNavigate } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import {
  Button, Tag, Typography, Space, Breadcrumb, Skeleton, Descriptions, Card,
} from 'antd'
import { ArrowLeftOutlined, HomeOutlined } from '@ant-design/icons'
import { workflowApi, BuildRun } from '@/api/workflow'
import { projectApi } from '@/api/project'
import dayjs from 'dayjs'

const { Title } = Typography

const BUILD_STATUS_MAP: Record<string, { color: string; label: string }> = {
  pending: { color: 'default', label: '等待中' },
  running: { color: 'processing', label: '运行中' },
  succeeded: { color: 'success', label: '成功' },
  failed: { color: 'error', label: '失败' },
  cancelled: { color: 'warning', label: '已取消' },
}

export default function BuildLog() {
  const { projectId, configId, runId } = useParams<{
    projectId: string; configId: string; runId: string
  }>()
  const navigate = useNavigate()

  const { data: projectData } = useQuery({
    queryKey: ['project', projectId],
    queryFn: () => projectApi.get(projectId!),
    enabled: !!projectId,
  })

  const { data: configData } = useQuery({
    queryKey: ['buildConfig', configId],
    queryFn: () => workflowApi.getBuildConfig(configId!),
    enabled: !!configId,
  })

  const { data: runData, isLoading } = useQuery({
    queryKey: ['buildRun', runId],
    queryFn: () => workflowApi.getBuildRun(runId!),
    enabled: !!runId,
    refetchInterval: (query) => {
      const status = (query.state.data as any)?.data?.status
      return (status === 'running' || status === 'pending') ? 5000 : false
    },
  })

  const projectName = (projectData as any)?.data?.name ?? '项目'
  const configName = (configData as any)?.data?.name ?? '构建配置'
  const run: BuildRun | undefined = (runData as any)?.data
  const statusCfg = BUILD_STATUS_MAP[run?.status ?? ''] ?? { color: 'default', label: run?.status ?? '' }

  if (isLoading) {
    return (
      <div style={{ padding: 24 }}>
        <Skeleton active paragraph={{ rows: 12 }} />
      </div>
    )
  }

  return (
    <div style={{ padding: 24 }}>
      <Breadcrumb
        style={{ marginBottom: 16 }}
        items={[
          { title: <><HomeOutlined /> 项目</>, href: '/projects', onClick: (e) => { e.preventDefault(); navigate('/projects') } },
          { title: projectName, href: `/projects/${projectId}`, onClick: (e) => { e.preventDefault(); navigate(`/projects/${projectId}`) } },
          { title: '构建管理', href: `/projects/${projectId}/builds`, onClick: (e) => { e.preventDefault(); navigate(`/projects/${projectId}/builds`) } },
          { title: configName, href: `/projects/${projectId}/builds/${configId}/runs`, onClick: (e) => { e.preventDefault(); navigate(`/projects/${projectId}/builds/${configId}/runs`) } },
          { title: `#${run?.run_number ?? ''}` },
        ]}
      />

      <div style={{ marginBottom: 16 }}>
        <Space>
          <Button icon={<ArrowLeftOutlined />} onClick={() => navigate(`/projects/${projectId}/builds/${configId}/runs`)}>
            返回
          </Button>
          <Title level={4} style={{ margin: 0 }}>构建日志 - #{run?.run_number}</Title>
        </Space>
      </div>

      <Card style={{ marginBottom: 16 }}>
        <Descriptions column={{ xs: 1, sm: 2, md: 3 }} size="small">
          <Descriptions.Item label="状态"><Tag color={statusCfg.color}>{statusCfg.label}</Tag></Descriptions.Item>
          <Descriptions.Item label="分支">{run?.branch || '-'}</Descriptions.Item>
          <Descriptions.Item label="Commit">{run?.commit_sha ? <code>{run.commit_sha.substring(0, 8)}</code> : '-'}</Descriptions.Item>
          <Descriptions.Item label="镜像 Tag">{run?.image_tag || '-'}</Descriptions.Item>
          <Descriptions.Item label="耗时">{run?.duration_sec ? `${run.duration_sec}s` : '-'}</Descriptions.Item>
          <Descriptions.Item label="触发者">{run?.triggered_by || '-'}</Descriptions.Item>
          <Descriptions.Item label="开始时间">{run?.started_at ? dayjs(run.started_at).format('YYYY-MM-DD HH:mm:ss') : '-'}</Descriptions.Item>
          <Descriptions.Item label="结束时间">{run?.finished_at ? dayjs(run.finished_at).format('YYYY-MM-DD HH:mm:ss') : '-'}</Descriptions.Item>
        </Descriptions>
      </Card>

      <div style={{
        background: '#1e1e1e',
        borderRadius: 8,
        padding: 24,
        minHeight: 400,
        fontFamily: "'Cascadia Code', 'Fira Code', 'Consolas', monospace",
        fontSize: 13,
        lineHeight: 1.6,
        color: '#d4d4d4',
        overflowX: 'auto',
        whiteSpace: 'pre-wrap',
        wordBreak: 'break-all',
      }}>
        {run?.status === 'pending' && (
          <span style={{ color: '#6a9955' }}>{'> '}等待构建开始...</span>
        )}
        {run?.status === 'running' && (
          <span style={{ color: '#dcdcaa' }}>{'> '}WebSocket 连接中... 日志加载中...</span>
        )}
        {run?.status === 'succeeded' && (
          <span style={{ color: '#6a9955' }}>{'> '}构建成功完成。日志详情将在 M4 版本通过 WebSocket 实时展示。</span>
        )}
        {run?.status === 'failed' && (
          <span style={{ color: '#f44747' }}>{'> '}构建失败。日志详情将在 M4 版本通过 WebSocket 实时展示。</span>
        )}
        {run?.status === 'cancelled' && (
          <span style={{ color: '#ce9178' }}>{'> '}构建已取消。</span>
        )}
        {!run?.status && (
          <span style={{ color: '#6a9955' }}>{'> '}日志加载中...</span>
        )}
      </div>
    </div>
  )
}
