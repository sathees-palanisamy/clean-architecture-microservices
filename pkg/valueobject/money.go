package valueobject

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
)

// Money represents a monetary value in a specific currency.
// For simplicity in this microservices setup, we'll assume USD and just track the amount.
// In a real-world scenario, you would track Currency Code as well.
type Money struct {
	amount float64 // Storing as float64 for simplicity in this demo, but typically int64 (cents) or big.Decimal is better.
}

// NewMoney creates a new Money instance safely.
func NewMoney(amount float64) Money {
	// Round to 2 decimal places to avoid standard float errors
	rounded := math.Round(amount*100) / 100
	return Money{amount: rounded}
}

// Amount returns the float value of the money
func (m Money) Amount() float64 {
	return m.amount
}

// IsNegative returns true if the amount is less than zero
func (m Money) IsNegative() bool {
	return m.amount < 0
}

// IsZero returns true if the amount is zero
func (m Money) IsZero() bool {
	return m.amount == 0
}

// Add adds two Money instances
func (m Money) Add(other Money) Money {
	return NewMoney(m.amount + other.amount)
}

// Multiply multiplies Money by a factor (e.g. quantity)
func (m Money) Multiply(factor int) Money {
	return NewMoney(m.amount * float64(factor))
}

// MarshalJSON implements json.Marshaler
func (m Money) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.amount)
}

// UnmarshalJSON implements json.Unmarshaler
func (m *Money) UnmarshalJSON(data []byte) error {
	var val float64
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	*m = NewMoney(val)
	return nil
}

// Value implements driver.Valuer for database storage
func (m Money) Value() (driver.Value, error) {
	return m.amount, nil
}

// Scan implements sql.Scanner for database retrieval
func (m *Money) Scan(value interface{}) error {
	if value == nil {
		*m = NewMoney(0)
		return nil
	}
	switch v := value.(type) {
	case float64:
		*m = NewMoney(v)
	case []byte:
		// Postgres might return decimals as strings/bytes
		// We need to parse it. For now let's try a simple approach assuming driver handles it or it's simple bytes
		var f float64
		_, err := fmt.Sscan(string(v), &f)
		if err != nil {
			return fmt.Errorf("failed to scan Money: %v", err)
		}
		*m = NewMoney(f)
	default:
		return fmt.Errorf("failed to scan Money: unexpected type %T", value)
	}
	return nil
}
