package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/blog/internal/models"
	"github.com/blog/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PostController struct {
	Db *gorm.DB
}

func NewPostController(db *gorm.DB) *PostController {
	return &PostController{Db: db}
}

type PostListItem struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Title     string    `gorm:"size:64;not null" json:"title"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	UserID    uint      `json:"userId"`
}

func (pc *PostController) GetAllPosts(ctx *gin.Context) {
	pageNow, _ := strconv.Atoi(ctx.DefaultQuery("pageNow", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))
	if pageNow < 1 {
		pageNow = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var posts []PostListItem
	err := pc.Db.Model(&models.Post{}).Order("created_at desc, id desc").Offset((pageNow - 1) * pageSize).Limit(pageSize).Find(&posts).Error
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code": utils.FAIL.Code,
			"msg":  "查询文章失败",
		})
		return
	}

	var totalCount int64
	err = pc.Db.Model(&models.Post{}).Count(&totalCount).Error
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code": utils.FAIL.Code,
			"msg":  "查询文章数量失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": utils.SUCCESS.Code,
		"msg":  utils.SUCCESS.Msg,
		"data": gin.H{
			"pageNow":    pageNow,
			"pageSize":   pageSize,
			"totalCount": totalCount,
			"list":       posts,
		},
	})
}

func (pc *PostController) GetPostDetail(ctx *gin.Context) {
	idStr := ctx.Query("id")
	var post models.Post

	id, _ := strconv.Atoi(idStr)
	err := pc.Db.Where("id = ?", id).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "user_name", "nick_name")
	}).First(&post).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"code": utils.DATA_NOT_EXISTS.Code,
				"msg":  "文章不存在",
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code": utils.FAIL.Code,
			"msg":  "查询文章失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": utils.SUCCESS.Code,
		"msg":  utils.SUCCESS.Msg,
		"data": post,
	})
}

func (pc *PostController) Create(ctx *gin.Context) {
	userID, exists := ctx.Get(utils.USER_ID_CTX_KEY)
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": utils.NO_AUTH.Code,
			"msg":  "用户未认证",
		})
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code": utils.FAIL.Code,
			"msg":  "服务错误",
		})
		return
	}

	req := &models.CreatePostReq{}

	if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": utils.PARAM_ERROR.Code,
			"msg":  err.Error(),
		})
		return
	}

	post := &models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  uid,
	}

	err := pc.Db.Create(post).Error
	if err != nil {
		logrus.Errorf("创建文章失败, userId:%v, err:%v", uid, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": utils.FAIL.Code,
			"msg":  "创建文章失败",
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"code": utils.SUCCESS.Code,
		"msg":  "创建文章成功",
	})
}

func (pc *PostController) Update(ctx *gin.Context) {
	req := &models.UpdatePostReq{}
	if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": utils.PARAM_ERROR.Code,
			"msg":  err.Error(),
		})
		return
	}

	postCtx, exists := ctx.Get(utils.POST_CTX_KEY)
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code": utils.DATA_NOT_EXISTS.Code,
			"msg":  "文章不存在",
		})
		return
	}

	post := postCtx.(models.Post)
	post.Title = req.Title
	post.Content = req.Content

	err := pc.Db.Save(&post).Error

	if err != nil {
		logrus.Errorf("更新文章失败, postID:%v, err:%v", post.ID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": utils.FAIL.Code,
			"msg":  "更新文章失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": utils.SUCCESS.Code,
		"msg":  "更新文章成功",
	})
}

func (pc *PostController) Delete(ctx *gin.Context) {
	postCtx, exists := ctx.Get(utils.POST_CTX_KEY)
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code": utils.DATA_NOT_EXISTS.Code,
			"msg":  "文章不存在",
		})
		return
	}

	post := postCtx.(models.Post)

	err := pc.Db.Delete(&post).Error

	if err != nil {
		logrus.Errorf("删除文章失败, postID:%v, err:%v", post.ID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": utils.FAIL.Code,
			"msg":  "删除文章失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": utils.SUCCESS.Code,
		"msg":  "删除文章成功",
	})
}
