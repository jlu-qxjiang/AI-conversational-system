package models

import (
	"time"
)

// 用户注册结构体
type SignUpRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// 用户登录结构体
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// User 用户模型
type User struct {
	UserID    int64     `db:"user_id"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	Role      int       `db:"role"`
	CreatedAt time.Time `db:"created_at"`
	Token     string
}
