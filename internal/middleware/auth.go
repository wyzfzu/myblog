package middleware

import (
	"net/http"

	"github.com/blog/configs"
	"github.com/blog/internal/models"
	"github.com/blog/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func JWTMiddleware(config *configs.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("token")
		if tokenString == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": utils.PARAM_ERROR.Code,
				"msg":  "缺少令牌",
			})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrHashUnavailable
			}
			return []byte(config.JWT.SecretKey), nil
		})

		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": utils.FAIL.Code,
				"msg":  "令牌无效或已过期",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": utils.FAIL.Code,
				"msg":  "令牌解析失败",
			})
			return
		}
		userID, ok := claims["userId"].(float64)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": utils.FAIL.Code,
				"msg":  "令牌解析成用户ID失败",
			})
			return
		}

		ctx.Set(utils.USER_ID_CTX_KEY, uint(userID))
		ctx.Next()
	}
}

func PostAuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get(utils.USER_ID_CTX_KEY)
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": utils.NO_AUTH.Code,
				"msg":  "用户未认证",
			})
			return
		}
		var reqParam struct {
			ID uint
		}
		ctx.ShouldBindBodyWithJSON(&reqParam)
		var post models.Post
		err := db.Model(&models.Post{}).Where("id = ?", reqParam.ID).First(&post).Error
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"code": utils.DATA_NOT_EXISTS.Code,
				"msg":  "文章不存在",
			})
			return
		}

		if userID.(uint) != post.UserID {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": utils.NO_AUTH.Code,
				"msg":  "没有权限",
			})
			return
		}

		ctx.Set(utils.POST_CTX_KEY, post)
		ctx.Next()
	}
}
