#!/usr/bin/env ruby
# frozen_string_literal: true

require 'net/http'
require 'json'
require 'uri'

# ============================================================
# 配置区域 - 修改这里的配置
# ============================================================

CONFIG = {
  # API 配置
  api_url: 'https://ops-nav.devops.com/api/v1',
  token: 'kk_14b4d1a80a84f17d31954e9ba1877620398faa20d0948d199d541fffbbc7ed11',

  # 分类和链接配置
  data: [
    {
      category: {
        name: '基础架构',
        icon: '🏗️',
        color: '#f97316',
        description: '基础设施和运维平台',
        sort_order: 1
      },
      links: [
        {
          title: 'JumpServer',
          url: 'https://jumpserver.devops.com/',
          description: '开源堡垒机，统一运维审计平台'
        }
      ]
    },
    {
      category: {
        name: '开发工具',
        icon: '💻',
        color: '#3b82f6',
        description: '常用的开发工具和平台',
        sort_order: 2
      },
      links: [
        {
          title: 'GitHub',
          url: 'https://github.com',
          description: '全球最大的代码托管平台'
        },
        {
          title: 'GitLab',
          url: 'https://gitlab.com',
          description: '完整的 DevOps 生命周期工具'
        },
        {
          title: 'VS Code',
          url: 'https://code.visualstudio.com',
          description: '微软开源的轻量级代码编辑器'
        }
      ]
    },
    {
      category: {
        name: '云服务',
        icon: '☁️',
        color: '#10b981',
        description: '云计算平台和服务',
        sort_order: 3
      },
      links: [
        {
          title: 'AWS Console',
          url: 'https://console.aws.amazon.com',
          description: '亚马逊云服务管理控制台'
        },
        {
          title: '阿里云',
          url: 'https://www.aliyun.com',
          description: '阿里巴巴云计算平台'
        },
        {
          title: '腾讯云',
          url: 'https://cloud.tencent.com',
          description: '腾讯云计算平台'
        }
      ]
    },
    {
      category: {
        name: '监控工具',
        icon: '📊',
        color: '#f59e0b',
        description: '系统监控和性能分析',
        sort_order: 4
      },
      links: [
        {
          title: 'Grafana',
          url: 'https://grafana.com',
          description: '开源的数据可视化和监控平台'
        },
        {
          title: 'Prometheus',
          url: 'https://prometheus.io',
          description: '开源的监控和告警系统'
        },
        {
          title: 'Zabbix',
          url: 'https://www.zabbix.com',
          description: '企业级开源监控解决方案'
        }
      ]
    },
    {
      category: {
        name: '容器化',
        icon: '🐳',
        color: '#8b5cf6',
        description: 'Docker 和 Kubernetes 相关',
        sort_order: 5
      },
      links: [
        {
          title: 'Docker Hub',
          url: 'https://hub.docker.com',
          description: 'Docker 官方镜像仓库'
        },
        {
          title: 'Kubernetes',
          url: 'https://kubernetes.io',
          description: '容器编排和管理平台'
        },
        {
          title: 'Portainer',
          url: 'https://www.portainer.io',
          description: 'Docker 和 Kubernetes 可视化管理工具'
        }
      ]
    },
    {
      category: {
        name: '数据库',
        icon: '🗄️',
        color: '#ef4444',
        description: '数据库管理和工具',
        sort_order: 6
      },
      links: [
        {
          title: 'MySQL',
          url: 'https://www.mysql.com',
          description: '最流行的开源关系型数据库'
        },
        {
          title: 'PostgreSQL',
          url: 'https://www.postgresql.org',
          description: '强大的开源对象关系型数据库'
        },
        {
          title: 'Redis',
          url: 'https://redis.io',
          description: '高性能的内存数据库和缓存'
        }
      ]
    }
  ]
}.freeze

# ============================================================
# HTTP 请求类
# ============================================================

