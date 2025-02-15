package models

import "errors"

var (
	ErrInsufficientFunds      = errors.New("insufficient funds for the purchase")
	ErrUserNotFound           = errors.New("user not found")
	ErrMerchItemNotFound      = errors.New("merch_item not found")
	ErrMerchItemAlreadyExists = errors.New("merch_item already exists")
	ErrNotEnoughFunds         = errors.New("not enough funds for the purchase")
)
