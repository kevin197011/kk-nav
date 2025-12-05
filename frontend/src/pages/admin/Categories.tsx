// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useEffect, useState } from 'react'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
  DialogClose,
} from '@/components/ui/dialog'
import api from '@/lib/api'
import type { Category } from '@/types'
import { Plus, Edit, Trash2, ArrowUp, ArrowDown } from 'lucide-react'

export default function Categories() {
  const [categories, setCategories] = useState<Category[]>([])
  const [loading, setLoading] = useState(true)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [editingCategory, setEditingCategory] = useState<Category | null>(null)
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    icon: '',
    color: '#007bff',
    sort_order: 0,
    active: true,
  })

  useEffect(() => {
    loadCategories()
  }, [])

  const loadCategories = async () => {
    try {
      setLoading(true)
      const response: any = await api.get('/admin/categories')
      if (response.code === 0 && response.data) {
        setCategories(response.data.categories || [])
      }
    } catch (error) {
      console.error('Failed to load categories:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleOpenDialog = (category?: Category) => {
    if (category) {
      setEditingCategory(category)
      setFormData({
        name: category.name,
        description: category.description || '',
        icon: category.icon || '',
        color: category.color || '#007bff',
        sort_order: category.sort_order,
        active: category.active,
      })
    } else {
      setEditingCategory(null)
      setFormData({
        name: '',
        description: '',
        icon: '',
        color: '#007bff',
        sort_order: categories.length + 1,
        active: true,
      })
    }
    setDialogOpen(true)
  }

  const handleSave = async () => {
    try {
      if (editingCategory) {
        const response: any = await api.put(`/admin/categories/${editingCategory.id}`, formData)
        if (response.code === 0) {
          await loadCategories()
          setDialogOpen(false)
        } else {
          alert(response.message || '更新失败')
        }
      } else {
        const response: any = await api.post('/admin/categories', formData)
        if (response.code === 0) {
          await loadCategories()
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
    if (!confirm('确定要删除这个分类吗？')) return
    try {
      const response: any = await api.delete(`/admin/categories/${id}`)
      if (response.code === 0) {
        await loadCategories()
      } else {
        alert(response.message || '删除失败')
      }
    } catch (error: any) {
      alert(error.message || '删除失败')
    }
  }

  const handleMove = async (id: number, direction: 'up' | 'down') => {
    try {
      console.log(`Moving category ${id} ${direction}`)
      const response: any = await api.patch(`/admin/categories/${id}/move-${direction}`)
      console.log('Move response:', response)
      if (response.code === 0) {
        await loadCategories()
      } else {
        alert(response.message || '移动失败')
      }
    } catch (error: any) {
      console.error('Move error:', error)
      alert(error.message || '移动失败')
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
          <h2 className="text-3xl font-bold">分类管理</h2>
          <p className="text-muted-foreground">管理网站分类</p>
        </div>
        <Button onClick={() => handleOpenDialog()}>
          <Plus className="h-4 w-4 mr-2" />
          新建分类
        </Button>
      </div>

      <Card>
        <CardContent className="p-0">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b">
                  <th className="p-4 text-left">排序</th>
                  <th className="p-4 text-left">图标</th>
                  <th className="p-4 text-left">名称</th>
                  <th className="p-4 text-left">描述</th>
                  <th className="p-4 text-left">颜色</th>
                  <th className="p-4 text-left">状态</th>
                  <th className="p-4 text-right">操作</th>
                </tr>
              </thead>
              <tbody>
                {categories.map((category) => (
                  <tr key={category.id} className="border-b">
                    <td className="p-4">
                      <div className="flex items-center gap-2">
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => handleMove(category.id, 'up')}
                        >
                          <ArrowUp className="h-4 w-4" />
                        </Button>
                        <span>{category.sort_order}</span>
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => handleMove(category.id, 'down')}
                        >
                          <ArrowDown className="h-4 w-4" />
                        </Button>
                      </div>
                    </td>
                    <td className="p-4">{category.icon || '📁'}</td>
                    <td className="p-4 font-medium">{category.name}</td>
                    <td className="p-4 text-muted-foreground">
                      {category.description || '-'}
                    </td>
                    <td className="p-4">
                      <div
                        className="w-6 h-6 rounded"
                        style={{ backgroundColor: category.color }}
                      />
                    </td>
                    <td className="p-4">
                      <span
                        className={`px-2 py-1 rounded text-xs ${
                          category.active
                            ? 'bg-green-100 text-green-800'
                            : 'bg-gray-100 text-gray-800'
                        }`}
                      >
                        {category.active ? '启用' : '禁用'}
                      </span>
                    </td>
                    <td className="p-4">
                      <div className="flex items-center justify-end gap-2">
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => handleOpenDialog(category)}
                        >
                          <Edit className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => handleDelete(category.id)}
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
        <DialogContent>
          <DialogHeader>
            <DialogTitle>
              {editingCategory ? '编辑分类' : '新建分类'}
            </DialogTitle>
          </DialogHeader>
          <div className="space-y-4 py-4">
            <div>
              <Label htmlFor="name">名称 *</Label>
              <Input
                id="name"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              />
            </div>
            <div>
              <Label htmlFor="description">描述</Label>
              <Input
                id="description"
                value={formData.description}
                onChange={(e) =>
                  setFormData({ ...formData, description: e.target.value })
                }
              />
            </div>
            <div>
              <Label htmlFor="icon">图标</Label>
              <Input
                id="icon"
                value={formData.icon}
                onChange={(e) => setFormData({ ...formData, icon: e.target.value })}
                placeholder="📁"
              />
            </div>
            <div>
              <Label htmlFor="color">颜色</Label>
              <Input
                id="color"
                type="color"
                value={formData.color}
                onChange={(e) => setFormData({ ...formData, color: e.target.value })}
              />
            </div>
            <div>
              <Label htmlFor="sort_order">排序</Label>
              <Input
                id="sort_order"
                type="number"
                value={formData.sort_order}
                onChange={(e) =>
                  setFormData({ ...formData, sort_order: parseInt(e.target.value) || 0 })
                }
              />
            </div>
            <div className="flex items-center gap-2">
              <input
                type="checkbox"
                id="active"
                checked={formData.active}
                onChange={(e) => setFormData({ ...formData, active: e.target.checked })}
              />
              <Label htmlFor="active">启用</Label>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDialogOpen(false)}>
              取消
            </Button>
            <Button onClick={handleSave}>保存</Button>
          </DialogFooter>
          <DialogClose onClose={() => setDialogOpen(false)} />
        </DialogContent>
      </Dialog>
    </div>
  )
}

