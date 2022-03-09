package utils

import (
	"sort"
	"strings"

	"github.com/shopspring/decimal"
)

func SortStrAsc(values []string) {
	sort.Slice(values, func(i, j int) bool {
		if len(values[i]) == 0 && len(values[j]) == 0 {
			return false
		}
		if len(values[i]) == 0 || len(values[j]) == 0 {
			return len(values[i]) == 0
		}

		val1 := decimal.RequireFromString(values[i])
		val2 := decimal.RequireFromString(values[j])

		return val1.LessThan(val2)
	})
}

func SortStrDesc(values []string) {
	sort.Slice(values, func(i, j int) bool {
		if len(values[i]) == 0 && len(values[j]) == 0 {
			return false
		}
		if len(values[i]) == 0 || len(values[j]) == 0 {
			return len(values[i]) == 0
		}

		val1 := decimal.RequireFromString(values[i])
		val2 := decimal.RequireFromString(values[j])

		return val2.LessThan(val1)
	})
}

func PairQuote(value string) (pair string, quote string) {
	itens := strings.Split(value, "-") // 0 = PAIR 1 = QUOTE
	pair = itens[0]
	quote = itens[1]
	return
}
