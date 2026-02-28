import { useState } from 'react'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import { Layout, Menu, Avatar, Dropdown, theme } from 'antd'
import {
  ProjectOutlined,
  DashboardOutlined,
  SettingOutlined,
  LogoutOutlined,
  UserOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  CloudServerOutlined,
  ApiOutlined,
  AuditOutlined,
  BellOutlined,
  InboxOutlined,
} from '@ant-design/icons'
import { useAuthStore } from '@/store/auth'

const { Header, Sider, Content } = Layout

const menuItems = [
  { key: '/dashboard', icon: <DashboardOutlined />, label: '仪表盘' },
  { key: '/projects', icon: <ProjectOutlined />, label: '项目管理' },
  { key: '/artifacts/registries', icon: <InboxOutlined />, label: '制品管理' },
  {
    key: 'system',
    icon: <SettingOutlined />,
    label: '系统管理',
    children: [
      { key: '/system/clusters', icon: <CloudServerOutlined />, label: '集群管理' },
      { key: '/system/integrations', icon: <ApiOutlined />, label: '集成管理' },
      { key: '/system/notifications', icon: <BellOutlined />, label: '通知管理' },
      { key: '/system/audit-logs', icon: <AuditOutlined />, label: '审计日志' },
    ],
  },
]

const getSelectedKey = (pathname: string): string => {
  if (pathname.startsWith('/system/')) return pathname
  if (pathname.startsWith('/artifacts/')) return '/artifacts/registries'
  if (pathname.startsWith('/projects')) return '/projects'
  if (pathname.startsWith('/dashboard')) return '/dashboard'
  return pathname
}

const getOpenKeys = (pathname: string): string[] => {
  if (pathname.startsWith('/system/')) return ['system']
  return []
}

export default function MainLayout() {
  const [collapsed, setCollapsed] = useState(false)
  const navigate = useNavigate()
  const location = useLocation()
  const { user, logout } = useAuthStore()
  const { token: { colorBgContainer, borderRadiusLG } } = theme.useToken()

  const userMenuItems = [
    { key: 'profile', icon: <UserOutlined />, label: '个人设置' },
    { type: 'divider' as const },
    { key: 'logout', icon: <LogoutOutlined />, label: '退出登录', danger: true },
  ]

  const handleUserMenu = ({ key }: { key: string }) => {
    if (key === 'logout') {
      logout()
      navigate('/login')
    }
  }

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider trigger={null} collapsible collapsed={collapsed} theme="light"
        style={{ borderRight: '1px solid #f0f0f0' }}>
        <div style={{ height: 64, display: 'flex', alignItems: 'center', justifyContent: 'center',
          fontSize: collapsed ? 16 : 20, fontWeight: 700, color: '#1677ff' }}>
          {collapsed ? 'Z' : 'ZCICD'}
        </div>
        <Menu mode="inline"
          selectedKeys={[getSelectedKey(location.pathname)]}
          defaultOpenKeys={getOpenKeys(location.pathname)}
          items={menuItems}
          onClick={({ key }) => navigate(key)} />
      </Sider>
      <Layout>
        <Header style={{ padding: '0 24px', background: colorBgContainer,
          display: 'flex', alignItems: 'center', justifyContent: 'space-between',
          borderBottom: '1px solid #f0f0f0' }}>
          <div style={{ cursor: 'pointer', fontSize: 18 }}
            onClick={() => setCollapsed(!collapsed)}>
            {collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
          </div>
          <Dropdown menu={{ items: userMenuItems, onClick: handleUserMenu }}>
            <div style={{ cursor: 'pointer', display: 'flex', alignItems: 'center', gap: 8 }}>
              <Avatar icon={<UserOutlined />} />
              <span>{user?.username || '用户'}</span>
            </div>
          </Dropdown>
        </Header>
        <Content style={{ margin: 24, padding: 24, background: colorBgContainer,
          borderRadius: borderRadiusLG, minHeight: 280 }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}
