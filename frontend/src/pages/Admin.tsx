// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { useAuthStore } from '@/stores/authStore'
import Dashboard from './admin/Dashboard'
import Categories from './admin/Categories'
import Links from './admin/Links'
import Tags from './admin/Tags'
import Users from './admin/Users'
import Settings from './admin/Settings'
import Tokens from './admin/Tokens'
import {
  LayoutDashboard,
  Folder,
  Link2,
  Tag,
  Users as UsersIcon,
  Settings as SettingsIcon,
  Key,
} from 'lucide-react'

type TabType = 'dashboard' | 'categories' | 'links' | 'tags' | 'users' | 'tokens' | 'settings'

export default function Admin() {
  const navigate = useNavigate()
  const { user, isAuthenticated } = useAuthStore()
  const [activeTab, setActiveTab] = useState<TabType>('dashboard')

  useEffect(() => {
    if (!isAuthenticated || user?.role !== 'admin') {
      navigate('/')
    }
  }, [isAuthenticated, user, navigate])

  if (!isAuthenticated || user?.role !== 'admin') {
    return null
  }

  const tabs = [
    { id: 'dashboard' as TabType, label: '仪表盘', icon: LayoutDashboard },
    { id: 'categories' as TabType, label: '分类管理', icon: Folder },
    { id: 'links' as TabType, label: '链接管理', icon: Link2 },
    { id: 'tags' as TabType, label: '标签管理', icon: Tag },
    { id: 'users' as TabType, label: '用户管理', icon: UsersIcon },
    { id: 'tokens' as TabType, label: 'Token 管理', icon: Key },
    { id: 'settings' as TabType, label: '系统设置', icon: SettingsIcon },
  ]

  const renderContent = () => {
    switch (activeTab) {
      case 'dashboard':
        return <Dashboard />
      case 'categories':
        return <Categories />
      case 'links':
        return <Links />
      case 'tags':
        return <Tags />
      case 'users':
        return <Users />
      case 'tokens':
        return <Tokens />
      case 'settings':
        return <Settings />
      default:
        return <Dashboard />
    }
  }

  return (
    <div className="min-h-screen bg-background">
    <div className="container mx-auto px-4 py-8">
        <div className="flex gap-6">
          {/* 侧边栏 */}
          <div className="w-64 flex-shrink-0">
      <Card>
              <CardContent className="p-4">
                <h2 className="text-xl font-bold mb-4">管理后台</h2>
                <nav className="space-y-2">
                  {tabs.map((tab) => {
                    const Icon = tab.icon
                    return (
                      <Button
                        key={tab.id}
                        variant={activeTab === tab.id ? 'default' : 'ghost'}
                        className="w-full justify-start"
                        onClick={() => setActiveTab(tab.id)}
                      >
                        <Icon className="h-4 w-4 mr-2" />
                        {tab.label}
                      </Button>
                    )
                  })}
                </nav>
        </CardContent>
      </Card>
          </div>

          {/* 主内容区 */}
          <div className="flex-1">{renderContent()}</div>
        </div>
      </div>
    </div>
  )
}
