package mysql

import (
	"ai-rag-system/go-backend/mysql/models"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// 全局盐值（不用了）
// const secret = "liwenzhou.com"

// CheckUserExist 检查用户是否存在
func CheckUserExist(username string) (err error) {
	sqlStr := "select count(user_id) from user where username=?"
	var count int
	if err := db.QueryRow(sqlStr, username).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return errors.New("用户已存在")
	}
	return nil
}

// InsertUser 插入用户数据 —— **不加密密码**
func InsertUser(user *models.User) (err error) {
	// ❌ 去掉加密
	// user.Password = encryptPassword(user.Password)

	// 注册默认普通用户
	user.Role = 0
	user.CreatedAt = time.Now()

	// SQL
	sqlStr := "insert into user(user_id, username, password, role, created_at) values (?, ?, ?, ?, ?)"
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password, user.Role, user.CreatedAt)
	return
}

// Login 登录 —— **不加密密码，直接对比明文**
// Login 登录 —— 最终正确版
func Login(user *models.User) (err error) {
	// 1. 保存用户输入的密码
	inputPwd := user.Password

	// 2. 【重要】用临时变量接收查询结果，不污染原结构体
	var dbUser struct {
		UserID   int64  `db:"user_id"`
		Password string `db:"password"`
		Role     int    `db:"role"`
	}

	// 3. 查询
	sqlStr := "select user_id, password, role from user where username=?"
	err = db.Get(&dbUser, sqlStr, user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrorUserNotExist
		}
		return err
	}

	// 4. 正确对比 bcrypt 密码哈希
	if err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(inputPwd)); err != nil {
		return ErrorInvalidPassword
	}

	// 6. 赋值回去
	user.UserID = dbUser.UserID
	user.Password = dbUser.Password
	user.Role = dbUser.Role

	return nil
}

// 以下函数全部不用了
// 统一加密函数
// func encryptPassword(oPassword string) string {
// 	hash := md5.Sum([]byte(oPassword + secret))
// 	return hex.EncodeToString(hash[:])
// }

// 密码校验
// func checkPassword(password string, hashedPassword string) bool {
// 	return encryptPassword(password) == hashedPassword
// }
