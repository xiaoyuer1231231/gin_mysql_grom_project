package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/xiaoyuer1231231/gin_mysql_grom_project/config"
	"github.com/xiaoyuer1231231/gin_mysql_grom_project/controller"
	"github.com/xiaoyuer1231231/gin_mysql_grom_project/database"
	"github.com/xiaoyuer1231231/gin_mysql_grom_project/middleware"
)

// @title          博客系统 API
// @version        1.0
// @description    这是一个博客系统的后端 API 文档
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT 认证令牌，格式: Bearer <token>
func main() {
	cfg, err := config.LoadFromFile("config/config.yaml")
	if err != nil {
		fmt.Errorf("failed to migrate config: %w", err)
	}
	db, error := database.InitDataBase(cfg)
	if error != nil {
		panic(error)
	}

	fmt.Println("ssssss", cfg.JWT.ExpirationHours)
	//Initialize controllers
	authController := controller.NewAuthController(db, cfg)
	postController := controller.NewPostController(db)
	commentController := controller.NewCommentController(db)
	if _, err := os.Stat("./docs/swagger.json"); os.IsNotExist(err) {
		fmt.Printf("❌ docs/swagger.json 不存在")
	}
	router := gin.Default()
	router.Use(middleware.LoggerMiddleware())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}
		//  创建文章需要认证认证路由
		post := api.Group("/post")
		post.Use(middleware.AuthMiddleware(cfg))
		{
			post.POST("/createPost", postController.CreatePost)
			post.GET("/queryPost", postController.QueryPost)
			post.POST("/uptDateById", postController.UptDateById)
			post.DELETE("/deleteById", postController.DeleteById)
		}
		//评论功能
		comment := api.Group("/comment")
		comment.Use(middleware.AuthMiddleware(cfg))
		{
			comment.POST("/createComment", commentController.CreateComment)
			comment.GET("/queryComment", commentController.QueryComment)
		}
	}

	port := ":" + cfg.Server.Port
	fmt.Printf("🚀 服务器启动在 http://localhost%s\n", port)
	fmt.Printf("📊 健康检查: http://localhost%s/health\n", port)

	// 启动服务器
	if err := router.Run(port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
