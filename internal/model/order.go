package model

import (
	"encoding/json"
	"github.com/djokcik/gophermart/pkg/luhn"
	"time"
)

const (
	StatusNew        Status = "NEW"        // The order has been loaded into the system but has not been processed
	StatusRegistered Status = "REGISTERED" // The order was register from user but not calculated
	StatusProcessing Status = "PROCESSING" // The award for the order is calculated
	StatusProcessed  Status = "PROCESSED"  // The order data has been verified and the calculation information has been successfully obtained.
	StatusInvalid    Status = "INVALID"    // The remuneration calculation system refused to calculate
)

type (
	Status       string
	OrderId      string
	UploadedTime time.Time

	Order struct {
		Id         OrderId      `json:"number"`
		UserId     int          `json:"-"`
		Status     Status       `json:"status"`
		UploadedAt UploadedTime `json:"uploaded_at"`
		Accrual    int          `json:"accrual,omitempty"`
	}
)

func (s Status) Valid() bool {
	return s == StatusNew ||
		s == StatusRegistered ||
		s == StatusProcessed ||
		s == StatusProcessing ||
		s == StatusInvalid
}

func (o OrderId) Valid() bool {
	return luhn.Validate(string(o))
}

// MarshalJSON реализует интерфейс json.Marshaler.
func (s UploadedTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(s).Format(time.RFC3339))
}

// UnmarshalJSON реализует интерфейс json.Unmarshaler.
func (s *UploadedTime) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &s)
}
