package db

import "errors"

var (
	ErrCouponAlreadyTaken = errors.New("coupon already taken")
)