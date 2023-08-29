package model

import (
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Setup initializes the database instance
func Setup() {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	Migration(db)
	DB = db
}

func migrate(db *gorm.DB, models ...interface{}) {
	err := db.AutoMigrate(models...)
	if err != nil {
		panic(err)
	}
}

type Model struct {
	ID        uint           `gorm:"primarykey" json:"id" example:"1" format:"int64"`                                     // 主键ID
	CreatedAt time.Time      `json:"createdAt" example:"2023-06-13T19:06:22.514+08:00"`                                   // 创建时间
	UpdatedAt time.Time      `json:"updatedAt" example:"2023-06-13T19:06:22.514+08:00"`                                   // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt" swaggertype:"string" example:"2023-06-13T19:06:22.514+08:00"` // 删除时间 - 软删除
}

// Migration migrate the schema
func Migration(db *gorm.DB) {
	// Migrate the schema
	migrate(db, &User{})
}
