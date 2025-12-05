// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useEffect, useState } from 'react'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
  DialogClose,
} from '@/components/ui/dialog'
import api from '@/lib/api'
import type { Link, Category } from '@/types'
import { Plus, Edit, Trash2, ExternalLink } from 'lucide-react'

export default function Links() {
  const [links, setLinks] = useState<Link[]>([])
  const [categories, setCategories] = useState<Category[]>([])
  const [loading, setLoading] = useState(true)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [editingLink, setEditingLink] = useState<Link | null>(null)
  const [formData, setFormData] = useState({
    title: '',
    url: '',
    description: '',
    category_id: 0,
    sort_order: 0,
    status: 'active' as 'active' | 'inactive' | 'error',
    tag_names: [] as string[],
  })

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    try {
      setLoading(true)
      const [linksRes, categoriesRes]: any[] = await Promise.all([
        api.get('/admin/links'),
        api.get('/admin/categories'),
      ])
      if (linksRes.code === 0 && linksRes.data) {
        setLinks(linksRes.data.links || [])
      }
      if (categoriesRes.code === 0 && categoriesRes.data) {
        setCategories(categoriesRes.data.categories || [])
      }
    } catch (error) {
      console.error('Failed to load data:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleOpenDialog = (link?: Link) => {
    if (link) {
      setEditingLink(link)
      setFormData({
        title: link.title,
        url: link.url,
        description: link.description || '',
        category_id: link.category_id,
        sort_order: 0,
        status: link.status,
        tag_names: link.tags?.map((t) => t.name) || [],
      })
    } else {
      setEditingLink(null)
      setFormData({
        title: '',
        url: '',
        description: '',
        category_id: categories[0]?.id || 0,
        sort_order: 0,
        status: 'active',
        tag_names: [],
      })
    }
    setDialogOpen(true)
  }

  const handleSave = async () => {
    try {
      if (editingLink) {
        const response: any = await api.put(`/admin/links/${editingLink.id}`, formData)
        if (response.code === 0) {
          await loadData()
          setDialogOpen(false)
        } else {
          alert(response.message || '更新失败')
        }
      } else {
        const response: any = await api.post('/admin/links', formData)
        if (response.code === 0) {
          await loadData()
          setDialogOpen(false)
        } else {
          alert(response.message || '创建失败')
        }
      }
    } catch (error: any) {
      alert(error.message || '操作失败')
    }
  }

  const handleDelete = async (id: number) => {
    if (!confirm('确定要删除这个链接吗？')) return
    try {
      const response: any = await api.delete(`/admin/links/${id}`)
      if (response.code === 0) {
        await loadData()
      } else {
        alert(response.message || '删除失败')
      }
    } catch (error: any) {
      alert(error.message || '删除失败')
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

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-3xl font-bold">链接管理</h2>
          <p className="text-muted-foreground">管理网站链接</p>
        </div>
        <Button onClick={() => handleOpenDialog()}>
          <Plus className="h-4 w-4 mr-2" />
          新建链接
        </Button>
      </div>

      <Card>
        <CardContent className="p-0">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b">
                  <th className="p-4 text-left">标题</th>
                  <th className="p-4 text-left">URL</th>
                  <th className="p-4 text-left">分类</th>
                  <th className="p-4 text-left">状态</th>
                  <th className="p-4 text-left">点击</th>
                  <th className="p-4 text-right">操作</th>
                </tr>
              </thead>
              <tbody>
                {links.map((link) => (
                  <tr key={link.id} className="border-b">
                    <td className="p-4 font-medium">{link.title}</td>
                    <td className="p-4">
                      <a
                        href={link.url}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-primary hover:underline flex items-center gap-1"
                      >
                        {link.url}
                        <ExternalLink className="h-3 w-3" />
                      </a>
                    </td>
                    <td className="p-4">{link.category?.name || '-'}</td>
                    <td className="p-4">
                      <span
                        className={`px-2 py-1 rounded text-xs ${
                          link.status === 'active'
                            ? 'bg-green-100 text-green-800'
                            : link.status === 'error'
                            ? 'bg-red-100 text-red-800'
                            : 'bg-gray-100 text-gray-800'
                        }`}
                      >
                        {link.status === 'active'
                          ? '正常'
                          : link.status === 'error'
                          ? '错误'
                          : '禁用'}
                      </span>
                    </td>
                    <td className="p-4">{link.click_count || 0}</td>
                    <td className="p-4">
                      <div className="flex items-center justify-end gap-2">
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => handleOpenDialog(link)}
                        >
                          <Edit className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => handleDelete(link.id)}
                        >
                          <Trash2 className="h-4 w-4 text-destructive" />
                        </Button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>

      <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
        <DialogContent className="max-w-3xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle className="text-2xl">{editingLink ? '编辑链接' : '新建链接'}</DialogTitle>
          </DialogHeader>
          <div className="space-y-5 py-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
              <div className="md:col-span-2">
                <Label htmlFor="title" className="text-base font-medium mb-2 block">
                  标题 *
                </Label>
                <Input
                  id="title"
                  className="h-11 text-base"
                  value={formData.title}
                  onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                  placeholder="请输入链接标题"
                />
              </div>
              <div className="md:col-span-2">
                <Label htmlFor="url" className="text-base font-medium mb-2 block">
                  URL *
                </Label>
                <Input
                  id="url"
                  type="url"
                  className="h-11 text-base"
                  value={formData.url}
                  onChange={(e) => setFormData({ ...formData, url: e.target.value })}
                  placeholder="https://example.com"
                />
              </div>
              <div className="md:col-span-2">
                <Label htmlFor="description" className="text-base font-medium mb-2 block">
                  描述
                </Label>
                <Textarea
                  id="description"
                  className="min-h-[100px] text-base"
                  value={formData.description}
                  onChange={(e) =>
                    setFormData({ ...formData, description: e.target.value })
                  }
                  placeholder="请输入链接描述（可选）"
                />
              </div>
              <div>
                <Label htmlFor="category_id" className="text-base font-medium mb-2 block">
                  分类 *
                </Label>
                <select
                  id="category_id"
                  className="flex h-11 w-full rounded-md border border-input bg-background px-4 py-2 text-base"
                  value={formData.category_id}
                  onChange={(e) =>
                    setFormData({ ...formData, category_id: parseInt(e.target.value) })
                  }
                >
                  <option value={0}>选择分类</option>
                  {categories.map((cat) => (
                    <option key={cat.id} value={cat.id}>
                      {cat.name}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <Label htmlFor="status" className="text-base font-medium mb-2 block">
                  状态
                </Label>
                <select
                  id="status"
                  className="flex h-11 w-full rounded-md border border-input bg-background px-4 py-2 text-base"
                  value={formData.status}
                  onChange={(e) =>
                    setFormData({
                      ...formData,
                      status: e.target.value as 'active' | 'inactive' | 'error',
                    })
                  }
                >
                  <option value="active">正常</option>
                  <option value="inactive">禁用</option>
                  <option value="error">错误</option>
                </select>
              </div>
              <div className="md:col-span-2">
                <Label htmlFor="tag_names" className="text-base font-medium mb-2 block">
                  标签（用逗号分隔）
                </Label>
                <Input
                  id="tag_names"
                  className="h-11 text-base"
                  value={formData.tag_names.join(', ')}
                  onChange={(e) =>
                    setFormData({
                      ...formData,
                      tag_names: e.target.value.split(',').map((s) => s.trim()),
                    })
                  }
                  placeholder="监控, 开源, 可视化"
                />
                <p className="text-sm text-muted-foreground mt-1">
                  多个标签请用逗号分隔，例如：监控, 开源, 可视化
                </p>
              </div>
            </div>
          </div>
          <DialogFooter className="mt-6">
            <Button
              variant="outline"
              onClick={() => setDialogOpen(false)}
              className="h-11 px-6"
            >
              取消
            </Button>
            <Button onClick={handleSave} className="h-11 px-6">
              保存
            </Button>
          </DialogFooter>
          <DialogClose onClose={() => setDialogOpen(false)} />
        </DialogContent>
      </Dialog>
    </div>
  )
}

