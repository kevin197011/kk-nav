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
import { Plus, Edit, Trash2, Copy, Eye, EyeOff } from 'lucide-react'

interface APIToken {
  id: number
  name: string
  token: string
  user_id: number
  user?: {
    id: number
    username: string
    email: string
  }
  last_used_at?: string
  expires_at?: string
  active: boolean
  created_at: string
  updated_at: string
}

export default function Tokens() {
  const [tokens, setTokens] = useState<APIToken[]>([])
  const [users, setUsers] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [showToken, setShowToken] = useState<Record<number, boolean>>({})
  const [editingToken, setEditingToken] = useState<APIToken | null>(null)
  const [formData, setFormData] = useState({
    name: '',
    user_id: 0,
    expires_at: '',
    active: true,
  })

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    try {
      setLoading(true)
      const [tokensRes, usersRes]: any[] = await Promise.all([
        api.get('/admin/tokens'),
        api.get('/admin/users'),
      ])
      if (tokensRes.code === 0 && tokensRes.data) {
        setTokens(tokensRes.data.tokens || [])
      }
      if (usersRes.code === 0 && usersRes.data) {
        setUsers(usersRes.data.users || [])
      }
    } catch (error) {
      console.error('Failed to load data:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleOpenDialog = (token?: APIToken) => {
    if (token) {
      setEditingToken(token)
      setFormData({
        name: token.name,
        user_id: token.user_id,
        expires_at: token.expires_at ? token.expires_at.split('T')[0] : '',
        active: token.active,
      })
    } else {
      setEditingToken(null)
      setFormData({
        name: '',
        user_id: users[0]?.id || 0,
        expires_at: '',
        active: true,
      })
    }
    setDialogOpen(true)
  }

  const handleSave = async () => {
    try {
      const data: any = {
        name: formData.name,
        user_id: formData.user_id,
        active: formData.active,
      }
      if (formData.expires_at) {
        data.expires_at = new Date(formData.expires_at).toISOString()
      }

      if (editingToken) {
        const response: any = await api.put(`/admin/tokens/${editingToken.id}`, data)
        if (response.code === 0) {
          await loadData()
          setDialogOpen(false)
        } else {
          alert(response.message || '更新失败')
        }
      } else {
        const response: any = await api.post('/admin/tokens', data)
        if (response.code === 0) {
          await loadData()
          setDialogOpen(false)
          // 如果是新建，显示 token
          if (response.data?.token) {
            setShowToken({ [response.data.id]: true })
          }
        } else {
          alert(response.message || '创建失败')
        }
      }
    } catch (error: any) {
      alert(error.message || '操作失败')
    }
  }

  const handleDelete = async (id: number) => {
    if (!confirm('确定要删除这个 Token 吗？删除后无法恢复！')) return
    try {
      const response: any = await api.delete(`/admin/tokens/${id}`)
      if (response.code === 0) {
        await loadData()
      } else {
        alert(response.message || '删除失败')
      }
    } catch (error: any) {
      alert(error.message || '删除失败')
    }
  }

  const handleCopyToken = (token: string) => {
    navigator.clipboard.writeText(token)
    alert('Token 已复制到剪贴板')
  }

  const formatDate = (dateString?: string) => {
    if (!dateString) return '-'
    return new Date(dateString).toLocaleString('zh-CN')
  }

  const isExpired = (token: APIToken) => {
    if (!token.expires_at) return false
    return new Date(token.expires_at) < new Date()
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
          <h2 className="text-3xl font-bold">Token 管理</h2>
          <p className="text-muted-foreground">管理 API 访问令牌</p>
        </div>
        <Button onClick={() => handleOpenDialog()}>
          <Plus className="h-4 w-4 mr-2" />
          新建 Token
        </Button>
      </div>

      <Card>
        <CardContent className="p-0">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b">
                  <th className="p-4 text-left">名称</th>
                  <th className="p-4 text-left">Token</th>
                  <th className="p-4 text-left">用户</th>
                  <th className="p-4 text-left">状态</th>
                  <th className="p-4 text-left">过期时间</th>
                  <th className="p-4 text-left">最后使用</th>
                  <th className="p-4 text-right">操作</th>
                </tr>
              </thead>
              <tbody>
                {tokens.map((token) => (
                  <tr key={token.id} className="border-b">
                    <td className="p-4 font-medium">{token.name}</td>
                    <td className="p-4">
                      <div className="flex items-center gap-2">
                        <code className="text-sm bg-muted px-2 py-1 rounded">
                          {showToken[token.id]
                            ? token.token
                            : token.token.substring(0, 20) + '...'}
                        </code>
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() =>
                            setShowToken({
                              ...showToken,
                              [token.id]: !showToken[token.id],
                            })
                          }
                        >
                          {showToken[token.id] ? (
                            <EyeOff className="h-4 w-4" />
                          ) : (
                            <Eye className="h-4 w-4" />
                          )}
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => handleCopyToken(token.token)}
                        >
                          <Copy className="h-4 w-4" />
                        </Button>
                      </div>
                    </td>
                    <td className="p-4">
                      {token.user ? (
                        <div>
                          <div className="font-medium">{token.user.username}</div>
                          <div className="text-sm text-muted-foreground">
                            {token.user.email}
                          </div>
                        </div>
                      ) : (
                        '-'
                      )}
                    </td>
                    <td className="p-4">
                      <span
                        className={`px-2 py-1 rounded text-xs ${
                          !token.active
                            ? 'bg-gray-100 text-gray-800'
                            : isExpired(token)
                            ? 'bg-red-100 text-red-800'
                            : 'bg-green-100 text-green-800'
                        }`}
                      >
                        {!token.active
                          ? '已禁用'
                          : isExpired(token)
                          ? '已过期'
                          : '正常'}
                      </span>
                    </td>
                    <td className="p-4 text-sm">{formatDate(token.expires_at)}</td>
                    <td className="p-4 text-sm">{formatDate(token.last_used_at)}</td>
                    <td className="p-4">
                      <div className="flex items-center justify-end gap-2">
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => handleOpenDialog(token)}
                        >
                          <Edit className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => handleDelete(token.id)}
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
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle className="text-2xl">
              {editingToken ? '编辑 Token' : '新建 Token'}
            </DialogTitle>
          </DialogHeader>
          <div className="space-y-5 py-6">
            <div>
              <Label htmlFor="name" className="text-base font-medium mb-2 block">
                名称 *
              </Label>
              <Input
                id="name"
                className="h-11 text-base"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                placeholder="例如：API 客户端 Token"
              />
            </div>
            <div>
              <Label htmlFor="user_id" className="text-base font-medium mb-2 block">
                用户 *
              </Label>
              <select
                id="user_id"
                className="flex h-11 w-full rounded-md border border-input bg-background px-4 py-2 text-base"
                value={formData.user_id}
                onChange={(e) =>
                  setFormData({ ...formData, user_id: parseInt(e.target.value) })
                }
                disabled={!!editingToken}
              >
                <option value={0}>选择用户</option>
                {users.map((user) => (
                  <option key={user.id} value={user.id}>
                    {user.username} ({user.email})
                  </option>
                ))}
              </select>
            </div>
            <div>
              <Label htmlFor="expires_at" className="text-base font-medium mb-2 block">
                过期时间（可选）
              </Label>
              <Input
                id="expires_at"
                type="date"
                className="h-11 text-base"
                value={formData.expires_at}
                onChange={(e) => setFormData({ ...formData, expires_at: e.target.value })}
              />
              <p className="text-sm text-muted-foreground mt-1">
                留空表示永不过期
              </p>
            </div>
            <div className="flex items-center gap-2">
              <input
                type="checkbox"
                id="active"
                checked={formData.active}
                onChange={(e) => setFormData({ ...formData, active: e.target.checked })}
              />
              <Label htmlFor="active" className="text-base font-medium">
                启用
              </Label>
            </div>
            {!editingToken && (
              <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
                <p className="text-sm text-yellow-800">
                  <strong>注意：</strong>Token 创建后只会显示一次，请妥善保管。如果丢失，需要删除后重新创建。
                </p>
              </div>
            )}
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

