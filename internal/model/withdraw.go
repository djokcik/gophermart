package model

import (
	"errors"
)

var (
	ErrInvalidOrderID = errors.New("service: invalid orderID")
)

type (
	WithdrawRequestDto struct {
		OrderID OrderID `json:"order"`
		Sum     Amount  `json:"sum"`
	}

	Withdraw struct {
		ID          int          `json:"-"`
		OrderID     OrderID      `json:"order"`
		Sum         Amount       `json:"sum"`
		ProcessedAt UploadedTime `json:"processed_at"`
		UserID      int          `json:"-"`
	}
)

func (w WithdrawRequestDto) Validate() error {
	if !w.OrderID.Valid() {
		return ErrInvalidOrderID
	}

	if w.Sum <= 0 {
		return errors.New("withdraw validate: invalid sum")
	}

	return nil
}
