// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { Link, useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { useAuthStore } from '@/stores/authStore'
import { LogOut, User, Settings } from 'lucide-react'
import ThemeToggle from '@/components/ThemeToggle'

interface LayoutProps {
  children: React.ReactNode
}

export default function Layout({ children }: LayoutProps) {
  const { user, isAuthenticated, logout } = useAuthStore()
  const navigate = useNavigate()

  const handleLogout = async () => {
    await logout()
    navigate('/')
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Navigation */}
      <nav className="border-b bg-card">
        <div className="container mx-auto px-4">
          <div className="flex items-center justify-between h-16">
            <Link to="/" className="flex items-center gap-2 text-xl font-bold">
              <span>🧭</span>
              运维工具导航
            </Link>

            <div className="flex items-center gap-4">
              <ThemeToggle />
              {isAuthenticated ? (
                <>
                  {user?.role === 'admin' && (
                    <Button variant="ghost" onClick={() => navigate('/admin')}>
                      <Settings className="h-4 w-4 mr-2" />
                      管理后台
                    </Button>
                  )}
                  <div className="flex items-center gap-2">
                    <User className="h-4 w-4 text-muted-foreground" />
                    <span className="text-sm">{user?.username}</span>
                  </div>
                  <Button variant="ghost" size="sm" onClick={handleLogout}>
                    <LogOut className="h-4 w-4 mr-2" />
                    退出
                  </Button>
                </>
              ) : (
                <Button variant="ghost" onClick={() => navigate('/login')}>
                  登录
                </Button>
              )}
            </div>
          </div>
        </div>
      </nav>

      {/* Main Content */}
      <main>{children}</main>

      {/* Footer */}
      <footer className="border-t bg-card mt-12">
        <div className="container mx-auto px-4 py-6 text-center text-sm text-muted-foreground">
          <p>© 2025 运维工具导航系统. All rights reserved.</p>
        </div>
      </footer>
    </div>
  )
}

