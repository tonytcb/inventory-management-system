package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildPairs(t *testing.T) {
	type args struct {
		currencies []Currency
	}
	tests := []struct {
		name string
		args args
		want []CurrencyPair
	}{
		{
			name: "single currency",
			args: args{currencies: []Currency{USD}},
			want: []CurrencyPair{},
		},
		{
			name: "two currencies",
			args: args{currencies: []Currency{USD, EUR}},
			want: []CurrencyPair{
				{From: USD, To: EUR},
				{From: EUR, To: USD},
			},
		},
		{
			name: "three currencies",
			args: args{currencies: []Currency{USD, EUR, JPY}},
			want: []CurrencyPair{
				{From: USD, To: EUR},
				{From: USD, To: JPY},
				{From: EUR, To: USD},
				{From: EUR, To: JPY},
				{From: JPY, To: USD},
				{From: JPY, To: EUR},
			},
		},
		{
			name: "four currencies",
			args: args{currencies: []Currency{USD, EUR, JPY, GBP}},
			want: []CurrencyPair{
				{From: USD, To: EUR},
				{From: USD, To: JPY},
				{From: USD, To: GBP},
				{From: EUR, To: USD},
				{From: EUR, To: JPY},
				{From: EUR, To: GBP},
				{From: JPY, To: USD},
				{From: JPY, To: EUR},
				{From: JPY, To: GBP},
				{From: GBP, To: USD},
				{From: GBP, To: EUR},
				{From: GBP, To: JPY},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, BuildPairs(tt.args.currencies), "BuildPairs(%v)", tt.args.currencies)
		})
	}
}
