// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useEffect, useState } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import api from '@/lib/api'
import { Save } from 'lucide-react'

export default function Settings() {
  const [settings, setSettings] = useState<Record<string, string>>({})
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    loadSettings()
  }, [])

  const loadSettings = async () => {
    try {
      setLoading(true)
      const response: any = await api.get('/admin/settings')
      if (response.code === 0 && response.data) {
        setSettings(response.data.settings || {})
      }
    } catch (error) {
      console.error('Failed to load settings:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleSave = async () => {
    try {
      setSaving(true)
      const response: any = await api.put('/admin/settings', { settings })
      if (response.code === 0) {
        alert('设置保存成功')
      } else {
        alert(response.message || '保存失败')
      }
    } catch (error: any) {
      alert(error.message || '保存失败')
    } finally {
      setSaving(false)
    }
  }

  const handleChange = (key: string, value: string) => {
    setSettings({ ...settings, [key]: value })
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

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-3xl font-bold">系统设置</h2>
          <p className="text-muted-foreground">管理系统配置</p>
        </div>
        <Button onClick={handleSave} disabled={saving}>
          <Save className="h-4 w-4 mr-2" />
          {saving ? '保存中...' : '保存设置'}
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>基本设置</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div>
            <Label htmlFor="site_name">网站名称</Label>
            <Input
              id="site_name"
              value={settings.site_name || ''}
              onChange={(e) => handleChange('site_name', e.target.value)}
            />
          </div>
          <div>
            <Label htmlFor="site_description">网站描述</Label>
            <Input
              id="site_description"
              value={settings.site_description || ''}
              onChange={(e) => handleChange('site_description', e.target.value)}
            />
          </div>
          <div>
            <Label htmlFor="primary_color">主题色</Label>
            <Input
              id="primary_color"
              type="color"
              value={settings.primary_color || '#007bff'}
              onChange={(e) => handleChange('primary_color', e.target.value)}
            />
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>功能设置</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              id="enable_registration"
              checked={settings.enable_registration === 'true'}
              onChange={(e) =>
                handleChange('enable_registration', e.target.checked ? 'true' : 'false')
              }
            />
            <Label htmlFor="enable_registration">允许用户注册</Label>
          </div>
          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              id="enable_link_check"
              checked={settings.enable_link_check === 'true'}
              onChange={(e) =>
                handleChange('enable_link_check', e.target.checked ? 'true' : 'false')
              }
            />
            <Label htmlFor="enable_link_check">启用链接状态检测</Label>
          </div>
          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              id="enable_analytics"
              checked={settings.enable_analytics === 'true'}
              onChange={(e) =>
                handleChange('enable_analytics', e.target.checked ? 'true' : 'false')
              }
            />
            <Label htmlFor="enable_analytics">启用统计分析</Label>
          </div>
          <div>
            <Label htmlFor="check_interval_hours">链接检测间隔（小时）</Label>
            <Input
              id="check_interval_hours"
              type="number"
              value={settings.check_interval_hours || '24'}
              onChange={(e) => handleChange('check_interval_hours', e.target.value)}
            />
          </div>
          <div>
            <Label htmlFor="links_per_page">每页链接数</Label>
            <Input
              id="links_per_page"
              type="number"
              value={settings.links_per_page || '12'}
              onChange={(e) => handleChange('links_per_page', e.target.value)}
            />
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