class ApiClient
  def initialize(base_url, token)
    @base_url = base_url
    @token = token
  end

  def get(path)
    request(:get, path)
  end

  def post(path, body)
    request(:post, path, body)
  end

  private

  def request(method, path, body = nil)
    # 确保路径正确拼接
    full_url = @base_url + path
    uri = URI.parse(full_url)
    http = Net::HTTP.new(uri.host, uri.port)
    http.use_ssl = uri.scheme == 'https'
    http.read_timeout = 30
    http.verify_mode = OpenSSL::SSL::VERIFY_NONE if uri.scheme == 'https'

    request = case method
              when :get
                Net::HTTP::Get.new(uri.request_uri)
              when :post
                Net::HTTP::Post.new(uri.request_uri)
              end

    request['Authorization'] = "Bearer #{@token}"
    request['Content-Type'] = 'application/json'
    request.body = body.to_json if body

    response = http.request(request)

    # 调试输出
    puts "    调试: HTTP #{response.code}, Body: #{response.body[0..200]}" if response.code != '200'

    JSON.parse(response.body)
  rescue JSON::ParserError => e
    { 'code' => -1, 'message' => "JSON 解析错误: #{e.message}, Response: #{response.body[0..200]}" }
  rescue StandardError => e
    { 'code' => -1, 'message' => e.message }
  end
end

# ============================================================
# 导航设置类
# ============================================================

class NavigationSetup
  def initialize(config)
    @config = config
    @client = ApiClient.new(config[:api_url], config[:token])
    @stats = {
      categories_created: 0,
      categories_existed: 0,
      links_created: 0,
      links_failed: 0
    }
  end

  def run
    puts '=' * 60
    puts '网站导航自动配置脚本'
    puts '=' * 60
    puts
    puts "目标网站: #{@config[:api_url]}"
    puts "配置数量: #{@config[:data].length} 个分类"
    puts
    puts '开始配置...'
    puts

    @config[:data].each_with_index do |item, index|
      puts "[#{index + 1}/#{@config[:data].length}] 处理分类: #{item[:category][:name]}"
      process_category(item)
      puts
    end

    print_summary
  end

  private

  def process_category(item)
    category_id = find_or_create_category(item[:category])
    return unless category_id

    puts "  添加 #{item[:links].length} 个链接..."
    item[:links].each_with_index do |link, index|
      add_link(link, category_id, index + 1)
    end
  end

  def find_or_create_category(category_data)
    # 获取现有分类（使用公开 API）
    result = @client.get('/categories')

    if result['code'] != 0
      puts "  ❌ 获取分类列表失败: #{result['message']}"
      return nil
    end

    # 处理不同的返回格式
    categories = if result['data'].is_a?(Hash)
                   result['data']['categories'] || []
                 elsif result['data'].is_a?(Array)
                   result['data']
                 else
                   []
                 end
    existing = categories.find { |c| c['name'] == category_data[:name] }

    if existing
      puts "  ✓ 分类已存在 (ID: #{existing['id']})"
      @stats[:categories_existed] += 1
      return existing['id']
    end

    # 创建新分类
    result = @client.post('/admin/categories', {
                            name: category_data[:name],
                            icon: category_data[:icon],
                            color: category_data[:color],
                            description: category_data[:description],
                            sort_order: category_data[:sort_order],
                            active: true
                          })

    if result['code'].zero?
      puts "  ✓ 分类创建成功 (ID: #{result['data']['id']})"
      @stats[:categories_created] += 1
      result['data']['id']
    else
      puts "  ❌ 创建分类失败: #{result['message']}"
      nil
    end
  end

  def add_link(link_data, category_id, sort_order)
    result = @client.post('/admin/links', {
                            title: link_data[:title],
                            url: link_data[:url],
                            description: link_data[:description],
                            category_id: category_id,
                            sort_order: sort_order,
                            status: 'active'
                          })

    if result['code'].zero?
      puts "    ✓ #{link_data[:title]}"
      @stats[:links_created] += 1
    else
      puts "    ❌ #{link_data[:title]}: #{result['message']}"
      @stats[:links_failed] += 1
    end
  end

  def print_summary
    puts '=' * 60
    puts '配置完成'
    puts '=' * 60
    puts
    puts '统计信息:'
    puts "  分类创建: #{@stats[:categories_created]}"
    puts "  分类已存在: #{@stats[:categories_existed]}"
    puts "  链接创建成功: #{@stats[:links_created]}"
    puts "  链接创建失败: #{@stats[:links_failed]}" if @stats[:links_failed].positive?
    puts
    puts "访问 #{@config[:api_url].gsub('/api/v1', '')} 查看结果"
    puts
  end
end

# ============================================================
# 主程序
# ============================================================

if __FILE__ == $PROGRAM_NAME
  begin
    setup = NavigationSetup.new(CONFIG)
    setup.run
  rescue Interrupt
    puts "\n\n操作已取消"
    exit 1
  rescue StandardError => e
    puts "\n❌ 错误: #{e.message}"
    puts e.backtrace.first(5).join("\n")
    exit 1
  end
end
