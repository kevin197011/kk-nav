// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useEffect, useState } from 'react'
import api from '@/lib/api'

interface SiteSettings {
  site_name?: string
  site_description?: string
  primary_color?: string
  theme?: string
}

export function useSiteSettings() {
  const [settings, setSettings] = useState<SiteSettings>({})
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadSettings()
  }, [])

  const loadSettings = async () => {
    try {
      const response: any = await api.get('/settings')
      if (response.code === 0 && response.data) {
        const siteSettings = response.data.settings || {}
        setSettings(siteSettings)

        // 更新页面标题
        if (siteSettings.site_name) {
          document.title = siteSettings.site_name
        }
      }
    } catch (error) {
      console.error('Failed to load site settings:', error)
      // 使用默认标题
      document.title = '运维工具导航'
    } finally {
      setLoading(false)
    }
  }

  return { settings, loading }
}

