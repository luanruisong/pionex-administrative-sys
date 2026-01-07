package db

// CouponType 卡券类型
type CouponType struct {
	Type int    `json:"type"`
	Name string `json:"name"`
}

// 卡券类型定义
const (
	couponTypeFitness = iota + 1
)

var (
	CouponTypeFitness = CouponType{Type: couponTypeFitness, Name: "健身卡"}
)

// AllCouponTypes 获取所有卡券类型
func AllCouponTypes() []CouponType {
	return []CouponType{
		CouponTypeFitness,
	}
}

// GetCouponTypeName 根据类型获取名称
func GetCouponTypeName(t int) string {
	for _, ct := range AllCouponTypes() {
		if ct.Type == t {
			return ct.Name
		}
	}
	return "未知类型"
}

// IsValidCouponType 检查是否为有效的卡券类型
func IsValidCouponType(t int) bool {
	for _, ct := range AllCouponTypes() {
		if ct.Type == t {
			return true
		}
	}
	return false
}