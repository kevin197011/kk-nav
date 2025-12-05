// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useEffect, useState } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Search, Star, ExternalLink, Mouse } from 'lucide-react'
import type { Category, Link, Tag, Stats } from '@/types'
import api from '@/lib/api'
import { useAuthStore } from '@/stores/authStore'
import OpsBackground from '@/components/OpsBackground'

export default function Home() {
  const [categories, setCategories] = useState<Category[]>([])
  const [links, setLinks] = useState<Link[]>([])
  const [tags, setTags] = useState<Tag[]>([])
  const [stats, setStats] = useState<Stats | null>(null)
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedTag, setSelectedTag] = useState<string | null>(null)
  const [favorites, setFavorites] = useState<number[]>([])
  const [loading, setLoading] = useState(true)
  const { isAuthenticated } = useAuthStore()

  useEffect(() => {
    loadData()
  }, [])

  useEffect(() => {
    if (isAuthenticated) {
      loadFavorites()
    }
  }, [isAuthenticated])

  const loadData = async () => {
    try {
      setLoading(true)
      const [categoriesRes, linksRes, tagsRes, statsRes] = await Promise.allSettled([
        api.get('/categories'),
        api.get('/links'),
        api.get('/tags'),
        api.get('/stats'),
      ])

      // 处理响应，检查 code 字段
      const getData = (result: PromiseSettledResult<any>) => {
        if (result.status === 'fulfilled') {
          const value = result.value as any
          return value?.code === 0 ? value.data : null
        }
        return null
      }

      // 从 API 响应中提取数据
      const categoriesData = getData(categoriesRes)
      const linksData = getData(linksRes)
      const tagsData = getData(tagsRes)
      const statsData = getData(statsRes)

      // API 返回的数据结构：{categories: [...], total: ...} 或 {links: [...], pagination: ...}
      setCategories(
        Array.isArray(categoriesData?.categories)
          ? categoriesData.categories
          : Array.isArray(categoriesData)
          ? categoriesData
          : []
      )
      setLinks(
        Array.isArray(linksData?.links)
          ? linksData.links
          : Array.isArray(linksData)
          ? linksData
          : []
      )
      setTags(
        Array.isArray(tagsData?.tags)
          ? tagsData.tags
          : Array.isArray(tagsData)
          ? tagsData
          : []
      )
      // stats 数据结构：{stats: {...}, popular_links: [...]}
      setStats(statsData?.stats || statsData || null)
    } catch (error) {
      console.error('Failed to load data:', error)
      // 设置默认值，避免页面崩溃
      setCategories([])
      setLinks([])
      setTags([])
      setStats(null)
    } finally {
      setLoading(false)
    }
  }

  const loadFavorites = async () => {
    try {
      const response = await api.get('/favorites') as any
      if (response.code === 0 && response.data?.links && Array.isArray(response.data.links)) {
        setFavorites(response.data.links.map((link: Link) => link.id))
      }
    } catch (error) {
      console.error('Failed to load favorites:', error)
    }
  }

  const toggleFavorite = async (linkId: number) => {
    if (!isAuthenticated) {
      alert('请先登录')
      return
    }

    const isFavorited = Array.isArray(favorites) && favorites.includes(linkId)
    try {
      if (isFavorited) {
        await api.delete(`/links/${linkId}/unfavorite`)
        setFavorites(Array.isArray(favorites) ? favorites.filter((id) => id !== linkId) : [])
      } else {
        await api.post(`/links/${linkId}/favorite`)
        setFavorites([...favorites, linkId])
      }
    } catch (error) {
      console.error('Failed to toggle favorite:', error)
    }
  }

  const handleClick = async (linkId: number, url: string) => {
    try {
      await api.post(`/links/${linkId}/click`)
      window.open(url, '_blank')
    } catch (error) {
      console.error('Failed to record click:', error)
      window.open(url, '_blank')
    }
  }

  const filteredLinks = Array.isArray(links) ? links.filter((link) => {
    if (searchQuery) {
      const query = searchQuery.toLowerCase()
      if (
        !link.title?.toLowerCase().includes(query) &&
        !link.description?.toLowerCase().includes(query) &&
        !link.url?.toLowerCase().includes(query)
      ) {
        return false
      }
    }
    if (selectedTag) {
      return Array.isArray(link.tags) && link.tags.some((tag) => tag.name === selectedTag)
    }
    return true
  }) : []

  const categoryLinksMap = new Map<number, Link[]>()
  if (Array.isArray(filteredLinks)) {
    filteredLinks.forEach((link) => {
      if (link && link.category_id) {
        const categoryId = link.category_id
        if (!categoryLinksMap.has(categoryId)) {
          categoryLinksMap.set(categoryId, [])
        }
        categoryLinksMap.get(categoryId)!.push(link)
      }
    })
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto"></div>
          <p className="mt-4 text-muted-foreground">加载中...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Hero Section */}
      <div className="relative bg-gradient-to-br from-primary to-purple-600 text-white py-12 overflow-hidden">
        {/* 运维风格动态背景 */}
        <OpsBackground />

        {/* 内容层 */}
        <div className="relative z-10 container mx-auto px-4">
          <div className="max-w-4xl">
            <h1 className="text-4xl font-bold mb-4 flex items-center gap-3">
              <span>🛠️</span>
              运维工具导航
            </h1>
            <p className="text-xl mb-8 text-white/90">专业的运维工具网址导航系统</p>

            {/* Stats */}
            {stats && (
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div className="text-center">
                  <div className="text-3xl font-bold">{stats.total_links || 0}</div>
                  <div className="text-sm text-white/80">个工具</div>
                </div>
                <div className="text-center">
                  <div className="text-3xl font-bold">{stats.total_categories || 0}</div>
                  <div className="text-sm text-white/80">个分类</div>
                </div>
                <div className="text-center">
                  <div className="text-3xl font-bold">{stats.total_clicks || 0}</div>
                  <div className="text-sm text-white/80">总点击</div>
                </div>
                <div className="text-center">
                  <div className="text-3xl font-bold">{stats.today_clicks || 0}</div>
                  <div className="text-sm text-white/80">今日点击</div>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      <div className="container mx-auto px-4 py-8">
        {/* Search and Filters */}
        <div className="max-w-2xl mx-auto mb-8">
          <div className="flex gap-2 mb-4">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
              <input
                type="text"
                placeholder="搜索工具名称、描述或网址..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full pl-10 pr-4 py-2 border rounded-md bg-background"
              />
            </div>
          </div>

          {/* Tags */}
          <div className="flex flex-wrap gap-2 justify-center">
            <Button
              variant={selectedTag === null ? 'default' : 'outline'}
              size="sm"
              onClick={() => setSelectedTag(null)}
            >
              全部
            </Button>
            {Array.isArray(tags) && tags.map((tag) => (
              <Button
                key={tag.id}
                variant={selectedTag === tag.name ? 'default' : 'outline'}
                size="sm"
                onClick={() => setSelectedTag(tag.name)}
                style={
                  selectedTag === tag.name && tag.color
                    ? { backgroundColor: tag.color, borderColor: tag.color }
                    : tag.color
                    ? { borderColor: tag.color, color: tag.color }
                    : {}
                }
              >
                {tag.name}
              </Button>
            ))}
          </div>
        </div>

        {/* Categories and Links */}
        {Array.isArray(categories) && categories
          .filter((category) => category && category.id && categoryLinksMap.has(category.id))
          .map((category) => {
            const categoryLinks = categoryLinksMap.get(category.id) || []
            const safeCategoryLinks = Array.isArray(categoryLinks) ? categoryLinks : []
            return (
              <div key={category.id} className="mb-8">
                <Card>
                  <CardHeader
                    className="text-white rounded-t-lg"
                    style={{
                      background: category.color
                        ? `linear-gradient(135deg, ${category.color} 0%, #6f42c1 100%)`
                        : 'linear-gradient(135deg, #007bff 0%, #6f42c1 100%)',
                    }}
                  >
                    <div className="flex justify-between items-center">
                      <CardTitle className="flex items-center gap-2">
                        <span className="text-2xl">{category.icon || '📁'}</span>
                        {category.name}
                      </CardTitle>
                      <CardDescription className="text-white/80">
                        {safeCategoryLinks.length} 个工具
                      </CardDescription>
                    </div>
                  </CardHeader>
                  <CardContent className="p-6">
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                      {safeCategoryLinks.map((link) => (
                        <Card key={link.id} className="hover:shadow-lg transition-shadow">
                          <CardHeader>
                            <div className="flex justify-between items-start">
                              <CardTitle className="text-lg flex items-center gap-2">
                                <span
                                  className={`w-2 h-2 rounded-full ${
                                    link.status === 'active'
                                      ? 'bg-green-500'
                                      : link.status === 'error'
                                      ? 'bg-red-500'
                                      : 'bg-gray-500'
                                  }`}
                                />
                                {link.title}
                              </CardTitle>
                              {isAuthenticated && (
                                <Button
                                  variant="ghost"
                                  size="icon"
                                  onClick={() => toggleFavorite(link.id)}
                                >
                                  <Star
                                    className={`h-4 w-4 ${
                                      Array.isArray(favorites) && favorites.includes(link.id)
                                        ? 'fill-yellow-400 text-yellow-400'
                                        : ''
                                    }`}
                                  />
                                </Button>
                              )}
                            </div>
                            {link.description && (
                              <CardDescription className="line-clamp-2">
                                {link.description}
                              </CardDescription>
                            )}
                          </CardHeader>
                          <CardContent>
                            <div className="flex justify-between items-center">
                              <Button
                                size="sm"
                                onClick={() => handleClick(link.id, link.url)}
                              >
                                <ExternalLink className="h-4 w-4 mr-1" />
                                访问
                              </Button>
                              <div className="flex items-center gap-1 text-sm text-muted-foreground">
                                <Mouse className="h-4 w-4" />
                                {link.click_count || 0}
                              </div>
                            </div>
                            {Array.isArray(link.tags) && link.tags.length > 0 && (
                              <div className="flex flex-wrap gap-1 mt-3">
                                {link.tags.map((tag) => (
                                  <span
                                    key={tag.id}
                                    className="text-xs px-2 py-1 rounded-full bg-secondary"
                                    style={tag.color ? { color: tag.color } : {}}
                                  >
                                    {tag.name}
                                  </span>
                                ))}
                              </div>
                            )}
                          </CardContent>
                        </Card>
                      ))}
                    </div>
                  </CardContent>
                </Card>
              </div>
            )
          })}

        {filteredLinks.length === 0 && (
          <Card>
            <CardContent className="py-12 text-center">
              <p className="text-muted-foreground">没有找到相关工具</p>
              <Button
                variant="outline"
                className="mt-4"
                onClick={() => {
                  setSearchQuery('')
                  setSelectedTag(null)
                }}
              >
                清除筛选
              </Button>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  )
}

