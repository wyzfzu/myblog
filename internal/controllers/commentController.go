package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/blog/internal/models"
	"github.com/blog/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CommentController struct {
	Db *gorm.DB
}

func NewCommentController(db *gorm.DB) *CommentController {
	return &CommentController{Db: db}
}

type CommentItem struct {
	Content   string      `json:"content"`
	PostID    uint        `json:"postId"`
	CreatedAt time.Time   `json:"createdAt"`
	UserID    uint        `json:"userId"`
	User      models.User `json:"user"`
}

func (cc *CommentController) GetPostComments(ctx *gin.Context) {
	idStr := ctx.Query("postId")
	if idStr == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": utils.PARAM_ERROR.Code,
			"msg":  "文章ID为空",
		})
	}
	postID, _ := strconv.Atoi(idStr)
	var comments []CommentItem
	cc.Db.Model(&models.Comment{}).Where("post_id = ?", postID).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "user_name", "nick_name")
	}).Find(&comments)

	ctx.JSON(http.StatusOK, gin.H{
		"code": utils.SUCCESS.Code,
		"msg":  utils.SUCCESS.Msg,
		"data": comments,
	})
}

type CreateCommentReq struct {
	PostID  uint   `json:"postId" binding:"required"`
	Content string `json:"content" binding:"required,min=2"`
}

func (cc *CommentController) Create(ctx *gin.Context) {
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

	req := &CreateCommentReq{}
	if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": utils.PARAM_ERROR.Code,
			"msg":  err.Error(),
		})
		return
	}

	comment := &models.Comment{
		Content: req.Content,
		PostID:  req.PostID,
		UserID:  uid,
	}

	err := cc.Db.Create(comment).Error
	if err != nil {
		logrus.Errorf("评论文章失败, userId:%v, postId:%v, err:%v", uid, req.PostID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": utils.FAIL.Code,
			"msg":  "评论失败",
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"code": utils.SUCCESS.Code,
		"msg":  "评论文章成功",
	})
}
