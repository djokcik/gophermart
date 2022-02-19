package model

import (
	"errors"
	"time"
)

var (
	ErrInvalidOrderId = errors.New("service: invalid orderId")
)

type (
	WithdrawRequestDto struct {
		Order OrderId `json:"order"`
		Sum   Amount  `json:"sum"`
	}

	Withdraw struct {
		Id          int       `json:"-"`
		OrderId     OrderId   `json:"order"`
		Sum         Amount    `json:"sum"`
		ProcessedAt time.Time `json:"processed_at"`
		UserId      int       `json:"-"`
	}
)

func (w WithdrawRequestDto) Validate() error {
	if !w.Order.Valid() {
		return ErrInvalidOrderId
	}

	if w.Sum <= 0 {
		return errors.New("withdraw validate: invalid sum")
	}

	return nil
}
