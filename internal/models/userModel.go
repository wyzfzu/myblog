package models

import (
	"time"

	"github.com/blog/internal/utils"
	"gorm.io/gorm"
)

/*
*
设计数据库表结构，至少包含以下几个表：
users 表：存储用户信息，包括 id 、 username 、 password 、 email 等字段。
posts 表：存储博客文章信息，包括 id 、 title 、 content 、 user_id （关联 users 表的 id ）、 created_at 、 updated_at 等字段。
comments 表：存储文章评论信息，包括 id 、 content 、 user_id （关联 users 表的 id ）、 post_id （关联 posts 表的 id ）、 created_at 等字段。
*/
type User struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	UserName   string         `gorm:"size:64;not null;unique:uniqueIndex" json:"userName"`
	Password   string         `gorm:"size:64;not null" json:"-"`
	NickName   string         `gorm:"size:64;not null" json:"nickName"`
	Gender     uint8          `json:"gender"`
	Age        uint8          `json:"age"`
	Email      string         `json:"email"`
	UserStatus uint8          `json:"userStatus"`

	Posts []Post `json:"posts,omitempty"`
}

func (u *User) BeforeCreate(db *gorm.DB) error {
	newPwd, err := utils.EncryptPassword(u.Password)
	if err != nil {
		return err
	}

	u.Password = newPwd

	return nil
}

type UserRegisterReq struct {
	UserName string `json:"userName" binding:"required,min=2,max=64"`
	Password string `json:"password" binding:"required"`
	NickName string `json:"nickName" binding:"max=64"`
	Gender   uint8  `json:"gender"`
	Age      uint8  `json:"age"`
	Email    string `json:"email"`
}

type UserLoginReq struct {
	UserName string `json:"userName" binding:"required,min=2,max=64"`
	Password string `json:"password" binding:"required,min=1"`
}
