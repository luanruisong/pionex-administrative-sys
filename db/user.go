package db

import (
	"context"
	"errors"
	"pionex-administrative-sys/utils"
)

type User struct {
	Id         int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Name       string `gorm:"column:name;type:varchar(64);not null"`
	Account    string `gorm:"column:account;type:varchar(128);uniqueIndex;not null"`
	Md5Pwd     string `gorm:"column:md5_pwd;type:varchar(32);not null"`
	Role       int    `gorm:"column:role;default:0"` // 权限位: 1=admin, 2=login
	PrivateKey string `gorm:"column:private_key;type:varchar(128)"`
	CreatedAt  int64  `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt  int64  `gorm:"column:updated_at;autoUpdateTime:milli"`
}

func (User) TableName() string {
	return "users"
}

// HasRole 检查是否拥有指定权限
func (u User) HasRole(role CommonRole) bool {
	return u.Role&role.Role != 0
}

func (u User) CheckPwd(pwd string) error {
	if u.Md5Pwd != utils.MD5(pwd) {
		return errors.New("wrong password")
	}
	return nil
}

// CreateUser 创建用户
func CreateUser(ctx context.Context, user *User) error {
	return getDb(ctx).Create(user).Error
}

// GetUserById 根据ID查询用户
func GetUserById(ctx context.Context, id int64) (*User, error) {
	var user User
	err := getDb(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByAccount 根据账号查询用户
func GetUserByAccount(ctx context.Context, account string) (*User, error) {
	var user User
	err := getDb(ctx).Where("account = ?", account).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserList 查询用户列表
func GetUserList(ctx context.Context, offset, limit int) ([]*User, error) {
	var users []*User
	err := getDb(ctx).Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// UpdateUser 更新用户
func UpdateUser(ctx context.Context, user *User) error {
	return getDb(ctx).Save(user).Error
}

// UpdateUserFields 更新用户指定字段
func UpdateUserFields(ctx context.Context, id int64, fields map[string]interface{}) error {
	return getDb(ctx).Model(&User{}).Where("id = ?", id).Updates(fields).Error
}

// DeleteUser 删除用户
func DeleteUser(ctx context.Context, id int64) error {
	return getDb(ctx).Where("id = ?", id).Delete(&User{}).Error
}

// CountUsers 统计用户总数
func CountUsers(ctx context.Context) (int64, error) {
	var count int64
	err := getDb(ctx).Model(&User{}).Count(&count).Error
	return count, err
}
