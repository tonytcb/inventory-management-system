package domain

type Currency string

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	JPY Currency = "JPY"
	GBP Currency = "GBP"
	AUD Currency = "AUD"
)

func (c Currency) IsValid() bool {
	valid := []Currency{USD, EUR, JPY, GBP, AUD}

	for _, v := range valid {
		if v == c {
			return true
		}
	}

	return false
}

type CurrencyPair struct {
	From Currency
	To   Currency
}

func BuildPairs(currencies []Currency) []CurrencyPair {
	var pairs = make([]CurrencyPair, 0)

	for _, from := range currencies {
		for _, to := range currencies {
			if from == to {
				continue
			}

			pairs = append(pairs, CurrencyPair{From: from, To: to})
		}
	}

	return pairs
}
