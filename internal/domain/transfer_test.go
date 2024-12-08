package domain

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestTransfer_ConvertAmounts(t *testing.T) {
	tests := []struct {
		name                    string
		originalAmount          decimal.Decimal
		rate                    decimal.Decimal
		margin                  decimal.Decimal
		expectedConvertedAmount decimal.Decimal
		expectedFinalAmount     decimal.Decimal
	}{
		{
			name:                    "No margin",
			originalAmount:          decimal.NewFromFloat(100.0),
			rate:                    decimal.NewFromFloat(1.2),
			margin:                  decimal.Zero,
			expectedConvertedAmount: decimal.NewFromFloat(120.0),
			expectedFinalAmount:     decimal.NewFromFloat(120.0),
		},
		{
			name:                    "With margin",
			originalAmount:          decimal.NewFromFloat(100.0),
			rate:                    decimal.NewFromFloat(1.2),
			margin:                  decimal.NewFromFloat(0.05),
			expectedConvertedAmount: decimal.NewFromFloat(120.0),
			expectedFinalAmount:     decimal.NewFromFloat(126.0),
		},
		{
			name:                    "Zero original amount",
			originalAmount:          decimal.Zero,
			rate:                    decimal.NewFromFloat(1.2),
			margin:                  decimal.NewFromFloat(0.05),
			expectedConvertedAmount: decimal.Zero,
			expectedFinalAmount:     decimal.Zero,
		},
		{
			name:                    "Zero rate",
			originalAmount:          decimal.NewFromFloat(100.0),
			rate:                    decimal.Zero,
			margin:                  decimal.NewFromFloat(0.05),
			expectedConvertedAmount: decimal.Zero,
			expectedFinalAmount:     decimal.Zero,
		},
		{
			name:                    "Negative margin",
			originalAmount:          decimal.NewFromFloat(100.0),
			rate:                    decimal.NewFromFloat(1.2),
			margin:                  decimal.NewFromFloat(-0.05),
			expectedConvertedAmount: decimal.NewFromFloat(120.0),
			expectedFinalAmount:     decimal.NewFromFloat(114.0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transfer := Transfer{
				OriginalAmount: tt.originalAmount,
			}

			transfer.ConvertAmounts(tt.rate, tt.margin)

			assert.Equal(t, tt.expectedConvertedAmount.String(), transfer.ConvertedAmount.String())
			assert.Equal(t, tt.expectedFinalAmount.String(), transfer.FinalAmount.String())
		})
	}
}

func TestTransfer_Margin(t *testing.T) {
	tests := []struct {
		name            string
		convertedAmount decimal.Decimal
		finalAmount     decimal.Decimal
		expectedMargin  decimal.Decimal
	}{
		{
			name:            "Positive margin",
			convertedAmount: decimal.NewFromFloat(100.0),
			finalAmount:     decimal.NewFromFloat(105.0),
			expectedMargin:  decimal.NewFromFloat(5.0),
		},
		{
			name:            "Zero margin",
			convertedAmount: decimal.NewFromFloat(100.0),
			finalAmount:     decimal.NewFromFloat(100.0),
			expectedMargin:  decimal.Zero,
		},
		{
			name:            "Negative margin",
			convertedAmount: decimal.NewFromFloat(100.0),
			finalAmount:     decimal.NewFromFloat(95.0),
			expectedMargin:  decimal.Zero,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transfer := Transfer{
				ConvertedAmount: tt.convertedAmount,
				FinalAmount:     tt.finalAmount,
			}

			assert.Equal(t, tt.expectedMargin.String(), transfer.Margin().String())
		})
	}
}
