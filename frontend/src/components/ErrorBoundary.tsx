// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { Component } from 'react'
import type { ErrorInfo, ReactNode } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'

interface Props {
  children: ReactNode
}

interface State {
  hasError: boolean
  error: Error | null
}

export class ErrorBoundary extends Component<Props, State> {
  public state: State = {
    hasError: false,
    error: null,
  }

  public static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error }
  }

  public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error('Uncaught error:', error, errorInfo)
  }

  public render() {
    if (this.state.hasError) {
      return (
        <div className="min-h-screen flex items-center justify-center bg-background p-4">
          <Card className="w-full max-w-md">
            <CardHeader>
              <CardTitle className="text-destructive">出现错误</CardTitle>
              <CardDescription>应用遇到了一个错误</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {this.state.error && (
                  <div className="p-3 bg-destructive/10 rounded-md">
                    <p className="text-sm font-mono text-destructive">
                      {this.state.error.message}
                    </p>
                  </div>
                )}
                <div className="flex gap-2">
                  <Button
                    onClick={() => {
                      this.setState({ hasError: false, error: null })
                      window.location.reload()
                    }}
                  >
                    重新加载
                  </Button>
                  <Button
                    variant="outline"
                    onClick={() => {
                      window.location.href = '/'
                    }}
                  >
                    返回首页
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      )
    }

    return this.props.children
  }
}

