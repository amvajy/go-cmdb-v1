package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 全局数据库连接
var DB *gorm.DB

// 初始化数据库连接
func initDatabase() error {
	// 从环境变量获取数据库连接信息，如果没有则使用默认值
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(127.0.0.1:3306)/cmdb?charset=utf8mb4&parseTime=True&loc=Local"
	}

	// 配置GORM日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // 慢SQL阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略记录未找到错误
			Colorful:                  true,        // 彩色打印
		},
	)

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// 设置连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(10)
	// 设置最大打开连接数
	sqlDB.SetMaxOpenConns(100)
	// 设置连接最大生存时间
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db

	// 自动迁移数据库表结构
	if err := migrateDatabase(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}

// 数据库迁移
func migrateDatabase() error {
	// 这里会自动创建或更新数据库表结构
	// 由于我们还没有定义模型，这里暂时不做实际迁移
	// 稍后会添加模型定义和迁移代码
	return nil
}

// 初始化路由
func setupRoutes(router *gin.Engine) {
	// 设置静态文件目录
	router.Static("/static", "./static")

	// 加载HTML模板
	router.LoadHTMLGlob("./templates/**/*")

	// 前端页面路由
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/dashboard")
	})

	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"pageTitle": "登录 - CMDB系统",
		})
	})

	router.GET("/dashboard", authMiddleware(), func(c *gin.Context) {
		// 确保使用正确的模板名称
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"pageTitle": "系统概览",
			"user": gin.H{
				"name": "管理员",
				"role": "admin",
			},
		})
	})

	router.GET("/asset-management", authMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "asset_management.html", gin.H{
			"pageTitle": "资产管理",
			"user": gin.H{
				"name": "管理员",
				"role": "admin",
			},
		})
	})

	router.GET("/idc-cabinet-management", authMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "idc_cabinet_management.html", gin.H{
			"pageTitle": "机房机柜管理",
			"user": gin.H{
				"name": "管理员",
				"role": "admin",
			},
		})
	})

	router.GET("/ip-management", authMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "ip_management.html", gin.H{
			"pageTitle": "IP地址管理",
			"user": gin.H{
				"name": "管理员",
				"role": "admin",
			},
		})
	})

	router.GET("/report-management", authMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "report_management.html", gin.H{
			"pageTitle": "报表管理",
			"user": gin.H{
				"name": "管理员",
				"role": "admin",
			},
		})
	})

	router.GET("/system-management", authMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "system_management.html", gin.H{
			"pageTitle": "系统管理",
			"user": gin.H{
				"name": "管理员",
				"role": "admin",
			},
		})
	})

	// API路由组
	api := router.Group("/api")
	{
		// 认证相关路由
		auth := api.Group("/auth")
		{
			auth.POST("/login", handleLogin)
			auth.POST("/logout", handleLogout)
			auth.GET("/me", authMiddleware(), handleGetCurrentUser)
		}

		// 仪表盘数据路由
		dashboard := api.Group("/dashboard", authMiddleware())
		{
			dashboard.GET("/data", handleGetDashboardData)
		}

		// 设备管理路由
		devices := api.Group("/devices", authMiddleware())
		{
			devices.GET("/", handleListDevices)
			devices.GET("/search", handleSearchDevices)
			devices.GET("/:id", handleGetDevice)
			devices.POST("/", handleCreateDevice)
			devices.PUT("/:id", handleUpdateDevice)
			devices.DELETE("/:id", handleDeleteDevice)
		}

		// 机房管理路由
		idcs := api.Group("/idcs", authMiddleware())
		{
			idcs.GET("/", handleListIDCs)
			idcs.GET("/:id", handleGetIDC)
			idcs.POST("/", handleCreateIDC)
			idcs.PUT("/:id", handleUpdateIDC)
			idcs.DELETE("/:id", handleDeleteIDC)
		}

		// 机柜管理路由
		cabinets := api.Group("/cabinets", authMiddleware())
		{
			cabinets.GET("/", handleListCabinets)
			cabinets.GET("/:id", handleGetCabinet)
			cabinets.POST("/", handleCreateCabinet)
			cabinets.PUT("/:id", handleUpdateCabinet)
			cabinets.DELETE("/:id", handleDeleteCabinet)
		}

		// 子网管理路由
		subnets := api.Group("/subnets", authMiddleware())
		{
			subnets.GET("/", handleListSubnets)
			subnets.GET("/:id", handleGetSubnet)
			subnets.POST("/", handleCreateSubnet)
			subnets.PUT("/:id", handleUpdateSubnet)
			subnets.DELETE("/:id", handleDeleteSubnet)
		}

		// IP地址管理路由
		ipPool := api.Group("/ip-pool", authMiddleware())
		{
			ipPool.GET("/", handleListIPAddresses)
			ipPool.GET("/:id", handleGetIPAddress)
			ipPool.POST("/", handleCreateIPAddress)
			ipPool.PUT("/:id", handleUpdateIPAddress)
			ipPool.DELETE("/:id", handleDeleteIPAddress)
			ipPool.PUT("/:id/allocate", handleAllocateIPAddress)
			ipPool.PUT("/:id/release", handleReleaseIPAddress)
		}

		// 报表管理路由
		reports := api.Group("/reports", authMiddleware())
		{
			reports.GET("/summary", handleGetAssetSummary)
			reports.GET("/type-distribution", handleGetAssetTypeDistribution)
			reports.GET("/idc-stats", handleGetIDCStatistics)
			reports.GET("/cabinet-usage", handleGetCabinetUsage)
			reports.GET("/ip-usage", handleGetIPUsage)
			reports.GET("/maintenance", handleGetMaintenanceReport)
			reports.GET("/export", handleExportReport)
		}

		// 用户管理路由
		users := api.Group("/users", authMiddleware(), adminMiddleware())
		{
			users.GET("/", handleListUsers)
			users.GET("/:id", handleGetUser)
			users.POST("/", handleCreateUser)
			users.PUT("/:id", handleUpdateUser)
			users.DELETE("/:id", handleDeleteUser)
		}

		// 角色管理路由
		roles := api.Group("/roles", authMiddleware(), adminMiddleware())
		{
			roles.GET("/", handleListRoles)
			roles.GET("/:id", handleGetRole)
			roles.POST("/", handleCreateRole)
			roles.PUT("/:id", handleUpdateRole)
			roles.DELETE("/:id", handleDeleteRole)
		}

		// 系统设置路由
		settings := api.Group("/settings", authMiddleware())
		{
			settings.GET("/", handleGetSettings)
			settings.PUT("/", handleUpdateSettings)
		}

		// 系统日志路由
		logs := api.Group("/logs", authMiddleware(), adminMiddleware())
		{
			logs.GET("/", handleListSystemLogs)
		}
	}
}

