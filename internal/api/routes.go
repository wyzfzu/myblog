package api

import (
	"github.com/blog/configs"
	"github.com/blog/internal/controllers"
	"github.com/blog/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(engine *gin.Engine, db *gorm.DB, cfg *configs.Config) {

	userController := controllers.NewUserController(cfg, db)
	postController := controllers.NewPostController(db)
	commentController := controllers.NewCommentController(db)

	apiPublic := engine.Use(middleware.ErrorHandler())
	{
		apiPublic.POST("/api/user/register", userController.Register)
		apiPublic.POST("/api/user/login", userController.Login)
		apiPublic.GET("/api/posts/all", postController.GetAllPosts)
		apiPublic.GET("/api/posts/detail", postController.GetPostDetail)
		apiPublic.GET("/api/comments/all", commentController.GetPostComments)
	}

	apiPrivate := engine.Use(middleware.JWTMiddleware(cfg)).Use(middleware.ErrorHandler())
	{
		apiPrivate.POST("/api/posts/create", postController.Create)
		apiPrivate.POST("/api/posts/update", middleware.PostAuthMiddleware(db), postController.Update)
		apiPrivate.POST("/api/posts/delete", middleware.PostAuthMiddleware(db), postController.Delete)
		apiPrivate.POST("/api/comments/create", commentController.Create)
	}
}
