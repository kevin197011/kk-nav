// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useEffect, useState } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import api from '@/lib/api'
import { Link2, Folder, Tag, Users, MousePointerClick } from 'lucide-react'

interface DashboardStats {
  total_links: number
  active_links: number
  inactive_links: number
  error_links: number
  total_categories: number
  total_tags: number
  total_users: number
  total_clicks: number
  today_clicks: number
  this_week_clicks: number
  this_month_clicks: number
}

export default function Dashboard() {
  const [stats, setStats] = useState<DashboardStats | null>(null)
  const [popularLinks, setPopularLinks] = useState<any[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadDashboard()
  }, [])

  const loadDashboard = async () => {
    try {
      setLoading(true)
      const response: any = await api.get('/admin/dashboard')
      if (response.code === 0 && response.data) {
        setStats(response.data.stats)
        setPopularLinks(response.data.popular_links || [])
      }
    } catch (error) {
      console.error('Failed to load dashboard:', error)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto"></div>
          <p className="mt-4 text-muted-foreground">加载中...</p>
        </div>
      </div>
    )
  }

  if (!stats) {
    return <div className="text-center text-muted-foreground">暂无数据</div>
  }

  const statCards = [
    {
      title: '总链接数',
      value: stats.total_links,
      icon: Link2,
      color: 'text-blue-500',
      bgColor: 'bg-blue-50',
    },
    {
      title: '活跃链接',
      value: stats.active_links,
      icon: Link2,
      color: 'text-green-500',
      bgColor: 'bg-green-50',
    },
    {
      title: '分类数',
      value: stats.total_categories,
      icon: Folder,
      color: 'text-purple-500',
      bgColor: 'bg-purple-50',
    },
    {
      title: '标签数',
      value: stats.total_tags,
      icon: Tag,
      color: 'text-orange-500',
      bgColor: 'bg-orange-50',
    },
    {
      title: '用户数',
      value: stats.total_users,
      icon: Users,
      color: 'text-pink-500',
      bgColor: 'bg-pink-50',
    },
    {
      title: '总点击',
      value: stats.total_clicks,
      icon: MousePointerClick,
      color: 'text-indigo-500',
      bgColor: 'bg-indigo-50',
    },
  ]

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-3xl font-bold">仪表盘</h2>
        <p className="text-muted-foreground">系统概览和统计数据</p>
      </div>

      {/* 统计卡片 */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {statCards.map((stat, index) => {
          const Icon = stat.icon
          return (
            <Card key={index}>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">{stat.title}</CardTitle>
                <div className={`${stat.bgColor} p-2 rounded-lg`}>
                  <Icon className={`h-4 w-4 ${stat.color}`} />
                </div>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stat.value}</div>
              </CardContent>
            </Card>
          )
        })}
      </div>

      {/* 点击统计 */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card>
          <CardHeader>
            <CardTitle className="text-sm font-medium">今日点击</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.today_clicks}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle className="text-sm font-medium">本周点击</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.this_week_clicks}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle className="text-sm font-medium">本月点击</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.this_month_clicks}</div>
          </CardContent>
        </Card>
      </div>

      {/* 热门链接 */}
      {popularLinks.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>热门链接</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              {popularLinks.map((link) => (
                <div
                  key={link.id}
                  className="flex items-center justify-between p-2 border rounded"
                >
                  <span className="font-medium">{link.title}</span>
                  <span className="text-sm text-muted-foreground">
                    {link.click_count} 次点击
                  </span>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  )
}