// 认证中间件
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 演示模式：跳过严格的token验证
		// 实际项目中应该实现完整的JWT认证逻辑
		
		// 检查是否有token
		token := c.GetHeader("Authorization")
		if token == "" {
			// 在演示模式下，自动设置默认用户信息
			// 不返回401错误，允许用户访问
			log.Println("演示模式：自动登录为管理员用户")
		}

		// 设置用户信息到上下文
		c.Set("user_id", 1) // 假设用户ID为1
		c.Set("username", "admin") // 假设用户名为admin

		c.Next()
	}
}

// 管理员权限中间件
func adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 这里会实现管理员权限检查逻辑
		// 暂时简单实现，稍后会完善
		// 假设用户ID为1的是管理员
		userID, exists := c.Get("user_id")
		if !exists || userID.(int) != 1 {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Forbidden",
				"message": "没有权限执行此操作",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// 认证相关处理函数
func handleLogin(c *gin.Context) {
	// 这里会实现登录逻辑
	// 暂时返回模拟数据
	c.JSON(http.StatusOK, gin.H{
		"access_token": "mock-token-123456",
		"token_type": "Bearer",
		"expires_in": 3600,
		"user": gin.H{
			"id": 1,
			"username": "admin",
			"email": "admin@example.com",
			"name": "管理员",
			"role": "admin",
		},
	})
}

func handleLogout(c *gin.Context) {
	// 这里会实现登出逻辑
	c.JSON(http.StatusOK, gin.H{
		"message": "成功退出登录",
	})
}

func handleGetCurrentUser(c *gin.Context) {
	// 获取当前用户信息
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")

	c.JSON(http.StatusOK, gin.H{
		"id": userID,
		"username": username,
		"email": "admin@example.com",
		"name": "管理员",
		"role": "admin",
	})
}

// 仪表盘相关处理函数
func handleGetDashboardData(c *gin.Context) {
	// 获取仪表盘数据
	c.JSON(http.StatusOK, gin.H{
		"device_count": 128,
		"idc_count": 3,
		"cabinet_count": 56,
		"user_count": 25,
		"online_devices": 112,
		"offline_devices": 16,
		"alerts": 0,
		"device_types": []gin.H{
			{"type": "服务器", "count": 85},
			{"type": "网络设备", "count": 23},
			{"type": "存储设备", "count": 12},
			{"type": "其他", "count": 8},
		},
		"ip_usage": []gin.H{
			{"subnet": "192.168.1.0/24", "total": 254, "used": 128, "free": 126},
			{"subnet": "192.168.2.0/24", "total": 254, "used": 95, "free": 159},
			{"subnet": "192.168.3.0/24", "total": 254, "used": 230, "free": 24},
		},
		"recent_activities": []gin.H{
			{"id": 1, "type": "login", "user": "admin", "message": "管理员登录系统", "time": "2024-01-15 10:30:00"},
			{"id": 2, "type": "create", "user": "admin", "message": "新增设备服务器-001", "time": "2024-01-15 10:15:00"},
			{"id": 3, "type": "update", "user": "operator", "message": "更新设备状态", "time": "2024-01-15 09:45:00"},
			{"id": 4, "type": "login", "user": "operator", "message": "操作员登录系统", "time": "2024-01-15 09:30:00"},
		},
	})
}

// 设备管理相关处理函数
func handleListDevices(c *gin.Context) {
	// 获取设备列表
	c.JSON(http.StatusOK, gin.H{
		"total": 128,
		"page": 1,
		"page_size": 10,
		"data": []gin.H{
			{
				"id": 1,
				"hostname": "server-001",
				"ip": "192.168.1.101",
				"mac": "00:11:22:33:44:55",
				"type": "服务器",
				"model": "Dell R740",
				"cpu": "Intel Xeon Gold 6230",
				"memory": "128GB",
				"disk": "2TB SSD",
				"status": "在线",
				"idc": "北京机房",
				"cabinet": "A01",
				"u_position": "1-2",
				"created_at": "2024-01-15 10:15:00",
				"updated_at": "2024-01-15 10:15:00",
			},
			// 更多设备数据...
		},
	})
}

func handleSearchDevices(c *gin.Context) {
	// 搜索设备
	c.JSON(http.StatusOK, gin.H{
		"total": 5,
		"page": 1,
		"page_size": 10,
		"data": []gin.H{
			// 搜索结果...
		},
	})
}

func handleGetDevice(c *gin.Context) {
	// 获取设备详情
	c.JSON(http.StatusOK, gin.H{
		"id": 1,
		"hostname": "server-001",
		"ip": "192.168.1.101",
		"mac": "00:11:22:33:44:55",
		"type": "服务器",
		"model": "Dell R740",
		"cpu": "Intel Xeon Gold 6230",
		"memory": "128GB",
		"disk": "2TB SSD",
		"os": "Ubuntu 20.04 LTS",
		"os_version": "20.04.3",
		"status": "在线",
		"idc": "北京机房",
		"cabinet": "A01",
		"u_position": "1-2",
		"purchase_date": "2023-06-15",
		"warranty_end_date": "2026-06-15",
		"supplier": "Dell",
		"contact": "张工",
		"phone": "13800138000",
		"notes": "生产环境核心服务器",
		"created_at": "2024-01-15 10:15:00",
		"updated_at": "2024-01-15 10:15:00",
	})
}

func handleCreateDevice(c *gin.Context) {
	// 创建设备
	c.JSON(http.StatusOK, gin.H{
		"message": "设备创建成功",
		"data": gin.H{
			"id": 129,
			// 新创建设备信息...
		},
	})
}

func handleUpdateDevice(c *gin.Context) {
	// 更新设备
	c.JSON(http.StatusOK, gin.H{
		"message": "设备更新成功",
		"data": gin.H{
			"id": c.Param("id"),
			// 更新后的设备信息...
		},
	})
}

func handleDeleteDevice(c *gin.Context) {
	// 删除设备
	c.JSON(http.StatusOK, gin.H{
		"message": "设备删除成功",
	})
}

// 机房管理相关处理函数
func handleListIDCs(c *gin.Context) {
	// 获取机房列表
	c.JSON(http.StatusOK, gin.H{
		"total": 3,
		"data": []gin.H{
			{"id": 1, "name": "北京机房", "code": "BJ", "address": "北京市海淀区中关村科技园", "contact": "李工", "phone": "13900139000", "cabinet_count": 25},
			{"id": 2, "name": "上海机房", "code": "SH", "address": "上海市浦东新区张江高科技园区", "contact": "王工", "phone": "13700137000", "cabinet_count": 18},
			{"id": 3, "name": "广州机房", "code": "GZ", "address": "广东省广州市天河区软件园", "contact": "陈工", "phone": "13600136000", "cabinet_count": 13},
		},
	})
}

func handleGetIDC(c *gin.Context) {
	// 获取机房详情
	c.JSON(http.StatusOK, gin.H{
		"id": 1,
		"name": "北京机房",
		"code": "BJ",
		"address": "北京市海淀区中关村科技园",
		"contact": "李工",
		"phone": "13900139000",
		"email": "idc-bj@example.com",
		"bandwidth": "10Gbps",
		"power": "双路供电",
		"ups": "APC UPS",
		"air_conditioner": "精密空调",
		"fire_protection": "七氟丙烷气体灭火",
		"cabinet_count": 25,
		"device_count": 65,
		"created_at": "2023-01-01 00:00:00",
		"updated_at": "2023-06-15 10:00:00",
	})
}

func handleCreateIDC(c *gin.Context) {
	// 创建机房
	c.JSON(http.StatusOK, gin.H{
		"message": "机房创建成功",
		"data": gin.H{
			"id": 4,
			// 新创建机房信息...
		},
	})
}

func handleUpdateIDC(c *gin.Context) {
	// 更新机房
	c.JSON(http.StatusOK, gin.H{
		"message": "机房更新成功",
		"data": gin.H{
			"id": c.Param("id"),
			// 更新后的机房信息...
		},
	})
}

func handleDeleteIDC(c *gin.Context) {
	// 删除机房
	c.JSON(http.StatusOK, gin.H{
		"message": "机房删除成功",
	})
}

// 其他路由处理函数
// 注意：这里为了简化，其他处理函数暂时返回空或简单响应
// 实际项目中需要根据业务逻辑实现
func handleListCabinets(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"total": 0, "data": []gin.H{} }) }
func handleGetCabinet(c *gin.Context) { c.JSON(http.StatusOK, gin.H{}) }
func handleCreateCabinet(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "创建成功"}) }
func handleUpdateCabinet(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "更新成功"}) }
func handleDeleteCabinet(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "删除成功"}) }

