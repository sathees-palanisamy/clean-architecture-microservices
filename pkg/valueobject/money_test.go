package valueobject

import (
	"encoding/json"
	"testing"
)

func TestMoney_Arithmetic(t *testing.T) {
	m1 := NewMoney(10.50)
	m2 := NewMoney(5.25)

	// Addition
	sum := m1.Add(m2)
	if sum.Amount() != 15.75 {
		t.Errorf("Expected 15.75, got %v", sum.Amount())
	}

	// Multiplication
	prod := m1.Multiply(3)
	if prod.Amount() != 31.50 {
		t.Errorf("Expected 31.50, got %v", prod.Amount())
	}
}

func TestMoney_Rounding(t *testing.T) {
	m := NewMoney(10.5555)
	if m.Amount() != 10.56 {
		t.Errorf("Expected 10.56, got %v", m.Amount())
	}
}

func TestMoney_JSON(t *testing.T) {
	m := NewMoney(100.50)
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "100.5" {
		t.Errorf("Expected 100.5, got %s", string(data))
	}

	var m2 Money
	err = json.Unmarshal(data, &m2)
	if err != nil {
		t.Fatal(err)
	}

	if m2.Amount() != 100.50 {
		t.Errorf("Expected 100.50, got %v", m2.Amount())
	}
}
