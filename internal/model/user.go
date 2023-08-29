package model

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Model
	Username string `json:"username" gorm:"type:varchar(100);index" binding:"required" example:"张三"` // 用户名
	Password string `json:"password" gorm:"type:varchar(100)" example:"123456"`                      // 密码（哈希值）
}

// Query 查询
func (r *User) Query() *gorm.DB {
	return DB.Model(r)
}

// EncryptPassword 加密密码
func (r *User) EncryptPassword() {
	hash, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("EncryptPassword error: %v", err)
	}
	r.Password = string(hash)
}

// ComparePassword 比较密码
func (r *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(r.Password), []byte(password))
	return err == nil
}
