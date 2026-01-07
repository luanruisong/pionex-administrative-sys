package db

type CommonRole struct {
	Name string
	Role int
}

// 权限位定义
const (
	roleAdmin = iota
	roleLogin
	roleStock
	roleAApplyCoupon
)

var (
	RoleAdmin = newForRole(
		"管理员", roleAdmin,
	)
	RoleLogin = newForRole(
		"登录", roleLogin,
	)
	RoleStock = newForRole(
		"库存管理", roleStock,
	)
	RoleApplyCoupon = newForRole(
		"卡券申请", roleAApplyCoupon,
	)
)

func newForRole(name string, role int) CommonRole {
	return CommonRole{
		Name: name,
		Role: 1 << role,
	}
}

func AllRoles() []CommonRole {
	return []CommonRole{
		RoleAdmin,
		RoleLogin,
		RoleStock,
		RoleApplyCoupon,
	}
}

func MergeRole(rs ...CommonRole) int {
	var x = 0
	for _, v := range rs {
		x = x | v.Role
	}
	return x
}
