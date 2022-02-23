package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/djokcik/gophermart/internal/config"
	"github.com/djokcik/gophermart/internal/model"
	"github.com/djokcik/gophermart/pkg/logging"
	"github.com/rs/zerolog"
	"io"
	"net/http"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type accrualClient struct {
	accrualAddress string
	client         httpClient
}

func (o accrualClient) GetOrder(ctx context.Context, orderID model.OrderID) (AccrualResponse, error) {
	url := fmt.Sprintf("%s/api/orders/%s", o.accrualAddress, orderID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		o.Log(ctx).Error().Err(err).Msg("request was interrupted")
		return AccrualResponse{}, err
	}

	res, err := o.client.Do(req)
	if err != nil {
		o.Log(ctx).Error().Err(err).Msg("request ended with error")
		return AccrualResponse{}, err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return AccrualResponse{}, &ErrAccrualResponse{
			Code: res.StatusCode,
			Body: string(body),
		}
	}

	var accrual AccrualResponse
	err = json.NewDecoder(res.Body).Decode(&accrual)
	if err != nil {
		o.Log(ctx).Error().Err(err).Msg("")
		return AccrualResponse{}, fmt.Errorf("invalid decode response: %w", err)
	}

	o.Log(ctx).Info().Msgf("Finished accrual response: %+v", accrual)

	return accrual, nil
}

func (o accrualClient) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "accrualClient").Logger()

	return &logger
}

func NewAccrualClient(cfg config.Config) AccrualClient {
	return &accrualClient{
		client:         &http.Client{},
		accrualAddress: cfg.AccrualSystemAddress,
	}
}
