package models

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrMerchItemNotFound = errors.New("merch item not found")
	ErrNotEnoughCoins    = errors.New("not enough coins for the purchase")
)
