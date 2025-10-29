# 小型CMDB系统

基于Go语言和SQLite开发的小型配置管理数据库系统，专为离线环境设计，提供完整的IT基础设施资产管理功能。

## 功能特性

### 资产管理
- 服务器、网络设备、存储设备等IT资产的全生命周期管理
- 资产详情记录、状态跟踪、变更历史
- 资产查询、搜索和报表生成

### IP地址管理
- 子网管理与IP地址池管理
- IP地址分配、释放和回收
- IP地址冲突检测

### 报表管理
- 资产统计报表
- 设备类型分布报表
- 机房和机柜使用报表
- IP地址使用报表
- 导出报表功能

### 权限管理
- 用户角色和权限控制
- 细粒度的操作权限管理
- 系统日志审计

## 技术架构

### 后端
- **编程语言**: Go 1.19+
- **Web框架**: Gin
- **数据库**: SQLite（支持离线使用）
- **ORM**: GORM

### 前端
- **UI框架**: Bootstrap 5
- **JavaScript库**: jQuery
- **图标**: Font Awesome, Bootstrap Icons

## 安装与部署

### 环境要求
- Go 1.19+
- SQLite 3+

### 安装步骤
1. 克隆仓库
2. 安装依赖
```bash
go mod tidy
```
3. 创建并配置`.env`文件（复制`.env.example`并修改配置）
4. 运行应用
```bash
go run main.go
```

## 目录结构

```
├── config/          # 配置相关
├── data/            # 数据文件目录（SQLite数据库）
├── handlers/        # 请求处理器
├── middleware/      # 中间件
├── models/          # 数据模型
├── repositories/    # 数据访问层
├── services/        # 业务逻辑层
├── static/          # 静态资源文件
│   ├── css/         # 样式文件
│   ├── js/          # JavaScript文件
├── templates/       # HTML模板
│   ├── layouts/     # 布局模板
│   ├── pages/       # 页面模板
├── main.go          # 程序入口
├── go.mod           # Go模块定义
├── .env.example     # 环境变量示例
├── README.md        # 项目说明文档
```

## 数据库设计

系统使用SQLite数据库，主要包含以下核心表：

- 机房表 (idcs)
- 机柜表 (cabinets)
- 设备表 (devices)
- 机柜U位状态表 (cabinet_u_positions)
- 业务分类表 (business_categories)
- 设备变更记录表 (device_changes)
- 子网管理表 (subnets)
- IP地址池表 (ip_addresses)

详细的数据库表结构请参考代码中的模型定义。

## API接口

系统提供RESTful API接口，主要包括：

### 设备管理
- GET /api/devices - 获取设备列表
- GET /api/devices/search - 搜索设备
- GET /api/devices/:id - 获取设备详情
- POST /api/devices - 创建设备
- PUT /api/devices/:id - 更新设备
- DELETE /api/devices/:id - 删除设备

### 机房管理
- GET /api/idcs - 获取机房列表
- GET /api/idcs/:id - 获取机房详情
- POST /api/idcs - 创建机房
- PUT /api/idcs/:id - 更新机房
- DELETE /api/idcs/:id - 删除机房

### IP地址管理
- GET /api/subnets - 获取子网列表
- GET /api/ip-pool - 获取IP地址列表
- PUT /api/ip-pool/:id/allocate - 分配IP地址
- PUT /api/ip-pool/:id/release - 释放IP地址

更多API接口请参考源代码中的路由定义。

## 开发与贡献

### 开发进度
- 基础架构已搭建完成
- 核心功能模块正在开发中

### 贡献指南
欢迎提交Issue和Pull Request来帮助改进系统。

## 许可证

[MIT License](LICENSE)
