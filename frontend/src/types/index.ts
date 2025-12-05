// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

export interface User {
  id: number
  email: string
  username: string
  role: 'user' | 'admin'
  created_at: string
  updated_at: string
}

export interface Category {
  id: number
  name: string
  description?: string
  icon?: string
  color?: string
  sort_order: number
  active: boolean
  created_at: string
  updated_at: string
}

export interface Tag {
  id: number
  name: string
  color?: string
  created_at: string
  updated_at: string
}

export interface Link {
  id: number
  title: string
  url: string
  description?: string
  icon?: string
  status: 'active' | 'inactive' | 'error'
  click_count: number
  category_id: number
  category?: Category
  tags?: Tag[]
  created_at: string
  updated_at: string
}

export interface Stats {
  total_links: number
  total_categories: number
  total_clicks: number
  today_clicks: number
}

export interface ApiResponse<T = any> {
  code: number
  message: string
  data?: T
}