func handleListSubnets(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"total": 0, "data": []gin.H{} }) }
func handleGetSubnet(c *gin.Context) { c.JSON(http.StatusOK, gin.H{}) }
func handleCreateSubnet(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "创建成功"}) }
func handleUpdateSubnet(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "更新成功"}) }
func handleDeleteSubnet(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "删除成功"}) }

func handleListIPAddresses(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"total": 0, "data": []gin.H{} }) }
func handleGetIPAddress(c *gin.Context) { c.JSON(http.StatusOK, gin.H{}) }
func handleCreateIPAddress(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "创建成功"}) }
func handleUpdateIPAddress(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "更新成功"}) }
func handleDeleteIPAddress(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "删除成功"}) }
func handleAllocateIPAddress(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "分配成功"}) }
func handleReleaseIPAddress(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "释放成功"}) }

func handleGetAssetSummary(c *gin.Context) { c.JSON(http.StatusOK, gin.H{}) }
func handleGetAssetTypeDistribution(c *gin.Context) { c.JSON(http.StatusOK, gin.H{}) }
func handleGetIDCStatistics(c *gin.Context) { c.JSON(http.StatusOK, gin.H{}) }
func handleGetCabinetUsage(c *gin.Context) { c.JSON(http.StatusOK, gin.H{}) }
func handleGetIPUsage(c *gin.Context) { c.JSON(http.StatusOK, gin.H{}) }
func handleGetMaintenanceReport(c *gin.Context) { c.JSON(http.StatusOK, gin.H{}) }
func handleExportReport(c *gin.Context) { c.JSON(http.StatusOK, gin.H{}) }

