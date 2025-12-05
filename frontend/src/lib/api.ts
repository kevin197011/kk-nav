// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import axios from 'axios'

const API_BASE_URL = '/api/v1'

// 创建 axios 实例
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器：添加 Token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器：处理错误
api.interceptors.response.use(
  (response) => {
    return response.data as any
  },
  (error) => {
    // 网络错误或服务器无响应
    if (!error.response) {
      console.error('Network error:', error.message)
      // 不抛出错误，返回一个默认响应，让调用方处理
      return Promise.resolve({
        code: -1,
        message: '网络错误，请检查后端服务是否运行',
        data: null,
      })
    }

    if (error.response.status === 401) {
      // Token 过期，清除并跳转到登录页
      localStorage.removeItem('token')
      // 只在非登录页面才跳转
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
    }

    // 返回错误响应，而不是 reject
    return Promise.resolve({
      code: error.response.status || -1,
      message: error.response.data?.message || error.message || '请求失败',
      data: null,
    })
  }
)

export default api

