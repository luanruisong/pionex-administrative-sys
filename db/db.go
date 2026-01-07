package db

import (
	"context"
	"pionex-administrative-sys/utils"
	"pionex-administrative-sys/utils/app"
	"pionex-administrative-sys/utils/logger"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var db *gorm.DB

// Init 初始化 SQLite 数据库连接
func init() {
	var err error
	db, err = gorm.Open(sqlite.Open(app.DBPath("data.db")), &gorm.Config{
		Logger: logger.NewGormLogger(),
	})
	if err != nil {
		logger.Fatal(err.Error())
	}
	if err = autoMigrate(); err != nil {
		logger.Fatal(err.Error())
	}
	if err = initializeData(); err != nil {
		logger.Fatal(err.Error())
	}
}

// AutoMigrate 自动建表
func autoMigrate() error {
	return db.AutoMigrate(
		&User{},
		&Coupon{},
	)
}

// initializeData 初始化基础数据
func initializeData() error {
	return db.Clauses(clause.OnConflict{DoNothing: true}).Create(&User{
		Name:      "管理员",
		Account:   "admin",
		Md5Pwd:    utils.MD5("123456"),
		Role:      MergeRole(AllRoles()...),
		CreatedAt: time.Now().UnixMilli(),
		UpdatedAt: time.Now().UnixMilli(),
	}).Error
}

// GetDB 获取数据库实例
func getDb(ctx context.Context) *gorm.DB {
	return db.WithContext(ctx)
}
