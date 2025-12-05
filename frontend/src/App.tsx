// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useEffect } from 'react'
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import Layout from '@/components/Layout'
import Home from '@/pages/Home'
import Login from '@/pages/Login'
import Admin from '@/pages/Admin'
import { useAuthStore } from '@/stores/authStore'
import { useSiteSettings } from '@/hooks/useSiteSettings'

function App() {
  const { fetchUser, isAuthenticated } = useAuthStore()
  useSiteSettings() // 加载站点设置并更新页面标题

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (token) {
      fetchUser().catch((error) => {
        console.error('Failed to fetch user:', error)
      })
    }
  }, []) // 只在组件挂载时执行一次

  return (
    <BrowserRouter>
      <Layout>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route
            path="/login"
            element={isAuthenticated ? <Navigate to="/" replace /> : <Login />}
          />
          <Route
            path="/admin"
            element={
              isAuthenticated && useAuthStore.getState().user?.role === 'admin' ? (
                <Admin />
              ) : (
                <Navigate to="/" replace />
              )
            }
          />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </Layout>
    </BrowserRouter>
  )
}

export default App
