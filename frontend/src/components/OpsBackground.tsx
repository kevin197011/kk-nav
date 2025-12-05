// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

import { useEffect, useRef } from 'react'

export default function OpsBackground() {
  const canvasRef = useRef<HTMLCanvasElement>(null)

  useEffect(() => {
    const canvas = canvasRef.current
    if (!canvas) return

    const ctx = canvas.getContext('2d')
    if (!ctx) return

    // 设置画布大小
    const resizeCanvas = () => {
      canvas.width = canvas.offsetWidth
      canvas.height = canvas.offsetHeight
    }
    resizeCanvas()
    window.addEventListener('resize', resizeCanvas)

    // 流动的粒子
    const particles: Array<{
      x: number
      y: number
      vx: number
      vy: number
      size: number
      opacity: number
      life: number
    }> = []

    // 初始化粒子
    for (let i = 0; i < 30; i++) {
      particles.push({
        x: Math.random() * canvas.width,
        y: Math.random() * canvas.height,
        vx: (Math.random() - 0.5) * 0.3,
        vy: (Math.random() - 0.5) * 0.3,
        size: Math.random() * 2 + 1,
        opacity: Math.random() * 0.3 + 0.1,
        life: Math.random(),
      })
    }

    // 流动的线条
    const lines: Array<{
      x1: number
      y1: number
      x2: number
      y2: number
      progress: number
      speed: number
      opacity: number
    }> = []

    // 初始化线条
    for (let i = 0; i < 8; i++) {
      const angle = Math.random() * Math.PI * 2
      const length = Math.random() * 200 + 100
      const x1 = Math.random() * canvas.width
      const y1 = Math.random() * canvas.height
      lines.push({
        x1,
        y1,
        x2: x1 + Math.cos(angle) * length,
        y2: y1 + Math.sin(angle) * length,
        progress: Math.random(),
        speed: Math.random() * 0.01 + 0.005,
        opacity: Math.random() * 0.2 + 0.1,
      })
    }

    // 光晕效果
    const glows: Array<{
      x: number
      y: number
      radius: number
      opacity: number
      pulse: number
    }> = []

    // 初始化光晕
    for (let i = 0; i < 5; i++) {
      glows.push({
        x: Math.random() * canvas.width,
        y: Math.random() * canvas.height,
        radius: Math.random() * 100 + 50,
        opacity: Math.random() * 0.2 + 0.1,
        pulse: Math.random() * Math.PI * 2,
      })
    }

    let animationFrame: number

    const animate = () => {
      // 半透明背景（创造拖尾效果）
      ctx.fillStyle = 'rgba(0, 0, 0, 0.015)'
      ctx.fillRect(0, 0, canvas.width, canvas.height)

      // 绘制光晕（背景层）
      glows.forEach((glow) => {
        glow.pulse += 0.01
        const currentRadius = glow.radius + Math.sin(glow.pulse) * 20
        const currentOpacity = glow.opacity + Math.sin(glow.pulse) * 0.1

        const gradient = ctx.createRadialGradient(glow.x, glow.y, 0, glow.x, glow.y, currentRadius)
        gradient.addColorStop(0, `rgba(255, 255, 255, ${currentOpacity})`)
        gradient.addColorStop(0.5, `rgba(255, 255, 255, ${currentOpacity * 0.5})`)
        gradient.addColorStop(1, 'rgba(255, 255, 255, 0)')

        ctx.fillStyle = gradient
        ctx.beginPath()
        ctx.arc(glow.x, glow.y, currentRadius, 0, Math.PI * 2)
        ctx.fill()
      })

      // 绘制流动的线条
      lines.forEach((line) => {
        line.progress += line.speed
        if (line.progress > 1) {
          line.progress = 0
          // 重新生成线条
          const angle = Math.random() * Math.PI * 2
          const length = Math.random() * 200 + 100
          line.x1 = Math.random() * canvas.width
          line.y1 = Math.random() * canvas.height
          line.x2 = line.x1 + Math.cos(angle) * length
          line.y2 = line.y1 + Math.sin(angle) * length
        }

        // 计算当前线条位置
        const currentX1 = line.x1 + (line.x2 - line.x1) * line.progress
        const currentY1 = line.y1 + (line.y2 - line.y1) * line.progress
        const currentX2 = line.x1 + (line.x2 - line.x1) * (line.progress + 0.3)
        const currentY2 = line.y1 + (line.y2 - line.y1) * (line.progress + 0.3)

        // 绘制渐变线条
        const gradient = ctx.createLinearGradient(currentX1, currentY1, currentX2, currentY2)
        gradient.addColorStop(0, `rgba(255, 255, 255, 0)`)
        gradient.addColorStop(0.5, `rgba(255, 255, 255, ${line.opacity})`)
        gradient.addColorStop(1, `rgba(255, 255, 255, 0)`)

        ctx.strokeStyle = gradient
        ctx.lineWidth = 1.5
        ctx.beginPath()
        ctx.moveTo(currentX1, currentY1)
        ctx.lineTo(currentX2, currentY2)
        ctx.stroke()
      })

      // 更新和绘制粒子
      particles.forEach((particle) => {
        particle.x += particle.vx
        particle.y += particle.vy
        particle.life += 0.01

        // 边界处理（循环）
        if (particle.x < 0) particle.x = canvas.width
        if (particle.x > canvas.width) particle.x = 0
        if (particle.y < 0) particle.y = canvas.height
        if (particle.y > canvas.height) particle.y = 0

        // 生命周期的透明度变化
        const lifeOpacity = Math.sin(particle.life) * 0.5 + 0.5

        // 绘制粒子（带光晕）
        const glowGradient = ctx.createRadialGradient(particle.x, particle.y, 0, particle.x, particle.y, particle.size * 4)
        glowGradient.addColorStop(0, `rgba(255, 255, 255, ${particle.opacity * lifeOpacity})`)
        glowGradient.addColorStop(0.5, `rgba(255, 255, 255, ${particle.opacity * lifeOpacity * 0.5})`)
        glowGradient.addColorStop(1, 'rgba(255, 255, 255, 0)')

        ctx.fillStyle = glowGradient
        ctx.beginPath()
        ctx.arc(particle.x, particle.y, particle.size * 4, 0, Math.PI * 2)
        ctx.fill()

        // 粒子核心
        ctx.fillStyle = `rgba(255, 255, 255, ${particle.opacity * lifeOpacity * 0.8})`
        ctx.beginPath()
        ctx.arc(particle.x, particle.y, particle.size, 0, Math.PI * 2)
        ctx.fill()
      })

      // 绘制粒子间的连接线（近距离连接）
      ctx.strokeStyle = 'rgba(255, 255, 255, 0.1)'
      ctx.lineWidth = 0.5
      for (let i = 0; i < particles.length; i++) {
        for (let j = i + 1; j < particles.length; j++) {
          const dx = particles[i].x - particles[j].x
          const dy = particles[i].y - particles[j].y
          const dist = Math.sqrt(dx * dx + dy * dy)

          if (dist < 150) {
            const opacity = (1 - dist / 150) * 0.2
            ctx.strokeStyle = `rgba(255, 255, 255, ${opacity})`
            ctx.beginPath()
            ctx.moveTo(particles[i].x, particles[i].y)
            ctx.lineTo(particles[j].x, particles[j].y)
            ctx.stroke()
          }
        }
      }

      animationFrame = requestAnimationFrame(animate)
    }

    animate()

    return () => {
      window.removeEventListener('resize', resizeCanvas)
      if (animationFrame) {
        cancelAnimationFrame(animationFrame)
      }
    }
  }, [])

  return (
    <canvas
      ref={canvasRef}
      className="absolute inset-0 w-full h-full"
      style={{ pointerEvents: 'none' }}
    />
  )
}
