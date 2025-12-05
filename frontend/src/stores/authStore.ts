// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { create } from 'zustand'
import type { User } from '@/types'
import api from '@/lib/api'

interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  login: (username: string, password: string) => Promise<void>
  logout: () => Promise<void>
  fetchUser: () => Promise<void>
}

export const useAuthStore = create<AuthState>()((set) => ({
      user: null,
      token: null,
      isAuthenticated: false,

      login: async (username: string, password: string) => {
        const response = (await api.post('/auth/login', { username, password })) as any
        if (response.code === 0) {
          set({
            user: response.data.user,
            token: response.data.token,
            isAuthenticated: true,
          })
          localStorage.setItem('token', response.data.token)
        } else {
          throw new Error(response.message || '登录失败')
        }
      },

      logout: async () => {
        try {
          await api.post('/auth/logout')
        } catch (error) {
          console.error('Logout error:', error)
        } finally {
          set({
            user: null,
            token: null,
            isAuthenticated: false,
          })
          localStorage.removeItem('token')
        }
      },

      fetchUser: async () => {
        try {
          const response = (await api.get('/auth/me')) as any
          if (response.code === 0) {
            set({
              user: response.data,
              isAuthenticated: true,
            })
          }
        } catch (error) {
          set({
            user: null,
            token: null,
            isAuthenticated: false,
          })
        }
      },
    })
)

