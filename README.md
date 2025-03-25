# 实习进度管理系统

这是一个使用Go语言开发的实习进度管理系统，帮助用户跟踪和管理实习申请进度。

## 功能特点

- 用户注册和登录
- 实习申请记录管理
- 申请状态追踪（已投递/笔试中/面试中/已录用）
- 面试和笔试事件提醒
- 申请数据统计

## 技术栈

- 后端框架：Gin + GORM
- 数据库：MySQL + Redis
- 认证：JWT

## 项目结构

```bash
project-root/
├── config/            # 配置文件
├── internal/          # 核心业务代码
│   ├── handler/      # HTTP路由处理
│   ├── router/       # 路由
│   ├── service/      # 业务逻辑层
│   ├── model/        # 数据模型
│   └── middleware/   # 中间件
├── pkg/              # 公共组件
│   ├── database/     # 数据库连接
│   └── utils/        # 工具函数
├── scripts/          # 部署脚本
└── main.go           # 入口文件
```

## 环境要求

- Go 1.21+
- MySQL 5.7+
- Redis 6.0+

## 安装和运行

1. 克隆项目：
   ```bash
   git clone <repository-url>
   cd internship-manager
   ```

2. 安装依赖：
   ```bash
   go mod tidy
   ```

3. 配置数据库：
   - 创建MySQL数据库和表：
     ```bash
     mysql -u root -p < scripts/init.sql
     ```
   - 修改 `configs/config.yaml` 中的数据库配置

4. 运行项目：
   ```bash
   go run main.go
   ```

## API文档

### 用户相关

- POST /api/register - 用户注册
- POST /api/login - 用户登录

### 申请相关

- POST /api/applications - 创建申请记录
- PUT /api/applications/status - 更新申请状态
- PUT /api/applications/event - 更新面试/笔试事件
- GET /api/applications - 获取所有申请记录
- GET /api/applications/statistics - 获取申请统计信息
- GET /api/applications/upcoming-events - 获取即将到来的面试/笔试事件

## 贡献指南

1. Fork 项目
2. 创建特性分支
3. 提交代码
4. 发起 Pull Request

## 许可证

MIT License 