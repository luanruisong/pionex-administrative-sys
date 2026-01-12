package db

import (
	"context"

	"gorm.io/gorm"
)

type Coupon struct {
	Id        int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Coupon    string `gorm:"column:coupon;type:varchar(128);uniqueIndex;not null"`
	Type      int    `gorm:"column:type;index;not null;default:1"` // 卡券类型: 1=健身卡
	Creator   int64  `gorm:"column:creator;index"`
	Taker     int64  `gorm:"column:taker;index"`
	CreatedAt int64  `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64  `gorm:"column:updated_at;autoUpdateTime:milli"`
}

// CouponFilter 卡券筛选条件
type CouponFilter struct {
	Type  *int  // 卡券类型
	Taken *bool // 是否已领取
}

// applyFilter 应用筛选条件
func (f CouponFilter) applyFilter(db *gorm.DB) *gorm.DB {
	if f.Type != nil {
		db = db.Where("type = ?", *f.Type)
	}
	if f.Taken != nil {
		if *f.Taken {
			db = db.Where("taker > 0")
		} else {
			db = db.Where("taker = 0")
		}
	}
	return db
}

func (Coupon) TableName() string {
	return "coupons"
}

// IsTaken 检查卡券是否已被领取
func (c Coupon) IsTaken() bool {
	return c.Taker > 0
}

// CreateCoupon 创建卡券
func CreateCoupon(ctx context.Context, coupon *Coupon) error {
	return getDb(ctx).Create(coupon).Error
}

// BatchCreateCoupons 批量创建卡券
func BatchCreateCoupons(ctx context.Context, coupons []*Coupon) error {
	return getDb(ctx).CreateInBatches(coupons, 100).Error
}

// GetCouponById 根据ID查询卡券
func GetCouponById(ctx context.Context, id int64) (*Coupon, error) {
	var coupon Coupon
	err := getDb(ctx).Where("id = ?", id).First(&coupon).Error
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

// GetCouponByCode 根据卡券码查询卡券
func GetCouponByCode(ctx context.Context, code string) (*Coupon, error) {
	var coupon Coupon
	err := getDb(ctx).Where("coupon = ?", code).First(&coupon).Error
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

// GetCouponList 查询卡券列表
func GetCouponList(ctx context.Context, offset, limit int) ([]*Coupon, error) {
	var coupons []*Coupon
	err := getDb(ctx).Order("id DESC").Offset(offset).Limit(limit).Find(&coupons).Error
	if err != nil {
		return nil, err
	}
	return coupons, nil
}

// GetCouponListByCreator 根据创建者查询卡券列表
func GetCouponListByCreator(ctx context.Context, creator int64, offset, limit int) ([]*Coupon, error) {
	var coupons []*Coupon
	err := getDb(ctx).Where("creator = ?", creator).Order("id DESC").Offset(offset).Limit(limit).Find(&coupons).Error
	if err != nil {
		return nil, err
	}
	return coupons, nil
}

// GetAvailableCoupons 获取未被领取的卡券列表
func GetAvailableCoupons(ctx context.Context, offset, limit int) ([]*Coupon, error) {
	var coupons []*Coupon
	err := getDb(ctx).Where("taker = 0").Order("id DESC").Offset(offset).Limit(limit).Find(&coupons).Error
	if err != nil {
		return nil, err
	}
	return coupons, nil
}

// UpdateCoupon 更新卡券
func UpdateCoupon(ctx context.Context, coupon *Coupon) error {
	return getDb(ctx).Save(coupon).Error
}

// UpdateCouponFields 更新卡券指定字段
func UpdateCouponFields(ctx context.Context, id int64, fields map[string]interface{}) error {
	return getDb(ctx).Model(&Coupon{}).Where("id = ?", id).Updates(fields).Error
}

// TakeCoupon 领取卡券
func TakeCoupon(ctx context.Context, id int64, taker int64) error {
	result := getDb(ctx).Model(&Coupon{}).Where("id = ? AND taker = 0", id).Update("taker", taker)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrCouponAlreadyTaken
	}
	return nil
}

// DeleteCoupon 删除卡券
func DeleteCoupon(ctx context.Context, id int64) error {
	return getDb(ctx).Where("id = ?", id).Delete(&Coupon{}).Error
}

// CountCoupons 统计卡券总数
func CountCoupons(ctx context.Context) (int64, error) {
	var count int64
	err := getDb(ctx).Model(&Coupon{}).Count(&count).Error
	return count, err
}

// CountCouponsByCreator 统计指定创建者的卡券总数
func CountCouponsByCreator(ctx context.Context, creator int64) (int64, error) {
	var count int64
	err := getDb(ctx).Model(&Coupon{}).Where("creator = ?", creator).Count(&count).Error
	return count, err
}

// CountAvailableCoupons 统计未被领取的卡券总数
func CountAvailableCoupons(ctx context.Context) (int64, error) {
	var count int64
	err := getDb(ctx).Model(&Coupon{}).Where("taker = 0").Count(&count).Error
	return count, err
}

// GetCouponListWithFilter 根据筛选条件查询卡券列表
func GetCouponListWithFilter(ctx context.Context, filter CouponFilter, offset, limit int) ([]*Coupon, error) {
	var coupons []*Coupon
	query := filter.applyFilter(getDb(ctx).Model(&Coupon{}))
	err := query.Order("id DESC").Offset(offset).Limit(limit).Find(&coupons).Error
	if err != nil {
		return nil, err
	}
	return coupons, nil
}

// CountCouponsWithFilter 根据筛选条件统计卡券总数
func CountCouponsWithFilter(ctx context.Context, filter CouponFilter) (int64, error) {
	var count int64
	query := filter.applyFilter(getDb(ctx).Model(&Coupon{}))
	err := query.Count(&count).Error
	return count, err
}

// GetCouponsByTaker 根据领取者查询卡券列表
func GetCouponsByTaker(ctx context.Context, taker int64, typeFilter *int, offset, limit int) ([]*Coupon, error) {
	var coupons []*Coupon
	query := getDb(ctx).Where("taker = ?", taker)
	if typeFilter != nil {
		query = query.Where("type = ?", *typeFilter)
	}
	err := query.Order("updated_at DESC").Offset(offset).Limit(limit).Find(&coupons).Error
	if err != nil {
		return nil, err
	}
	return coupons, nil
}

// CountCouponsByTaker 统计领取者的卡券总数
func CountCouponsByTaker(ctx context.Context, taker int64, typeFilter *int) (int64, error) {
	var count int64
	query := getDb(ctx).Model(&Coupon{}).Where("taker = ?", taker)
	if typeFilter != nil {
		query = query.Where("type = ?", *typeFilter)
	}
	err := query.Count(&count).Error
	return count, err
}

// CountAvailableCouponsByType 统计指定类型未被领取的卡券总数
func CountAvailableCouponsByType(ctx context.Context, couponType int) (int64, error) {
	var count int64
	err := getDb(ctx).Model(&Coupon{}).Where("taker = 0 AND type = ?", couponType).Count(&count).Error
	return count, err
}

// GetOneAvailableCouponByType 获取一个指定类型的未领取卡券
func GetOneAvailableCouponByType(ctx context.Context, couponType int) (*Coupon, error) {
	var coupon Coupon
	err := getDb(ctx).Where("taker = 0 AND type = ?", couponType).Order("id ASC").First(&coupon).Error
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

// GetLastTakenCouponByTakerAndType 获取用户最后领取的指定类型卡券
func GetLastTakenCouponByTakerAndType(ctx context.Context, taker int64, couponType int) (*Coupon, error) {
	var coupon Coupon
	err := getDb(ctx).Where("taker = ? AND type = ?", taker, couponType).Order("updated_at DESC").First(&coupon).Error
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}