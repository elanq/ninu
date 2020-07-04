package ninu_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/elanq/ninu"
	"github.com/stretchr/testify/assert"
)

func TestReadTransaction(t *testing.T) {
	cases := []struct {
		name        string
		input       string
		expected    *ninu.InputTransaction
		shouldError bool
	}{
		{
			"normal_case",
			"Makanan 20000",
			&ninu.InputTransaction{
				Amount:   20000,
				Category: "Makanan",
			},
			false,
		},
		{
			"normal_case_with_multiple_words",
			"Makanan ringan coca cola pepsi kentang manis 200000",
			&ninu.InputTransaction{
				Amount:   200000,
				Category: "Makanan ringan coca cola pepsi kentang manis",
			},
			false,
		},
		{
			"parsing_error",
			"Makanan hehe",
			nil,
			true,
		},
		{
			"invalid_params",
			"Makanan",
			nil,
			true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual, err := ninu.ReadTransaction(c.input)
			assert.Equal(t, c.shouldError, err != nil)
			if c.expected != nil {
				assert.Equal(t, c.expected.Amount, actual.Amount)
				assert.Equal(t, c.expected.Category, actual.Category)
			}
		})
	}
}

func TestToValueRange(t *testing.T) {
	cases := []struct {
		name     string
		input    *ninu.InputTransaction
		expected string
	}{
		{
			"normal_case",
			&ninu.InputTransaction{
				Date:     time.Date(1992, 9, 6, 0, 0, 0, 0, time.Local),
				Category: "Makanan",
				Amount:   20000,
			},
			`
				{
					"majorDimension": "ROWS",
					"values": [["06/09/1992", "Makanan", 20000]]
				}
			`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.input.ToValueRange()
			actualJSON, err := json.Marshal(actual)
			assert.NoError(t, err)

			assert.JSONEq(t, c.expected, string(actualJSON))
		})
	}
}
