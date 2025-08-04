package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

func main() {
	// 创建一个新的 Gin 引擎
	r := gin.Default()

	// 定义一个 POST 路由来接收稽核日志
	r.POST("/audit", func(c *gin.Context) {
		var auditLog map[string]interface{}

		// 绑定 JSON 数据到 auditLog 变量
		if err := c.ShouldBindJSON(&auditLog); err != nil {
			// 如果解析失败，返回 400 错误
			c.JSON(400, gin.H{
				"status":  "error",
				"message": "Invalid JSON data",
			})
			return
		}

		// 打印收到的稽核日志
		fmt.Println("Received audit log:")
		fmt.Printf("%+v\n", auditLog)

		// 在这里可以进行进一步的处理，例如将日志存储到数据库、转发到其他系统等

		// 返回成功响应
		c.JSON(200, gin.H{
			"status": "success",
		})
	})

	// 设置一个 POST 路由用于接收 MinIO 的 Webhook 日志
	r.POST("/sever-log", func(c *gin.Context) {
		// 读取请求体（Request Body）
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			// 如果读取失败，返回 400 错误
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 打印请求的 body
		fmt.Println("Received MinIO log:", string(body))

		// 返回一个成功的响应
		c.JSON(http.StatusOK, gin.H{"status": "received"})
	})

	// 启动 HTTP 服务器，监听 8080 端口
	if err := r.Run(":8095"); err != nil {
		log.Fatal("Unable to start server: ", err)
	}
}