func handleListUsers(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"total": 0, "data": []gin.H{} }) }
func handleGetUser(c *gin.Context) { c.JSON(http.StatusOK, gin.H{}) }
func handleCreateUser(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "创建成功"}) }
func handleUpdateUser(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "更新成功"}) }
func handleDeleteUser(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "删除成功"}) }

func handleListRoles(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"total": 0, "data": []gin.H{} }) }
func handleGetRole(c *gin.Context) { c.JSON(http.StatusOK, gin.H{}) }
func handleCreateRole(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "创建成功"}) }
func handleUpdateRole(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "更新成功"}) }
func handleDeleteRole(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "删除成功"}) }

func handleGetSettings(c *gin.Context) { c.JSON(http.StatusOK, gin.H{}) }
func handleUpdateSettings(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "更新成功"}) }

func handleListSystemLogs(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"total": 0, "data": []gin.H{} }) }

func main() {
	// 初始化数据库
	err := initDatabase()
	if err != nil {
		log.Printf("Warning: Failed to initialize database: %v", err)
		log.Printf("Running in demo mode with mock data...")
	}

	// 创建Gin实例
	router := gin.Default()

	// 配置CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 设置路由
	setupRoutes(router)

	// 获取端口号，默认为8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 启动服务器
	log.Printf("Server starting on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}