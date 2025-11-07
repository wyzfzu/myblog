package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/blog/configs"
	"github.com/blog/internal/models"
	"github.com/blog/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type UserController struct {
	Cfg *configs.Config
	Db  *gorm.DB
}

func NewUserController(cfg *configs.Config, db *gorm.DB) *UserController {
	return &UserController{Cfg: cfg, Db: db}
}

func (uc *UserController) Register(ctx *gin.Context) {
	req := &models.UserRegisterReq{}
	if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": utils.PARAM_ERROR.Code,
			"msg":  err.Error(),
		})
		return
	}

	// 校验参数合法性
	var userCount int64
	err := uc.Db.Model(&models.User{}).Where("user_name = ?", req.UserName).Count(&userCount).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": utils.FAIL.Code,
			"msg":  "注册失败",
		})
		return
	}

	if userCount > 0 {
		ctx.JSON(http.StatusConflict, gin.H{
			"code": utils.USER_NAME_EXISTS.Code,
			"msg":  "用户已存在",
		})
		return
	}

	newUser := &models.User{
		UserName: req.UserName,
		Password: req.Password,
		NickName: req.NickName,
		Age:      req.Age,
		Gender:   req.Gender,
		Email:    req.Email,
	}

	err = uc.Db.Create(newUser).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": utils.FAIL.Code,
			"msg":  "注册失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": utils.SUCCESS.Code,
		"msg":  utils.SUCCESS.Msg,
	})
}

func (uc *UserController) Login(ctx *gin.Context) {
	req := &models.UserLoginReq{}
	if err := ctx.ShouldBindBodyWithJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": utils.PARAM_ERROR.Code,
			"data": err.Error(),
		})
		return
	}

	// 验证用户
	var user models.User
	err := uc.Db.Model(&models.User{}).Where("user_name = ?", req.UserName).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code": utils.USER_NOT_EXISTS.Code,
			"msg":  utils.USER_NOT_EXISTS.Msg,
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"code": utils.FAIL.Code,
			"msg":  "登录失败",
		})
		return
	}

	// 验证密码
	if !utils.EqualPassword(user.Password, req.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code": utils.FAIL.Code,
			"msg":  "用户名或密码错误",
		})
		return
	}

	// 生成jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   user.ID,
		"userName": user.UserName,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(uc.Cfg.JWT.SecretKey))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code": utils.FAIL.Code,
			"msg":  "生成token失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"code":  utils.SUCCESS.Code,
		"msg":   utils.SUCCESS.Msg,
	})
}
