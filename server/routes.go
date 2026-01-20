package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"liangxiong/demo/auth"
	"liangxiong/demo/controller"
	"liangxiong/demo/docs"
	"liangxiong/demo/internal/config"
	"liangxiong/demo/middleware"
)

func setupRoutes(engine *gin.Engine, cfg *config.Config, logger *zap.Logger, userController *controller.UserController, authController *controller.AuthController, jwtManager *auth.JWTManager) {
	docs.SwaggerInfo.Title = cfg.App.Name + " API"
	docs.SwaggerInfo.Version = "1.0.0"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Description = "Simplified Go Web API boilerplate"
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", cfg.Server.Address, cfg.Server.Port)

	engine.GET("/healthz", controller.Health)
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := engine.Group("/api/v1")
	{
		authGroup := api.Group("/auth")
		authGroup.POST("/login", authController.Login)

		userGroup := api.Group("/users")
		userGroup.Use(middleware.Auth(jwtManager, logger))
		userGroup.GET("", userController.List)
		userGroup.GET("/:id", userController.Get)
		userGroup.POST("", userController.Create)
		userGroup.PUT("/:id", userController.Update)
		userGroup.DELETE("/:id", userController.Delete)
	}
}
