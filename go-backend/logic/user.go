package logic

import (
	"ai-rag-system/go-backend/mysql"
	"ai-rag-system/go-backend/mysql/models"
	"ai-rag-system/go-backend/pkg/jwt"
	"ai-rag-system/go-backend/pkg/snowflake"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// 注册
func Signup(p *models.SignUpRequest) (err error) {
	//判断用户是否存在验证
	err = mysql.CheckUserExist(p.Username)
	if err != nil {
		return err
	}
	//生成uid
	userID := snowflake.GenID()
	//哈希密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	//构造一个User实例
	user := models.User{
		UserID:    userID,
		Username:  p.Username,
		Password:  string(hashedPassword),
		Role:      0, // 默认普通用户
		CreatedAt: time.Now(),
	}
	//将用户信息存入数据库
	err = mysql.InsertUser(&user)
	if err != nil {
		return err
	}
	return nil
}

// 登录
func Login(p *models.LoginRequest) (user *models.User, err error) {
	//1.根据用户登录的姓名密码组装一个结构体
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	//2.去数据库查询用户信息
	if err := mysql.Login(user); err != nil {
		return nil, err //直接返回：用户不存在/密码错误
	}
	//生成jwt
	token, err := jwt.GenToken(user.UserID, user.Username, user.Role)
	if err != nil {
		return
	}
	user.Token = token
	return
}
