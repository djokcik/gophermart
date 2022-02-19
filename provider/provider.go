package provider

import (
	"context"
	"fmt"
	"github.com/djokcik/gophermart/internal/model"
)

type (
	AccrualClient interface {
		GetOrder(ctx context.Context, orderId model.OrderId) (AccrualResponse, error)
	}

	AccrualResponse struct {
		Order   model.OrderId `json:"order"`
		Status  model.Status  `json:"status"`
		Accrual model.Accrual `json:"accrual"`
	}

	ErrAccrualResponse struct {
		Code int
		Body string
	}
)

func (e ErrAccrualResponse) Error() string {
	return fmt.Sprintf("accrual: failed to request with status: %d, body: %s", e.Code, e.Body)
}
